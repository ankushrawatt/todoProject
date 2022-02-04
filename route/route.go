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
			todo.Use(middleware.AuthMiddleware)
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
