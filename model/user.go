package model

type User struct {
	ID        string `db:"id" json:"id"`
	FirstName string `db:"firstname" json:"firstName"`
	LastName  string `db:"lastname" json:"lastName"`
	UserId    string `db:"userid" json:"userId"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"password"`
	MobileNo  string `db:"mobile" json:"mobileNo"`
	//CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type LoginUser struct {
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type TodoTask struct {
	Task string `db:"task" json:"task"`
	Des  string `db:"des" json:"des"`
	Date string `db:"date" json:"date"`
}

type FogetPassword struct {
	Email    string `db:"email" json:"email"`
	Userid   string `db:"userid" json:"userid"`
	Password string `db:"password" json:"password"`
	MobileNo string `db:"mobile" json:"mobileNo"`
}
