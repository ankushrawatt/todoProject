package helper

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
	"todoproject/database"
	"todoproject/model"
	"todoproject/utils"
)

func HashPassword(password string) string {
	pass, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(pass)
}

func CheckHashPassword(password, hashpass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewUser(id, email, userid, password, firstname, lastname, mobile string) (string, error) {
	SQL := `INSERT INTO users(id,email,userid,password,firstname,lastname,mobile)VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING userid`
	var user string
	pass := HashPassword(password)
	err := database.Todo.Get(&user, SQL, id, email, userid, pass, firstname, lastname, mobile)
	if err != nil {
		return "", err
	}
	return user, nil

}

//SESSION TABLE
//func GetUserSession(token string) (string, error) {
//	SQL := `SELECT userid from session where token=$1`
//	var user string
//	err := database.Todo.Get(&user, SQL, token)
//	if err != nil {
//		return "", err
//	}
//	return user, nil
//}

//SESSION table
func CreateSession(id, userid string) (string, error) {
	SQL := `INSERT INTO session(userid,token)VALUES($1,$2) RETURNING TOKEN`
	var token string
	err := database.Todo.Get(&token, SQL, userid, id)
	if err != nil {
		return "", err
	}
	return token, nil

}

// Login function takes email and password, returns userID
func Login(email, password string) (string, error) {
	SQL := `SELECT userid,password FROM users WHERE email=$1`
	var user, Hashpass string
	err := database.Todo.QueryRowx(SQL, email).Scan(&user, &Hashpass)
	if err != nil {
		return "", err
	}
	pass, passErr := CheckHashPassword(password, Hashpass)
	if pass != true && passErr != nil {
		return "", passErr
	}
	return user, nil

}

func CreateTodo(des, userid, task, date string) (string, error) {
	s, newErr := time.Parse("01-02-2006", date)
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

func TodoList(userid string) ([]model.TodoTask, error) {
	SQL := `SELECT id,task,des,date from todo WHERE userid=$1`
	user := make([]model.TodoTask, 0)
	err := database.Todo.Select(&user, SQL, userid)
	utils.CheckError(err)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, nil
}

func UpcomingTodoList(userid string, date time.Time) ([]model.TodoTask, error) {
	SQL := `SELECT id,task,des, date FROM todo WHERE userid=$1 AND date>$2`
	user := make([]model.TodoTask, 0)
	err := database.Todo.Select(&user, SQL, userid, date)
	utils.CheckError(err)
	return user, nil
}

func ExpiredTodo(userid string, date time.Time) ([]model.TodoTask, error) {
	SQL := `SELECT id,task,des, date FROM todo WHERE userid=$1 AND date<$2::DATE`
	user := make([]model.TodoTask, 0)
	err := database.Todo.Select(&user, SQL, userid, date)
	utils.CheckError(err)
	return user, nil
}

//SESSION TABLE
func DeleteSession(token string) error {
	SQL := `DELETE FROM session WHERE token=$1`
	_, err := database.Todo.Exec(SQL, token)
	utils.CheckError(err)
	return nil
}

func ForgetPass(email, userid, mobile, password string) error {
	SQL := `SELECT userid, mobile FROM users WHERE email=$1`
	var id, mob string
	err := database.Todo.QueryRowx(SQL, email).Scan(&id, &mob)
	utils.CheckError(err)
	if mob == mobile && userid == id {
		newSQL := `UPDATE users SET password=$1 where userid=$2`
		pass := HashPassword(password)
		_, newErr := database.Todo.Exec(newSQL, pass, userid)
		utils.CheckError(newErr)
	} else {
		return errors.New("WRONG CREDENTIALS")
	}

	return nil
}

func Deletetask(userid, task string) error {
	SQL := `DELETE FROM todo WHERE task=$1 AND userid=$2`
	_, err := database.Todo.Exec(SQL, task, userid)
	utils.CheckError(err)
	return nil
}

func UpdateTask(task, userid, des, date string) error {
	SQL := `UPDATE todo SET des=$1,date=$2 WHERE userid=$3 AND task=$4`
	newDate, err := time.Parse("01-02-2006", date)
	utils.CheckError(err)
	_, newErr := database.Todo.Exec(SQL, des, newDate, userid, task)
	utils.CheckError(newErr)
	return nil
}

func DeleteAccount(email, password, userid string) error {
	user, err := Login(email, password)
	utils.CheckError(err)
	SQL := `DELETE FROM users WHERE userid=$1`
	_, newErr := database.Todo.Exec(SQL, user)
	utils.CheckError(newErr)
	return nil
}
