package route

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"todoproject/handler"
	"todoproject/middleware"
)

type Server struct {
	chi.Router
}

func Route() *Server {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Route("/todo", func(todo chi.Router) {
			todo.Use(middleware.AuthMiddleware)
			todo.Get("/", handler.AllTodo)
			todo.Post("/", handler.CreateTodo)
			todo.Get("/upcoming", handler.UpcomingTodo)
			todo.Get("/expired", handler.ExpiredTodo)

			todo.Route("/{id}", func(todoAction chi.Router) {
				todoAction.Put("/", handler.UpdateTodo)
				todoAction.Delete("/", handler.DeleteTodo)
			})
		})
		r.Route("/user", func(userRoute chi.Router) {
			userRoute.Use(middleware.AuthMiddleware)
			userRoute.Delete("/", handler.DeleteUser)
			userRoute.Get("/", handler.Logout)
		})
		r.Route("/public", func(public chi.Router) {
			public.Post("/signup", handler.Signup)
			public.Post("/login", handler.Login)
			public.Post("/forget_password", handler.ResetPassword)
		})
	})
	return &Server{router}
}

func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
