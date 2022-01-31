package helper

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"todoproject/database"
	"todoproject/model"
)

func HashPassword(password string) string {
	pass, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(pass)
}

func CheckHashPassword(password, hashpass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hashpass))
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
	return "", nil
}

func Login(email, password string) (string, error) {
	SQL := `SELECT userid,password FROM users WHERE email=$1`
	var user, Hashpassword string
	err := database.Todo.QueryRowx(SQL, email).Scan(&user, &Hashpassword)
	if err != nil {
		return "", err
	}
	pass, passErr := CheckHashPassword(password, Hashpassword)
	if pass != true && passErr != nil {
		return "", passErr
	}
	return user, nil

}

func CreateTodo(des, userid, date, task string) (string, error) {
	SQL := `INSERT INTO todo(userid,des,date,task)VALUES($1,$2,$3,$4) returning task`
	var todo string
	err := database.Todo.Get(&todo, SQL, userid, des, date, task)
	if err != nil {
		return "", err
	}
	return todo, nil
}

func TodoList(userid string) ([]model.CreateTask, error) {
	SQL := `SELECT task,des,date from todo`
	user := make([]model.CreateTask, 0)
	err := database.Todo.Select(&user, SQL)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, nil
}

func TodoByDate(date string) ([]model.CreateTask, error) {
	SQL := `SELECT task,des, date FROM todo WHERE date=$1`
	user := make([]model.CreateTask, 0)
	err := database.Todo.Select(&user, SQL, date)
	if err != nil {
		return nil, err
	}
	return user, nil
}
