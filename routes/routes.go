package routes

import (
	"gochitest/controllers"
	"gochitest/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Setup() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", controllers.Index)

	r.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/register", controllers.Register)
			r.Post("/login", controllers.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(middlewares.IsAuthenticated)
			r.Get("/", controllers.Users)
			r.Post("/update", controllers.UpdateInfo)
			r.Post("/logout", controllers.Logout)
		})
	})

	return r
}
