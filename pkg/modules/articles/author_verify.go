package articles

import (
	"net/http"

	"go.uber.org/zap"

	"shanraq.org/pkg/modules/auth"
	"shanraq.org/pkg/site"
)

// VerifyAuthorPage backs the Studio "become a verified author" flow.
type VerifyAuthorPage struct {
	Base
	First         string
	Last          string
	Middle        string
	Phone         string
	PhoneVerified bool
	CanPublish    bool
	CodeSent      bool
	NeedPublish   bool
	Notice        string
	Error         string
}

func (m *Module) renderVerifyAuthor(w http.ResponseWriter, r *http.Request, p VerifyAuthorPage) {
	lang := site.ResolveLang(w, r)
	p.Base = m.base(r, site.T(lang, "author.verify_title"), lang)
	uid, ok := m.authorID(r)
	if ok {
		first, last, verified := m.auth.AuthorIdentity(r.Context(), uid)
		p.First, p.Last, p.PhoneVerified = first, last, verified
		p.Middle = m.auth.MiddleName(r.Context(), uid)
		p.CanPublish = m.auth.CanPublish(r.Context(), uid)
		if !verified {
			if pending, active, err := m.auth.PendingPhoneCode(r.Context(), uid); err != nil {
				m.rt.Logger.Error("load pending phone code", zap.Error(err))
			} else if active {
				p.CodeSent = true
				if p.Phone == "" {
					p.Phone = pending
				}
			}
		}
	}
	q := r.URL.Query()
	p.NeedPublish = q.Get("need") == "publish"
	switch {
	case q.Get("name") == "ok":
		p.Notice = site.T(lang, "author.name_saved")
	case q.Get("verified") == "ok":
		p.Notice = site.T(lang, "author.phone_ok")
	case q.Get("code") == "bad":
		p.Error = site.T(lang, "author.code_bad")
	case q.Get("phone") == "bad":
		p.Error = site.T(lang, "author.err_phone")
	}
	m.render(w, "verify_author", p)
}

func (m *Module) handleAuthorVerifyPage(w http.ResponseWriter, r *http.Request) {
	m.renderVerifyAuthor(w, r, VerifyAuthorPage{})
}

func (m *Module) handleAuthorName(w http.ResponseWriter, r *http.Request) {
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	_ = r.ParseForm()
	if err := m.auth.SetAuthorName(r.Context(), uid, r.FormValue("first_name"), r.FormValue("last_name"), r.FormValue("middle_name")); err != nil {
		m.renderVerifyAuthor(w, r, VerifyAuthorPage{Error: site.T(site.ResolveLang(w, r), "author.err_name")})
		return
	}
	http.Redirect(w, r, "/studio/author?name=ok", http.StatusSeeOther)
}

func (m *Module) handleAuthorPhone(w http.ResponseWriter, r *http.Request) {
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	_ = r.ParseForm()
	phone := r.FormValue("phone")
	if !m.auth.AllowAuthAttempt(r, "signin", r.FormValue("phone")) { // reuse the signin limiter for OTP sends
		m.renderVerifyAuthor(w, r, VerifyAuthorPage{Phone: phone, Error: site.T(site.ResolveLang(w, r), "form.err_rate_limit")})
		return
	}
	channel := r.FormValue("channel")
	if channel != "telegram" && channel != "sms" {
		channel = "sms"
	}
	if err := m.auth.StartPhoneVerificationVia(r.Context(), uid, phone, channel); err != nil {
		if err == auth.ErrInvalidPhone {
			m.renderVerifyAuthor(w, r, VerifyAuthorPage{Phone: phone, Error: site.T(site.ResolveLang(w, r), "author.err_phone")})
			return
		}
		m.rt.Logger.Error("start phone verification", zap.String("user_id", uid.String()), zap.String("channel", channel), zap.Error(err))
		m.renderVerifyAuthor(w, r, VerifyAuthorPage{Phone: phone, Error: site.T(site.ResolveLang(w, r), "author.err_delivery")})
		return
	}
	m.rt.Logger.Info("phone verification request accepted", zap.String("user_id", uid.String()), zap.String("channel", channel))
	http.Redirect(w, r, "/studio/author?code=sent", http.StatusSeeOther)
}

func (m *Module) handleAuthorConfirm(w http.ResponseWriter, r *http.Request) {
	uid, ok := m.authorID(r)
	if !ok {
		http.Redirect(w, r, "/studio/login", http.StatusSeeOther)
		return
	}
	_ = r.ParseForm()
	verified, err := m.auth.ConfirmPhoneVerification(r.Context(), uid, r.FormValue("code"))
	if err != nil {
		m.rt.Logger.Error("confirm phone verification", zap.Error(err))
	}
	if verified {
		http.Redirect(w, r, "/studio/author?verified=ok", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/studio/author?code=bad", http.StatusSeeOther)
}
