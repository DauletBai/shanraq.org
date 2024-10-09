package app

import (
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"shanraq.org/config"
	"shanraq.org/internal/handlers"
	"shanraq.org/internal/repositories"
	"shanraq.org/internal/services"
)

func SetupRouter(cfg *config.Config, logger *logrus.Logger) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Connection to database
	db, err := sql.Open("postgres", cfg.Database.DSN())
	if err != nil {
		logger.Fatalf("Database connection error: %v", err)
	}

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	// Other repositories...

	// Services
	iwtService := services.JWTService{SecretKey: cfg.JWT.SekretKey}
	authService := services.AuthService{UserRepo: userRepo, TokenService: services.JWTService}
	// Other services...

	// Heandlers
	authHeandler := handlers.NewAuthHeandler(authService)
	// Other heandlers...

	// Routers
	r.Post("/register", authHeandler.Register)
	r.Post("/login", authHeandler.login)

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(jwtService))

		r.Route("/articles", func(r chi.Router) {
			r.Get("/", handlers.GetArticles)
			r.Post("/", handlers.CreateArticle)
			r.Put("/{id}", handlers.UpdateArticle)
			r.Delete("/{id}", handlers.DeleteArticle)
		})

		// Admin dashboard routers
		r.Group(func(r chi.Router) {
			r.Use(AdminOnlyMiddleware)
			r.Route("/admin", func(r chi.Router) {
				// Heandler admin dashboard
			})
		})
	})

	return r

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