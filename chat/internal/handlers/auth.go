package handlers

import (
	"encoding/json"
	"net/http"

	"shanraq.org/chat/internal/db"
	"shanraq.org/chat/internal/models"
)

type AuthHandler struct {
	UserRepo *db.UserRepository
	JWTAuth *auth.JWTAuth
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body) .Decode(user)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid data format")
		return
	}

	// Password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Server error")
		return
	}
	user.Password = string(hashedPassword)

	// Create user
	err = h.UserRepo.CreateUser(&user)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, user)
}

// func for login & other...