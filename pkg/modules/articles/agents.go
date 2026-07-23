package articles

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Agent is a real-estate agent's public profile. One per user; registration is
// free and self-serve from the cabinet.
type Agent struct {
	UserID    string
	FirstName string
	LastName  string
	Name      string // composed display name (first + last), kept for the badge/join
	Agency    string
	Phone     string
	About     string
	Status    string // pending | verified | rejected
	Reason    string // rejection reason, shown to the agent
	Email     string // owner email, for the admin queue
}

// Verified reports whether the profile has been approved by an admin.
func (a Agent) Verified() bool { return a.Status == agentVerified }

const (
	agentPending  = "pending"
	agentVerified = "verified"
	agentRejected = "rejected"
)

func isAgentStatus(s string) bool {
	return s == agentPending || s == agentVerified || s == agentRejected
}

// composeName builds the display name from first + last.
func composeName(first, last string) string {
	return strings.TrimSpace(strings.TrimSpace(first) + " " + strings.TrimSpace(last))
}

// AgentStore persists real-estate agent profiles.
type AgentStore struct{ db *pgxpool.Pool }

func NewAgentStore(db *pgxpool.Pool) *AgentStore { return &AgentStore{db: db} }

// ByUser loads a user's agent profile, or nil if they have not registered.
func (s *AgentStore) ByUser(ctx context.Context, userID uuid.UUID) (*Agent, error) {
	var a Agent
	var id uuid.UUID
	err := s.db.QueryRow(ctx, `
		SELECT user_id, first_name, last_name, name, agency, phone, about, status, reject_reason
		FROM re_agents WHERE user_id = $1`, userID).
		Scan(&id, &a.FirstName, &a.LastName, &a.Name, &a.Agency, &a.Phone, &a.About, &a.Status, &a.Reason)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("agent by user: %w", err)
	}
	a.UserID = id.String()
	return &a, nil
}

// Pending returns agent profiles awaiting review, oldest first, with the
// owner's email for the admin verification queue.
func (s *AgentStore) Pending(ctx context.Context, limit int) ([]Agent, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := s.db.Query(ctx, `
		SELECT a.user_id, a.first_name, a.last_name, a.name, a.agency, a.phone, a.about, a.status, COALESCE(u.email,'')
		FROM re_agents a JOIN auth_users u ON u.id = a.user_id
		WHERE a.status = 'pending'
		ORDER BY a.created_at ASC LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("pending agents: %w", err)
	}
	defer rows.Close()
	var out []Agent
	for rows.Next() {
		var a Agent
		var id uuid.UUID
		if err := rows.Scan(&id, &a.FirstName, &a.LastName, &a.Name, &a.Agency, &a.Phone, &a.About, &a.Status, &a.Email); err != nil {
			return nil, err
		}
		a.UserID = id.String()
		out = append(out, a)
	}
	return out, rows.Err()
}

// SetStatus approves or rejects an agent profile, recording who decided and why.
func (s *AgentStore) SetStatus(ctx context.Context, userID uuid.UUID, status, reason string, reviewer *uuid.UUID) error {
	if !isAgentStatus(status) {
		return fmt.Errorf("bad agent status")
	}
	_, err := s.db.Exec(ctx, `
		UPDATE re_agents SET status = $2, reject_reason = $3, reviewed_at = NOW(), reviewed_by = $4
		WHERE user_id = $1`, userID, status, reason, reviewer)
	if err != nil {
		return fmt.Errorf("set agent status: %w", err)
	}
	return nil
}

// Save registers or updates the caller's agent profile (one per user).
func (s *AgentStore) Save(ctx context.Context, userID uuid.UUID, a Agent) error {
	name := composeName(a.FirstName, a.LastName)
	_, err := s.db.Exec(ctx, `
		INSERT INTO re_agents (user_id, first_name, last_name, name, agency, phone, about)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id) DO UPDATE SET
			first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name, name = EXCLUDED.name,
			agency = EXCLUDED.agency, phone = EXCLUDED.phone, about = EXCLUDED.about, updated_at = NOW()`,
		userID, strings.TrimSpace(a.FirstName), strings.TrimSpace(a.LastName), name, a.Agency, a.Phone, a.About)
	if err != nil {
		return fmt.Errorf("save agent: %w", err)
	}
	return nil
}

// ---- pages ----

// AgentCabinetPage backs the agent's own cabinet.
type AgentCabinetPage struct {
	Base
	Agent *Agent // nil until registered
	Draft Agent  // echoed values on validation failure / prefill
	Count int    // active listings
	Saved bool
	Error string
}

// AgentPublicPage backs the public agent page (profile + all their listings).
type AgentPublicPage struct {
	Base
	Agent    *Agent
	Listings []*Listing
}

func (m *Module) handleAgentCabinet(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	agent, err := m.reagents.ByUser(r.Context(), uid)
	if err != nil {
		m.rt.Logger.Error("agent load", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	page := AgentCabinetPage{Base: m.base(r, T(lang, "agent.title"), lang), Agent: agent}
	page.ActiveCat = "realestate"
	page.Saved = r.URL.Query().Get("saved") == "1"
	if agent != nil {
		page.Draft = *agent
		if items, lerr := m.listings.AgentListings(r.Context(), uid); lerr == nil {
			page.Count = len(items)
		}
	}
	m.render(w, "agent_cabinet", page)
}

func (m *Module) handleAgentSave(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	// Registration can be paused (invite-only beta) from the admin panel.
	if !m.flags.Available(SvcAgentReg) {
		http.Redirect(w, r, "/agent", http.StatusSeeOther)
		return
	}
	_ = r.ParseForm()
	a := Agent{
		FirstName: strings.TrimSpace(r.FormValue("first_name")),
		LastName:  strings.TrimSpace(r.FormValue("last_name")),
		Agency:    strings.TrimSpace(r.FormValue("agency")),
		Phone:     strings.TrimSpace(r.FormValue("phone")),
		About:     strings.TrimSpace(r.FormValue("about")),
	}
	a.Name = composeName(a.FirstName, a.LastName)

	fail := func(msg string) {
		existing, _ := m.reagents.ByUser(r.Context(), uid)
		page := AgentCabinetPage{Base: m.base(r, T(lang, "agent.title"), lang), Agent: existing, Draft: a, Error: msg}
		page.ActiveCat = "realestate"
		m.render(w, "agent_cabinet", page)
	}
	if a.FirstName == "" {
		fail(T(lang, "agent.err_name"))
		return
	}
	// An agent is a trust signal, so the bar is a verified email AND phone — the
	// same identity proof the platform already requires to post and to author.
	if !m.auth.IsEmailVerified(r.Context(), uid) || !m.auth.IsPhoneVerified(r.Context(), uid) {
		fail(T(lang, "agent.err_verify"))
		return
	}
	if err := m.reagents.Save(r.Context(), uid, a); err != nil {
		m.rt.Logger.Error("agent save", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/agent?saved=1", http.StatusSeeOther)
}

func (m *Module) handleAgentPublic(w http.ResponseWriter, r *http.Request) {
	lang := m.resolveLang(w, r)
	uid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	agent, err := m.reagents.ByUser(r.Context(), uid)
	if err != nil {
		m.rt.Logger.Error("agent public load", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	// Only verified agents have a public page; a pending/rejected profile is
	// private (no trust surface until an admin approves it).
	if agent == nil || !agent.Verified() {
		http.NotFound(w, r)
		return
	}
	items, err := m.listings.AgentListings(r.Context(), uid)
	if err != nil {
		m.rt.Logger.Warn("agent listings", zap.Error(err))
	}
	page := AgentPublicPage{Base: m.base(r, agent.Name, lang), Agent: agent, Listings: items}
	page.ActiveCat = "realestate"
	m.render(w, "agent_public", page)
}
