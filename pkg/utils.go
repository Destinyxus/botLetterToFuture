package pkg

import "github.com/Destinyxus/botLetterToFuture/internal/model"

func UpdateStruct(m *model.User) {
	m.Letter = ""
	m.Date = ""
	m.Email = ""
}
