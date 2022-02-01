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

var mySigningKey = []byte("secret_key")

//func Middleware(handle http.HandlerFunc) http.HandlerFunc {
//	return func(writer http.ResponseWriter, request *http.Request) {
//		apikey := request.Header.Get("x-api-key")
//		userid := request.Header.Get("userid")
//		token, TokenErr := jwt.Parse(apikey, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, fmt.Errorf("There was an error")
//			}
//			return mySigningKey, nil
//		})
//		if TokenErr != nil {
//			fmt.Fprintf(writer, TokenErr.Error())
//
//		}
//		if token.Valid {
//			user := helper.GetUserSession(apikey, userid)
//
//			ctx := context.WithValue(request.Context(), userContext, user)
//			handle.ServeHTTP(writer, request.WithContext(ctx))
//			handle(writer, request)
//
//		} else {
//			fmt.Fprintf(writer, " PLEASE LOGIN AGAIN")
//
//		}

//user, err := helper.GetUserSession(apikey, userid)
//if err != nil {
//	writer.WriteHeader(http.StatusUnauthorized)
//	writer.Write([]byte(fmt.Sprintf("Please Login")))
//	panic(err)
//}
//ctx := context.WithValue(request.Context(), userContext, user)
//handle.ServeHTTP(writer, request.WithContext(ctx))

//COOKIES
//	handler.CheckCookies(writer, request)
//	}
//}

func Middleware(handle http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		apikey := request.Header.Get("x-api-key")
		//fmt.Println(apikey)
		//	userid := request.Header.Get("userid")
		user, err := helper.GetUserSession(apikey)
		if err != nil || user == "" {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(fmt.Sprintf("Please Login")))
			panic(err)
		}
		fmt.Println(user)
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
		todo.Get("/upcomingtodo", Middleware(handler.UpcomingTodo))
		todo.Get("/expiredtodo", Middleware(handler.ExpiredTodo))
		todo.Delete("/logout", Middleware(handler.Logout))

	})
	return &Server{router}
}

func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
