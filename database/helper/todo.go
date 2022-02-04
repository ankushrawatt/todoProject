package helper

import (
	"database/sql"
	"time"
	"todoproject/database"
	"todoproject/model"
	"todoproject/utils"
)

//CreateTodo takes description,userid,task, and date and make new todo
func CreateTodo(des, userid, task, date string) (string, error) {
	s, newErr := time.Parse("01-02-2006", date) //date format is MM-DD-YYYY
	if newErr != nil {
		return "", newErr
	}
	SQL := `INSERT INTO todo(userid,des,date,task)VALUES($1,$2,$3,$4) returning task`
	var todo string
	err := database.Todo.Get(&todo, SQL, userid, des, s, task)
	if err != nil {
		return "", err
	}
	return todo, nil
}

//TodoList prints all the todo list of user
func TodoList(userid interface{}) ([]model.TodoTask, error) {
	SQL := `SELECT id,task,des,date from todo WHERE userid=$1`
	user := make([]model.TodoTask, 0)
	err := database.Todo.Select(&user, SQL, userid)
	utils.CheckError(err)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, nil
}

//UpcomingTodoList shows all the upcoming todo of user by comparing today date with the all the todo dates
func UpcomingTodoList(userid interface{}, date time.Time) ([]model.TodoTask, error) {
	SQL := `SELECT id,task,des,date FROM todo WHERE userid=$1 AND date>$2`
	user := make([]model.TodoTask, 0)
	err := database.Todo.Select(&user, SQL, userid, date)
	utils.CheckError(err)
	return user, nil
}

//ExpiredTodo shows all the expired todo of user by comparing today date with the all the todo dates
func ExpiredTodo(userid interface{}, date time.Time) ([]model.TodoTask, error) {
	SQL := `SELECT id,task,des, date FROM todo WHERE userid=$1 AND date<$2::DATE`
	user := make([]model.TodoTask, 0)
	err := database.Todo.Select(&user, SQL, userid, date)
	utils.CheckError(err)
	return user, nil
}

//Deletetask takes userid and id and delete the todo of user by compairing userid and id of todo
func Deletetask(userid string, id int) error {
	SQL := `DELETE FROM todo WHERE id=$1 AND userid=$2`
	_, err := database.Todo.Exec(SQL, id, userid)
	utils.CheckError(err)
	return nil
}

//UpdateTask will update the given task
func UpdateTask(userid, des, date string, id int) error {
	SQL := `UPDATE todo SET des=$1,date=$2 WHERE userid=$3 AND id=$4`
	newDate, err := time.Parse("01-02-2006", date) //date format is MM-DD-YYYY
	utils.CheckError(err)
	_, newErr := database.Todo.Exec(SQL, des, newDate, userid, id)
	utils.CheckError(newErr)
	return nil
}
