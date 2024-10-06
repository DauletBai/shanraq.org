package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"syscall/js"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5"
	"shanraq.org/internal/models"
)

type CategoryHandler struct {
	DB *sql.DB
}

func (h. *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		http.Error(w, "Unable to fetch categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			http.Error(w, "Error scaning category", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	json.NewEncoder(w) .Encode(categories)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	err := json.NewDecoder(r.Body) .Decode(&category)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO categories (name) VALUES ($1) RETURNING id"
	err = h.DB.QueryRow(query, category.Name) .Scan(&category.ID)
	if err != nil {
		http.Error(w, "Unable to create category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w) .Encode(category)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var category models.Category
	err := json.NewDecoder(r.Body) .Decode(&category)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := "DELETE FROM categories WHERE id = $1"
	_, err := h.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Unable to delete category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

