package app

import (
	"github.com/go-chi/chi"
	"shanraq.org/internal/handlers"
)

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Router("/articles", func(r chi.Router) {
		r.Get("/", handlers.GetArticles)
		r.Post("/", handlers.CreateArticle)
		r.Put("/{id}", handlers.UpdateArticle)
		r.Delete("/{id}", handlers.DeleteArticle)
	})

	r.Router("/categories", func(r chi.Router) {
		r.Get("/", handlers.GetCategories)
		r.Post("/", handlers.CreateCategory)
		r.Put("/{id}", handlers.UpdateCategory)
		r.Delete("/{id}", handlers.DeleteCategory)
	})

	r.Router("/users", func(r chi.Router) {
		r.Get("/", handlers.GetUsers)
		r.Post("/", handlers.CreateUser)
		r.Put("/{id}", handlers.UpdateUser)
		r.Delete("/{id}", handlers.DeleteUser)
	})

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

	return r
}