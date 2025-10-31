package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"shanraq.org/pkg/shanraq"
	"shanraq.org/pkg/transport/respond"
)

// Handler processes a job payload.
type Handler func(context.Context, *shanraq.Runtime, Job) error

// Module adds HTTP endpoints and background workers for the job queue.
type Module struct {
	rt           *shanraq.Runtime
	store        *Store
	workerCount  int
	pollInterval time.Duration
	handlers     map[string]Handler
}

// JobContext key used for context values.
type jobContextKey struct{}

// JobContext exposes ambient metadata during job execution.
type JobContext struct {
	WorkerIndex int
	Attempts    int
}

// InfoFromContext extracts metadata about the running job from the context.
func InfoFromContext(ctx context.Context) (JobContext, bool) {
	jc, ok := ctx.Value(jobContextKey{}).(JobContext)
	return jc, ok
}

// Option customizes the jobs module.
type Option func(*Module)

// WithWorkerCount overrides the default worker parallelism.
func WithWorkerCount(n int) Option {
	return func(m *Module) {
		if n > 0 {
			m.workerCount = n
		}
	}
}

// WithPollInterval sets the fetch cadence for workers.
func WithPollInterval(d time.Duration) Option {
	return func(m *Module) {
		if d > 0 {
			m.pollInterval = d
		}
	}
}

// New creates a jobs module with sane defaults (2 workers, 1s poll interval).
func New(opts ...Option) *Module {
	m := &Module{
		workerCount:  2,
		pollInterval: time.Second,
		handlers:     map[string]Handler{},
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Module) Name() string {
	return "jobs"
}

// Handle registers a handler for the named job.
func (m *Module) Handle(name string, handler Handler) {
	m.handlers[name] = handler
}

// HandleFunc registers a handler function without requiring direct runtime access.
func (m *Module) HandleFunc(name string, fn func(context.Context, Job) error) {
	m.Handle(name, func(ctx context.Context, _ *shanraq.Runtime, job Job) error {
		return fn(ctx, job)
	})
}

// Init wires runtime dependencies.
func (m *Module) Init(_ context.Context, rt *shanraq.Runtime) error {
	m.rt = rt
	m.store = NewStore(rt.DB)
	return nil
}

// Routes exposes enqueue + inspection endpoints.
func (m *Module) Routes(r chi.Router) {
	if m.rt == nil {
		return
	}

	r.Route("/jobs", func(r chi.Router) {
		r.Post("/", m.handleEnqueue)
		r.Get("/", m.handleList)
		r.Post("/{id}/retry", m.handleRetry)
		r.Post("/{id}/cancel", m.handleCancel)
	})
}

// Start launches worker goroutines consuming jobs until ctx cancels.
func (m *Module) Start(ctx context.Context, rt *shanraq.Runtime) error {
	if m.store == nil {
		return errors.New("jobs store uninitialized")
	}
	if rt != nil {
		m.rt = rt
	}

	for i := 0; i < m.workerCount; i++ {
		go m.workerLoop(ctx, i)
	}

	<-ctx.Done()
	return ctx.Err()
}

func (m *Module) workerLoop(ctx context.Context, idx int) {
	ticker := time.NewTicker(m.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			job, err := m.store.ClaimNextJob(ctx)
			if err != nil {
				if errors.Is(err, ErrNoJobs) {
					continue
				}
				m.rt.Logger.Error("claim job", zap.Error(err))
				continue
			}
			m.processJob(ctx, job, idx)
		}
	}
}

func (m *Module) processJob(ctx context.Context, job Job, workerIdx int) {
	handler, ok := m.handlers[job.Name]
	if !ok {
		m.rt.Logger.Warn("job handler missing", zap.String("name", job.Name))
		_ = m.store.MarkFailed(ctx, job.ID, "handler missing")
		return
	}

	input := job
	m.rt.Logger.Info("job started", zap.String("job_id", job.ID.String()), zap.String("name", job.Name), zap.Int("worker", workerIdx))

	ctx = context.WithValue(ctx, jobContextKey{}, JobContext{
		WorkerIndex: workerIdx,
		Attempts:    job.Attempts,
	})

	if err := handler(ctx, m.rt, input); err != nil {
		m.rt.Logger.Warn("job errored", zap.String("job_id", job.ID.String()), zap.String("name", job.Name), zap.Error(err))
		if job.Attempts >= job.MaxAttempts {
			_ = m.store.MarkFailed(ctx, job.ID, err.Error())
			return
		}
		if err := m.store.MarkRetry(ctx, job.ID, err.Error()); err != nil {
			m.rt.Logger.Error("mark retry", zap.Error(err))
		}
		return
	}

	if err := m.store.MarkDone(ctx, job.ID); err != nil {
		m.rt.Logger.Error("mark done", zap.Error(err))
		return
	}
	m.rt.Logger.Info("job completed", zap.String("job_id", job.ID.String()), zap.String("name", job.Name))
}

func (m *Module) handleEnqueue(w http.ResponseWriter, r *http.Request) {
	var req enqueueRequest
	if err := respond.Decode(r, &req); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	if req.Name == "" {
		respond.Error(w, http.StatusBadRequest, errors.New("job name required"))
		return
	}
	if req.Payload == nil {
		req.Payload = map[string]any{}
	}

	runAt := time.Now()
	if req.RunAt != nil {
		runAt = *req.RunAt
	}
	maxAttempts := req.MaxAttempts
	if maxAttempts == 0 {
		maxAttempts = 3
	}

	payload, err := json.Marshal(req.Payload)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	job := Job{
		ID:          uuid.New(),
		Name:        req.Name,
		Payload:     payload,
		RunAt:       runAt,
		MaxAttempts: maxAttempts,
	}

	if err := m.store.Enqueue(r.Context(), job); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.JSON(w, http.StatusAccepted, map[string]string{
		"id": job.ID.String(),
	})
}

func (m *Module) handleList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()
	status := query.Get("status")
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	jobs, err := m.store.List(ctx, ListOptions{
		Status: status,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.JSON(w, http.StatusOK, jobs)
}

func (m *Module) handleRetry(w http.ResponseWriter, r *http.Request) {
	jobID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, errors.New("invalid job id"))
		return
	}
	if err := m.store.MarkPending(r.Context(), jobID); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{"status": "queued"})
}

func (m *Module) handleCancel(w http.ResponseWriter, r *http.Request) {
	jobID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respond.Error(w, http.StatusBadRequest, errors.New("invalid job id"))
		return
	}
	var body struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil && !errors.Is(err, io.EOF) {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	if err := m.store.Cancel(r.Context(), jobID, body.Reason); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.JSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

var _ interface {
	shanraq.Module
	shanraq.RouterModule
	shanraq.InitializerModule
	shanraq.StarterModule
} = (*Module)(nil)
