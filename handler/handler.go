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

//session
//func CreateToken(userid string) (string, error) {

//id := uuid.New()
//err := helper.CreateSession(id.String(), userid)
//if err != nil {
//	return "", err
//}
//
//return id.String(), nil

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
//}

func Signup(writer http.ResponseWriter, request *http.Request) {
	var user model.User
	err := json.NewDecoder(request.Body).Decode(&user)
	utils.CheckError(err)
	id := uuid.New()
	userId, NewErr := helper.NewUser(id.String(), user.Email, user.UserId, user.Password, user.FirstName, user.LastName, user.MobileNo)
	utils.CheckError(NewErr)
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
	utils.CheckError(err)
	loginUser, loginErr := helper.Login(Cred.Email, Cred.Password)
	utils.CheckError(loginErr)

	_, jsonErr := json.Marshal(loginUser)
	if jsonErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	token, tokenErr := CreateToken(loginUser)
	utils.CheckError(tokenErr)
	_, sessionErr := helper.CreateSession(token, loginUser)
	utils.CheckError(sessionErr)
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
	utils.CheckError(err)
	userkey := request.Header.Get("userid")
	todo, todoErr := helper.CreateTodo(task.Des, userkey, task.Task, task.Date)
	utils.CheckError(todoErr)
	writer.Write([]byte(fmt.Sprintf("New task created %s", todo)))
}
func UpdateTask(writer http.ResponseWriter, request *http.Request) {
	var task model.TodoTask
	err := json.NewDecoder(request.Body).Decode(&task)
	utils.CheckError(err)
	userid := request.Header.Get("userid")
	newErr := helper.UpdateTask(task.Task, userid, task.Des, task.Date)
	utils.CheckError(newErr)
	writer.Write([]byte(fmt.Sprintf("%s task Updated successfulluy", task)))
}

func AllTask(writer http.ResponseWriter, request *http.Request) {

	userkey := request.Header.Get("userid")
	tasks, taskErr := helper.TodoList(userkey)
	utils.CheckError(taskErr)
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
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(jsonData)
}

func ExpiredTodo(writer http.ResponseWriter, request *http.Request) {
	date := time.Now()
	userKey := request.Header.Get("userid")
	todo, err := helper.ExpiredTodo(userKey, date)
	utils.CheckError(err)
	//jsonData, jsonErr := json.Marshal(todo)
	//if jsonErr != nil {
	//	writer.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	err = json.NewEncoder(writer).Encode(todo)
	utils.CheckError(err)
	//writer.WriteHeader(http.StatusCreated)
	//writer.Write(jsonData)
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

func ResetPassword(writer http.ResponseWriter, request *http.Request) {
	var cred model.FogetPassword
	err := json.NewDecoder(request.Body).Decode(&cred)
	utils.CheckError(err)
	NewErr := helper.ForgetPass(cred.Email, cred.Userid, cred.MobileNo, cred.Password)
	utils.CheckError(NewErr)
}

func DeleteTodo(writer http.ResponseWriter, request *http.Request) {
	var task model.DeleteTask
	err := json.NewDecoder(request.Body).Decode(&task)
	utils.CheckError(err)
	userid := request.Header.Get("userid")
	newErr := helper.Deletetask(userid, task.Task)
	utils.CheckError(newErr)
	writer.Write([]byte(fmt.Sprintf("%s task deleted", task)))
}

//Delete user error because of foreign key constraint

func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	var credentials model.LoginUser
	err := json.NewDecoder(request.Body).Decode(&credentials)
	utils.CheckError(err)
	userid := request.Header.Get("userid")
	newErr := helper.DeleteAccount(credentials.Email, credentials.Password, userid)
	utils.CheckError(newErr)
	writer.Write([]byte(fmt.Sprintf("%s user deleted", userid)))
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
