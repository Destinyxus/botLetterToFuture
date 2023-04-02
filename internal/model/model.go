package model

type Model struct {
	Id              int
	Email           string
	Date            string
	Letter          string
	EncryptedLetter string
	Sent            bool
}

type DeleteModel struct {
	LetterId int
	DateId   int
	EmailId  int
}

func NewDeleteModel() *DeleteModel {
	return &DeleteModel{}

}

func NewModel(email string, date string, letter string) *Model {
	return &Model{Email: email, Date: date, EncryptedLetter: letter}
}

func TemporaryModel() *Model {
	return &Model{}
}
