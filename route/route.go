package route

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"net/http"
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		apikey := request.Header.Get("x-api-key")
		//userid := request.Header.Get("userid")
		token, TokenErr := jwt.Parse(apikey, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return mySigningKey, nil
		})
		if TokenErr != nil {
			fmt.Fprintf(writer, TokenErr.Error())

		}
		if token.Valid {
			//user, sessionErr := helper.CheckSession(apikey, userid)
			//if sessionErr != nil {
			//	writer.WriteHeader(http.StatusBadRequest)
			//	return
			//}
			ctx := context.WithValue(request.Context(), userContext, "")
			next.ServeHTTP(writer, request.WithContext(ctx))

		} else {
			fmt.Fprintf(writer, " PLEASE LOGIN AGAIN")
			return
		}
	})
}

//COOKIES
//	handler.CheckCookies(writer, request)
//	}
//}

//SESSION TABLE
//func Middleware(handle http.HandlerFunc) http.HandlerFunc {
//	return func(writer http.ResponseWriter, request *http.Request) {
//		apikey := request.Header.Get("x-api-key")
//		//fmt.Println(apikey)
//		//	userid := request.Header.Get("userid")
//		user, err := helper.GetUserSession(apikey)
//		if err != nil || user == "" {
//			writer.WriteHeader(http.StatusUnauthorized)
//			writer.Write([]byte(fmt.Sprintf("Please Login")))
//			panic(err)
//		}
//		ctx := context.WithValue(request.Context(), userContext, user)
//		handle.ServeHTTP(writer, request.WithContext(ctx))
//	}
//}POST

func Route() *Server {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Route("/todo", func(todo chi.Router) {
			todo.Use(AuthMiddleware)
			todo.Get("/", handler.AllTask)
			todo.Post("/", handler.CreateTask)
			todo.Get("/upcoming", handler.UpcomingTodo)
			todo.Get("/expired", handler.ExpiredTodo)

			todo.Route("/{id}", func(todoAction chi.Router) {
				todoAction.Put("/", handler.UpdateTask)
				todoAction.Delete("/", handler.DeleteTodo)
			})
		})
		r.Route("/user", func(userRoute chi.Router) {
			userRoute.Use(AuthMiddleware)
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
