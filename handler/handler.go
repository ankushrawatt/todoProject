package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"todoproject/database/helper"
	"todoproject/model"
)

func CreateToken(user string) (string, error) {
	var err error
	os.Setenv("ACCESS_KEY", "")

	return "", err
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
		writer.Write([]byte(fmt.Sprintf("WRONG CREDENTIOALS")))
		return
	}
	jsonData, jsonErr := json.Marshal(loginUser)
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
	writer.Write(jsonData)
}

func CreateTask(writer http.ResponseWriter, request *http.Request) {
	var task model.CreateTask
	err := json.NewDecoder(request.Body).Decode(&task)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	todo, todoErr := helper.CreateTodo(task.Description, task.UserId, task.Date, task.TaskName)
	if todoErr != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	writer.Write([]byte(fmt.Sprintf("New task created %s", todo)))
}
func UpdateTask(writer http.ResponseWriter, request *http.Request) {

}

func AllTask(writer http.ResponseWriter, request *http.Request) {
	var todo model.CreateTask
	err := json.NewDecoder(request.Body).Decode(&todo)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	tasks, taskErr := helper.TodoList(todo.UserId)
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

func SearchByDate(writer http.ResponseWriter, request *http.Request) {
	var date string
	err := json.NewDecoder(request.Body).Decode(&date)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	todo, todoErr := helper.TodoByDate(date)
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

func Logout(writer http.ResponseWriter, request *http.Request) {

}
