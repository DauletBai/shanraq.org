package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"shanraq.org/internal/models"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, name, email, role FROM users")
	if err != nil {
		http.Error(w, "Unable to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role)
		if err != nil {
			http.Error(w, "Error scanning user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w) .Encode(users)
}

func (h *UserHandler) CreateUsers(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body) .Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) returning ID"
	err = h.DB.QueryRow(query, user.Name, user.Email, user.Password, user.Role) .Scan(&user.ID)
	if err != nil {
		http.Error(w, "Unable to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w) .Encode(user)
}

func (h *UserHandler) UpdateUsers(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var user models.User
	err := json.NewDecoder(r.Body) .Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := "UPDATE users SET name = $1, email = $2, role = $3, WHERE id = $4"
	_, err = h.DB.Exec(query, user.Name, user.Email, user.Role, id)
	if err != nil {
		http.Error(w, "Unable to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUsers(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	query := "DELETE FROM users WHERE id = $1"
	_, err := h.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Unable to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}