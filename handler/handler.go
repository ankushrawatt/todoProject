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
	"todoproject/utils"
)

var mySigningKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateToken(userid string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = userid
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	tokenString, err := token.SignedString(mySigningKey)
	utils.CheckError(err)
	return tokenString, nil
}

func GetUserId(request *http.Request) string {
	token := request.Header.Get("x-api-key")
	userID, err := helper.GetUser(token)
	utils.CheckError(err)
	return userID
	//return ""
}

func Signup(writer http.ResponseWriter, request *http.Request) {
	var user model.User
	err := json.NewDecoder(request.Body).Decode(&user)
	utils.CheckError(err)
	id := uuid.New()
	userId, NewErr := helper.NewUser(id.String(), user.Email, user.UserId, user.Password, user.FirstName, user.LastName, user.MobileNo)
	utils.CheckError(NewErr)

	utils.Encoder(writer, userId)
}

func Login(writer http.ResponseWriter, request *http.Request) {
	var Cred model.LoginUser
	err := json.NewDecoder(request.Body).Decode(&Cred)
	utils.CheckError(err)
	loginUser, loginErr := helper.Login(Cred.Email, Cred.Password)
	utils.CheckError(loginErr)
	token, tokenErr := CreateToken(loginUser)
	utils.CheckError(tokenErr)
	_, sessionErr := helper.CreateSession(token, loginUser)
	utils.CheckError(sessionErr)
	utils.Encoder(writer, token)
	utils.Encoder(writer, loginUser)

}

func CreateTask(writer http.ResponseWriter, request *http.Request) {
	var task model.TodoTask
	err := json.NewDecoder(request.Body).Decode(&task)
	utils.CheckError(err)
	userID := GetUserId(request)
	todo, todoErr := helper.CreateTodo(task.Des, userID, task.Task, task.Date)
	utils.CheckError(todoErr)
	utils.Encoder(writer, todo)
}

func AllTask(writer http.ResponseWriter, request *http.Request) {

	userID := GetUserId(request)
	todo, taskErr := helper.TodoList(userID)
	utils.CheckError(taskErr)
	utils.Encoder(writer, todo)
}
func UpcomingTodo(writer http.ResponseWriter, request *http.Request) {

	date := time.Now()
	userID := GetUserId(request)
	todo, todoErr := helper.UpcomingTodoList(userID, date)
	utils.CheckError(todoErr)
	utils.Encoder(writer, todo)
}

func ExpiredTodo(writer http.ResponseWriter, request *http.Request) {

	userID := GetUserId(request)
	date := time.Now()
	todo, err := helper.ExpiredTodo(userID, date)
	utils.CheckError(err)
	utils.Encoder(writer, todo)
	fmt.Println()

}

func UpdateTask(writer http.ResponseWriter, request *http.Request) {
	var task model.TodoTask
	err := json.NewDecoder(request.Body).Decode(&task)
	utils.CheckError(err)
	userID := GetUserId(request)
	newErr := helper.UpdateTask(userID, task.Des, task.Date, task.ID)
	utils.CheckError(newErr)
	utils.Encoder(writer, task)
}

func DeleteTodo(writer http.ResponseWriter, request *http.Request) {
	var taskID int
	err := json.NewDecoder(request.Body).Decode(&taskID)
	utils.CheckError(err)
	userID := GetUserId(request)
	newErr := helper.Deletetask(userID, taskID)
	utils.CheckError(newErr)
	//writer.Write([]byte(fmt.Sprintf("%s task deleted", task)))
	utils.Encoder(writer, taskID)
}

func Logout(writer http.ResponseWriter, request *http.Request) {

	//writer.Write([]byte(fmt.Sprintf("%s USER LOGGED OUT SUCCESSFULLY", user)))

	//userid := request.Header.Get("userid")
	//var token *jwt.Token
	//claims := token.Claims.(jwt.MapClaims)
	//claims["authorized"] = true
	//claims["user"] = userid
	//claims["exp"] = time.Now().Unix()
	//_, err := token.SignedString(mySigningKey)
	//if err != nil {
	//	fmt.Errorf("Something went wrong %s", err.Error())
	//	return
	//}
	//return

}

//ResetPassword helps user to reset password
func ResetPassword(writer http.ResponseWriter, request *http.Request) {
	var cred model.FogetPassword
	err := json.NewDecoder(request.Body).Decode(&cred)
	utils.CheckError(err)
	NewErr := helper.ForgetPass(cred.Email, cred.Userid, cred.MobileNo, cred.Password)
	utils.CheckError(NewErr)
}

//DeleteUser error because of foreign key constraint
func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	var credentials model.LoginUser
	err := json.NewDecoder(request.Body).Decode(&credentials)
	utils.CheckError(err)
	userID := GetUserId(request)
	newErr := helper.DeleteAccount(credentials.Email, credentials.Password)
	utils.CheckError(newErr)
	//writer.Write([]byte(fmt.Sprintf("%s user deleted", userid)))
	utils.Encoder(writer, userID)
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
