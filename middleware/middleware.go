package middleware

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

type ContextKeys string

const (
	userContext ContextKeys = "__userContext"
)

var mySigningKey = []byte("secret_key")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		apikey := request.Header.Get("x-api-key")
		//userid := request.Header.Get("userid")
		token, TokenErr := jwt.Parse(apikey, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error")
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
			//	user, err := helper.GetUser(apikey)
			//utils.CheckError(err)

			ctx := context.WithValue(request.Context(), userContext, "")
			next.ServeHTTP(writer, request.WithContext(ctx))

		} else {
			fmt.Fprintf(writer, " PLEASE LOGIN AGAIN")
			return
		}
	})
}

//func UserContext(request *http.Request) *model.User {
//	//	var user model.User
//	user, ok := request.Context().Value(userContext).(*model.User)
//	if ok && user != nil {
//		return user
//	}
//	return nil
//}
