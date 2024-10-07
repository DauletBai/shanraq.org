package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"shanraq.org/internal/models"
)

type AdminHandler struct {
	DB *sql.DB
	Tmpl *template.Template
}

func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Admin Dashboard",
	}
	h.Tmpl.ExecuteTemplate(w, "dashboard.htnl", data)
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context() .Value("user").(*models.User)
		if !ok || (user.Role != "admin" && user.Role != "editor") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *AdminHandler) ManageArticles(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, title FROM articles")
	if err != nil {
		http.Error(w, "Unable to fetch articles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.Title)
		if err != nil {
			http.Error(w, "Error scanning article", http.StatusInternalServerError)
			return
		}

		articles = append(articles, article)
	}

	data := map[string]interface{}{
		"Title": "ManageArticles",
		"Articles": articles,
	}
	h.Tmpl.ExecuteTemplate(w, "manage_articles.html", data)
}