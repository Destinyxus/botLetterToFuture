package model

type User struct {
	UserName        string
	Email           string
	Date            string
	Letter          string
	EncryptedLetter string
	Sent            bool
}

func NewConfigUser() *User {
	return &User{}
}

type DeleteUser struct {
	LetterId int
	DateId   int
	EmailId  int
}

func NewDeleteModel() *DeleteUser {
	return &DeleteUser{}

}

func NewUser(email string, date string, letter string) *User {
	return &User{Email: email, Date: date, EncryptedLetter: letter}
}

func TemporaryUser() *User {
	return &User{}
}
