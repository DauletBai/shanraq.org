package app

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

	return r
}