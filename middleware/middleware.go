package middleware

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"todoproject/utils"
)

type ContextKeys string

const (
	userContext ContextKeys = "__userContext"
)

var mySigningKey = []byte("secret_key")

type JWTClaims struct {
	Userid string `json:"user"`
	//role   string `json:"role"`
	jwt.StandardClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		apikey := request.Header.Get("x-api-key")
		claims := &JWTClaims{}
		token, TokenErr := jwt.ParseWithClaims(apikey, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error")
			}
			return mySigningKey, nil
		})
		if TokenErr != nil {
			utils.CheckError(TokenErr)

		}
		if token.Valid {
			fmt.Println(claims.Userid)
			ctx := context.WithValue(request.Context(), "claims", claims.Userid)
			next.ServeHTTP(writer, request.WithContext(ctx))

		} else {
			_, err := fmt.Fprintf(writer, " PLEASE LOGIN AGAIN")
			utils.CheckError(err)
			return
		}
	})
}
