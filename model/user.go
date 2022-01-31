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

type CreateTask struct {
	Description string `db:"description" json:"description"`
	Date        string `db:"date" json:"date"`
	UserId      string `db:"userid" json:"userId"`
	TaskName    string `db:"taskname" json:"taskName"`
}
