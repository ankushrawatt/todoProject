package route

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"todoproject/database/helper"
	"todoproject/handler"
)

type ContextKeys string

const (
	userContext ContextKeys = "__userContext"
)

type Server struct {
	chi.Router
}

func Middleware(handle http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		apikey := request.Header.Get("x-api-key")
		user, err := helper.GetUserSession(apikey)
		if err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(fmt.Sprintf("Please Login")))
			panic(err)
		}
		ctx := context.WithValue(request.Context(), userContext, user)
		handle.ServeHTTP(writer, request.WithContext(ctx))
	}
}

func Route() *Server {
	router := chi.NewRouter()
	router.Route("/todo", func(todo chi.Router) {
		todo.Post("/signup", handler.Signup)
		todo.Post("/login", handler.Login)
		todo.Post("/createtask", Middleware(handler.CreateTask))
		todo.Post("/updatetask", Middleware(handler.UpdateTask))
		todo.Get("/alltask", Middleware(handler.AllTask))
		todo.Get("/logout", Middleware(handler.Logout))
	})
	return &Server{router}
}

func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
