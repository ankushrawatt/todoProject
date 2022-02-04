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

//HashPassword helps to encrypt the password
func HashPassword(password string) string {
	pass, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(pass)
}

//CheckHashPassword it verify the password of user
func CheckHashPassword(password, hashpass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

//NewUser safe the data of new user
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

//GetUser will give you userid of user with token in header
func GetUser(token string) (string, error) {
	SQL := `SELECT userid from session where token=$1`
	var user string
	err := database.Todo.Get(&user, SQL, token)
	if err != nil {
		return "", err
	}
	return user, nil
}

//CreateSession makes record of user into session table
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

//ForgetPass helps user to reset password by taking userid, email and mobile no. for
//security purpose and match them with the users table and if everything matches it set new password
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

func DeleteAccount(email, password string) error {
	user, err := Login(email, password)
	utils.CheckError(err)
	SQL := `DELETE FROM users WHERE userid=$1`
	_, newErr := database.Todo.Exec(SQL, user)
	utils.CheckError(newErr)
	return nil
}
