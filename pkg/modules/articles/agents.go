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
	"shanraq.org/pkg/modules/auth"
)

// Agent is a real-estate agent's public profile. One per user; registration is
// free and self-serve from the cabinet.
type Agent struct {
	UserID string
	Name   string
	Agency string
	Phone  string
	About  string
}

// AgentStore persists real-estate agent profiles.
type AgentStore struct{ db *pgxpool.Pool }

func NewAgentStore(db *pgxpool.Pool) *AgentStore { return &AgentStore{db: db} }

// ByUser loads a user's agent profile, or nil if they have not registered.
func (s *AgentStore) ByUser(ctx context.Context, userID uuid.UUID) (*Agent, error) {
	var a Agent
	var id uuid.UUID
	err := s.db.QueryRow(ctx, `
		SELECT user_id, name, agency, phone, about FROM re_agents WHERE user_id = $1`, userID).
		Scan(&id, &a.Name, &a.Agency, &a.Phone, &a.About)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("agent by user: %w", err)
	}
	a.UserID = id.String()
	return &a, nil
}

// Save registers or updates the caller's agent profile (one per user).
func (s *AgentStore) Save(ctx context.Context, userID uuid.UUID, a Agent) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO re_agents (user_id, name, agency, phone, about)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
			name = EXCLUDED.name, agency = EXCLUDED.agency,
			phone = EXCLUDED.phone, about = EXCLUDED.about, updated_at = NOW()`,
		userID, a.Name, a.Agency, a.Phone, a.About)
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
	} else if claims, ok := auth.ClaimsFromContext(r.Context()); ok {
		page.Draft.Name = displayName(claims.Email) // sensible default for the form
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
	_ = r.ParseForm()
	a := Agent{
		Name:   strings.TrimSpace(r.FormValue("name")),
		Agency: strings.TrimSpace(r.FormValue("agency")),
		Phone:  strings.TrimSpace(r.FormValue("phone")),
		About:  strings.TrimSpace(r.FormValue("about")),
	}
	if a.Name == "" {
		existing, _ := m.reagents.ByUser(r.Context(), uid)
		page := AgentCabinetPage{Base: m.base(r, T(lang, "agent.title"), lang), Agent: existing, Draft: a, Error: T(lang, "agent.err_name")}
		page.ActiveCat = "realestate"
		m.render(w, "agent_cabinet", page)
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
	if agent == nil {
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
