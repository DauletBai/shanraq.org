// github.com/DauletBai/shanraq.org/http/router.go
package http

import (
	"net/http" 

	"github.com/go-chi/chi/v5"

	k "github.com/DauletBai/shanraq.org/core/kernel"
)

// Router is a wrapper around chi.Router.
type Router struct { // Тип Router экспортируемый
	chiRouter *chi.Mux
	kernel    *k.Kernel
}

// NewRouter creates a new Shanraq Router. (Функция NewRouter экспортируемая)
func NewRouter(kernel *k.Kernel) *Router {
	chiMux := chi.NewRouter()
	return &Router{
		chiRouter: chiMux,
		kernel:    kernel,
	}
}

func (r *Router) ChiRouter() *chi.Mux {
	return r.chiRouter
}

func (r *Router) Kernel() *k.Kernel {
	return r.kernel
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chiRouter.ServeHTTP(w, req)
}

// adaptHandler использует Context и HandlerFunc из текущего пакета http
func (r *Router) adaptHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		shanraqCtx := NewContext(w, req, r.kernel) // NewContext из context.go (тот же пакет)
		handler(shanraqCtx)
	}
}

// Use принимает chi-совместимые middleware
func (r *Router) Use(middlewares ...func(next http.Handler) http.Handler) {
	for _, mw := range middlewares {
		r.chiRouter.Use(mw)
	}
}

// GET, POST и т.д. используют HandlerFunc из текущего пакета http
func (r *Router) GET(pattern string, handler HandlerFunc) {
	r.chiRouter.Get(pattern, r.adaptHandler(handler))
}

func (r *Router) POST(pattern string, handler HandlerFunc) {
	r.chiRouter.Post(pattern, r.adaptHandler(handler))
}

func (r *Router) PUT(pattern string, handler HandlerFunc) {
	r.chiRouter.Put(pattern, r.adaptHandler(handler))
}

func (r *Router) DELETE(pattern string, handler HandlerFunc) {
	r.chiRouter.Delete(pattern, r.adaptHandler(handler))
}

func (r *Router) PATCH(pattern string, handler HandlerFunc) {
	r.chiRouter.Patch(pattern, r.adaptHandler(handler))
}

func (r *Router) OPTIONS(pattern string, handler HandlerFunc) {
	r.chiRouter.Options(pattern, r.adaptHandler(handler))
}

func (r *Router) Mount(pattern string, handler http.Handler) {
	r.chiRouter.Mount(pattern, handler)
}

func (r *Router) Group(fn func(sr *Router)) {
	r.chiRouter.Group(func(chiSubRouter chi.Router) {
		subRouter := &Router{
			chiRouter: chiSubRouter.(*chi.Mux),
			kernel:    r.kernel,
		}
		fn(subRouter)
	})
}

func (r *Router) Route(pattern string, fn func(sr *Router)) {
	r.chiRouter.Route(pattern, func(chiSubRouter chi.Router) {
		subRouter := &Router{
			chiRouter: chiSubRouter.(*chi.Mux),
			kernel:    r.kernel,
		}
		fn(subRouter)
	})
}