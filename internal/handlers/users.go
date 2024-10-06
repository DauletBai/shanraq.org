package handlers

import (
	"database/sql"
	"net/http"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) CreateUsers(w http.ResponseWriter, r *http.Request) {
	
}

func (h *UserHandler) UpdateUsers(w http.ResponseWriter, r *http.Request) {
	
}

func (h *UserHandler) DeleteUsers(w http.ResponseWriter, r *http.Request) {
	
}