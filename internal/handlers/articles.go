package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5"
	"shanraq.org/internal/models"
)

// database injection framework
type ArtticleHandler struct {
	DB *sql.DB
} 

func (h *ArtticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, title, content, author_id FROM articles")
	if err != nil {
		http.Error(w, "Unable to fetch articles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	var article []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.AuthorID)
		if err != nil {
			http.Error(w, "Error scanning article", http.StatusInternalServerError)
			return
		}
		articles = append(articles, article)
	}

	json.NewEncoder(w) .Encode(articles)
}

func (h *ArtticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var article models.Article
	err := json.NewDecoder(r.Body) .Decode(&article)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO articles(title, content, author_id) VALUES ($1, $2, $3) RETURNING id"
		err = h.DB.QueryRow(query, article.Title, article.Content, article.AuthorID) .Scan(&article.ID)
		if err != nil {
			http.Error(w, "Unable to create article", http.StatusInternalServerError)
			return
		}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w) .Encode(article)
}

func (h *ArtticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var article models.Article
	err := json.NewDecoder(r.Body) .Decode(&article)
	if err != nil {
		http.Error(w, "Invalid input",  http.StatusBadRequest)
		return
	}

	query := "UPDATE articles SET title = $1, content = $2 WHERE id = $3"
	_, err h.DB.Exec(query, article.Title, article.Content, id)
	if err != nil {
		http.Error(w, "Unable to update article", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ArtticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := "DELETE FROM articles WHERE id = $1"
	_, err := h.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Unable to delete article", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}