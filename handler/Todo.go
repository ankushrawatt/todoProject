package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"todoproject/database/helper"
	"todoproject/model"
	"todoproject/utils"
)

func CreateTodo(writer http.ResponseWriter, request *http.Request) {
	var task model.TodoTask
	err := json.NewDecoder(request.Body).Decode(&task)
	utils.CheckError(err)
	userID := GetUserId(request)
	todo, todoErr := helper.CreateTodo(task.Des, userID, task.Task, task.Date)
	utils.CheckError(todoErr)
	utils.Encoder(writer, todo)
}

func AllTodo(writer http.ResponseWriter, request *http.Request) {

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

func UpdateTodo(writer http.ResponseWriter, request *http.Request) {
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
