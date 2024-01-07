package bot_commander

import (
	"reflect"
	"testing"
	"time"
)

func FormatTime(t time.Time) time.Time {
	now := t.Format(DateFormat)

	parsed, _ := time.Parse(DateFormat, now)

	return parsed
}

func TestValidateMessage(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		want      Letter
		resultErr error
	}{
		{
			name: "all good",
			message: "aboba. hi, my name is aboba;aboba@gmail.ru;" +
				time.Now().Format(DateFormat),
			want: Letter{
				Letter: "aboba. hi, my name is aboba",
				Email:  "aboba@gmail.ru",
				Date:   FormatTime(time.Now()),
			},
			resultErr: nil,
		},
		{
			name:      "invalid input: 2 `;` required",
			message:   "aboba. hi, my name is aboba;wef;reg;erg;aboba@gmail.ru;2025-01-01",
			want:      Letter{},
			resultErr: ErrNotValidEmailOrDate,
		},
		{
			name:      "invalid email",
			message:   "aboba. hi, my name is aboba;abmail.ru;2025-01-01",
			want:      Letter{},
			resultErr: ErrNotValidEmailOrDate,
		},
		{
			name:      "invalid date",
			message:   "aboba. hi, my name is aboba;aboba@gmail.ru;2029-01",
			want:      Letter{},
			resultErr: ErrNotValidEmailOrDate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateMessage(tt.message)
			if err != tt.resultErr {
				t.Errorf("ValidateMessage() error = %v, resultErr %v", err, tt.resultErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
