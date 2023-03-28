package model

type Model struct {
	Id     int
	Email  string
	Date   string
	Letter string
	Sent   bool
}

type DeleteModel struct {
	MessageId int
}

func NewDeleteModel() *DeleteModel {
	return &DeleteModel{}

}

func NewModel(email string, date string, letter string) *Model {
	return &Model{Email: email, Date: date, Letter: letter}
}

func TemporaryModel() *Model {
	return &Model{}
}
