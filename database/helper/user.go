package helper

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
	"todoproject/database"
	"todoproject/model"
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

func GetUserSession(token string) (string, error) {
	SQL := `SELECT userid from session where token=$1`
	var user string
	err := database.Todo.Get(&user, SQL, token)
	if err != nil {
		return "", err
	}
	return user, nil
}

func CreateSession(id, userid string) error {
	SQL := `INSERT INTO session(userid,token)VALUES($1,$2) RETURNING TOKEN`
	var token string
	err := database.Todo.Get(&token, SQL, userid, id)
	if err != nil {
		return err
	}
	return nil

}

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
	SQL := `SELECT task,des,date from todo WHERE userid=$1`
	user := make([]model.TodoTask, 0)
	err := database.Todo.Select(&user, SQL, userid)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, nil
}

func UpcomingTodoList(userid string, date time.Time) ([]model.TodoTask, error) {
	SQL := `SELECT task,des, date FROM todo WHERE userid=$1 AND date>$1`
	user := make([]model.TodoTask, 0)
	err := database.Todo.Select(&user, SQL, userid, date)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func ExpiredTodo(userid string, date time.Time) ([]model.TodoTask, error) {
	SQL := `SELECT task,des, date FROM todo WHERE userid=$1 AND date=<$2`
	user := make([]model.TodoTask, 0)
	err := database.Todo.Select(&user, SQL, userid, date)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func DeleteSession(token string) error {
	SQL := `DELETE FROM session WHERE token=$1`
	_, err := database.Todo.Exec(SQL, token)
	if err != nil {
		return err
	}
	return nil
}
