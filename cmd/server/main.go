package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/admin"
)

func main() {
	router := app.SetupRouter()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))

	r.Route("/admin", func(r chi.Router) {
		r.Use(AdminOnly) // Middleware to check admin rights
		r.Get("/", adminHandler.Dashboard) // Admin Panel Start Page
		r.Route("/articles", func(r chi.Router) {
			r.Get("/", adminHandler.ManageArticles)
			r.Post("/", adminHandler.CreateArticle)
			r.Put("/{id}", adminHandler.UpdateArticle)
			r.Delete("/{id}", adminHandler.DeleteArticle)
		})
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", adminHandler.ManageCategories)
			r.Post("/", adminHandler.CreateCategory)
			r.Put("/{id}", adminHandler.UpdateCategory)
			r.Delete("/{id}", adminHandler.DeleteCategory)
		})
		r.Route("/users", func(r chi.Router) {
			r.Get("/", adminHandler.ManageUsers)
			r.Post("/", adminHandler.CreateUser)
			r.Put("/{id}", adminHandler.UpdateUser)
			r.Delete("/{id}", adminHandler.DeleteUser)
		})
	})
}
