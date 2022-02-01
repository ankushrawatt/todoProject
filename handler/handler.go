package handler

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"net/http"
	"time"
	"todoproject/database/helper"
	"todoproject/model"
)

var mySigningKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateToken(userid string) (string, error) {
	id := uuid.New()
	err := helper.CreateSession(id.String(), userid)
	if err != nil {
		return "", err
	}

	return id.String(), nil

	//token := jwt.New(jwt.SigningMethodHS256)
	//claims := token.Claims.(jwt.MapClaims)
	//claims["authorized"] = true
	//claims["user"] = userid
	//claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	//tokenString, err := token.SignedString(mySigningKey)
	//if err != nil {
	//	fmt.Errorf("Something went wrong %s", err.Error())
	//	return "", err
	//}
	//return tokenString, nil
}

func Signup(writer http.ResponseWriter, request *http.Request) {
	var user model.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	id := uuid.New()
	userId, NewErr := helper.NewUser(id.String(), user.Email, user.UserId, user.Password, user.FirstName, user.LastName, user.MobileNo)
	if NewErr != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	//_, jsonErr := json.Marshal(userId)
	//if jsonErr != nil {
	//	writer.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	writer.Write([]byte(fmt.Sprintf("Thank You %s For signing up...", userId)))
}

func Login(writer http.ResponseWriter, request *http.Request) {
	var Cred model.LoginUser
	err := json.NewDecoder(request.Body).Decode(&Cred)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	loginUser, loginErr := helper.Login(Cred.Email, Cred.Password)
	if loginErr != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(fmt.Sprintf("WRONG CREDENTIALS")))
		return
	}
	_, jsonErr := json.Marshal(loginUser)
	if jsonErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	token, tokenErr := CreateToken(loginUser)
	if tokenErr != nil {
		writer.WriteHeader(http.StatusBadGateway)
		return
	}
	writer.Write([]byte(fmt.Sprintf("Tokken: %s", token)))
	//writer.Write(jsonData)

	//trying cookie
	//expirationTime := time.Now().Add(time.Minute * 20)
	//claims := &Claims{
	//	Username: loginUser,
	//	StandardClaims: jwt.StandardClaims{
	//		ExpiresAt: expirationTime.Unix(),
	//	},
	//}
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//tokenString, err := token.SignedString(jwtkey)
	//if err != nil {
	//	writer.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//http.SetCookie(writer, &http.Cookie{
	//	Name:    "token",
	//	Value:   tokenString,
	//	Expires: expirationTime,
	//})

	writer.Write([]byte(fmt.Sprintf(" WELCOME userid: %s", loginUser)))

}

func CreateTask(writer http.ResponseWriter, request *http.Request) {
	var task model.TodoTask
	err := json.NewDecoder(request.Body).Decode(&task)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	userkey := request.Header.Get("userid")
	todo, todoErr := helper.CreateTodo(task.Des, userkey, task.Task, task.Date)
	if todoErr != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	writer.Write([]byte(fmt.Sprintf("New task created %s", todo)))
}
func UpdateTask(writer http.ResponseWriter, request *http.Request) {

}

func AllTask(writer http.ResponseWriter, request *http.Request) {

	userkey := request.Header.Get("userid")
	tasks, taskErr := helper.TodoList(userkey)
	if taskErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, jsonErr := json.Marshal(tasks)
	if jsonErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(jsonData)
}

func UpcomingTodo(writer http.ResponseWriter, request *http.Request) {
	//	date := time.Now()
	date := time.Now()
	userkey := request.Header.Get("userid")
	todo, todoErr := helper.UpcomingTodoList(userkey, date)
	if todoErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, jsonErr := json.Marshal(todo)
	if jsonErr != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.Write(jsonData)
}

func ExpiredTodo(writer http.ResponseWriter, request *http.Request) {
	date := time.Now()
	userKey := request.Header.Get("userid")
	todo, err := helper.ExpiredTodo(userKey, date)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, jsonErr := json.Marshal(todo)
	if jsonErr != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.Write(jsonData)
}

func Logout(writer http.ResponseWriter, request *http.Request) {
	//claim:=
	apikey := request.Header.Get("x-api-key")
	err := helper.DeleteSession(apikey)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	user := request.Header.Get("userid")
	writer.Write([]byte(fmt.Sprintf("%s USER LOGGED OUT SUCCESSFULLY", user)))
}

//cookie
//
//func CheckCookies(writer http.ResponseWriter, request *http.Request) {
//	cookie, err := request.Cookie("token")
//	if err != nil {
//		if err == http.ErrNoCookie {
//			writer.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//		writer.WriteHeader(http.StatusBadGateway)
//		return
//	}
//	tokenStr := cookie.Value
//	claims := &Claims{}
//	tkn, NewErr := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
//		return jwtkey, nil
//	})
//	if NewErr != nil {
//		if NewErr == jwt.ErrSignatureInvalid {
//			writer.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//		writer.WriteHeader(http.StatusBadGateway)
//		return
//	}
//	if !tkn.Valid {
//		writer.WriteHeader(http.StatusUnauthorized)
//		return
//	}
//}
