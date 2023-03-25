package model

type Model struct {
	Id     int
	Email  string
	Date   string
	Letter string
}

func NewModel(email string, date string, letter string) *Model {
	return &Model{Email: email, Date: date, Letter: letter}
}

func TemporaryModel() *Model {
	return &Model{}
}
