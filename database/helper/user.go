package helper

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"todoproject/database"
	"todoproject/utils"
)

//HashPassword helps to encrypt the password
func HashPassword(password string) string {
	pass, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(pass)
}

//CheckHashPassword  verify the password of user
func CheckHashPassword(password, hashpass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

//NewUser save the data of new user
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

//Login function takes email and password, returns userID
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

//ForgetPass helps user to reset password by taking userid, email and mobile no. for
//security purpose and match them with the users table and if everything matches it set new password
func ForgetPass(email, userid, mobile, password string) error {
	SQL := `SELECT userid, mobile FROM users WHERE email=$1`
	var id, mobileNo string
	err := database.Todo.QueryRowx(SQL, email).Scan(&id, &mobileNo)
	utils.CheckError(err)
	if mobileNo == mobile && userid == id {
		newSQL := `UPDATE users SET password=$1 where userid=$2`
		pass := HashPassword(password)
		_, newErr := database.Todo.Exec(newSQL, pass, userid)
		utils.CheckError(newErr)
	} else {
		return errors.New("WRONG CREDENTIALS")
	}

	return nil
}

func DeleteAccount(email, password string) error {
	userID, err := Login(email, password)
	utils.CheckError(err)
	SQL := `DELETE FROM users WHERE userid=$1`
	_, newErr := database.Todo.Exec(SQL, userID)
	utils.CheckError(newErr)
	return nil
}
