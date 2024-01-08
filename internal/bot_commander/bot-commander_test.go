package bot_commander_test

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
	"testing"
	"time"

	"github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	mocks "github.com/Destinyxus/botLetterToFuture/internal/bot_commander/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBotCommander_CheckForActualDate(t *testing.T) {
	currentDateFormatted := bot_commander.FormatTime(time.Now())

	testTable := []struct {
		name                  string
		configureDateIndex    func() map[time.Time]struct{}
		configureMockBehavior func(
			repo *mocks.MockRepository,
			sndr *mocks.MockEmailSender,
		)
		expectedError error
	}{
		{
			name: "everything ok",
			configureDateIndex: func() map[time.Time]struct{} {
				m := make(map[time.Time]struct{})
				m[currentDateFormatted] = struct{}{}

				return m
			},
			configureMockBehavior: func(
				repo *mocks.MockRepository,
				sndr *mocks.MockEmailSender,
			) {
				repo.EXPECT().GetLetter(currentDateFormatted).Return(
					[]bot_commander.Letter{
						{
							Date:   currentDateFormatted,
							Email:  "foo@gmail.com",
							Letter: "hii",
						},
						{
							Date:   currentDateFormatted,
							Email:  "hey@gmail.com",
							Letter: "lol",
						},
					},
					nil,
				).Times(1)

				sndr.EXPECT().SendEmail("foo@gmail.com", "hii").Return(nil).Times(1)
				sndr.EXPECT().SendEmail("hey@gmail.com", "lol").Return(nil).Times(1)

				repo.EXPECT().DeprecateLetter(currentDateFormatted).Return(nil).Times(1)
			},
			expectedError: nil,
		},
		{
			name: "no emails on this date",
			configureDateIndex: func() map[time.Time]struct{} {
				return make(map[time.Time]struct{})
			},
			configureMockBehavior: func(
				repo *mocks.MockRepository,
				sndr *mocks.MockEmailSender,
			) {
				repo.EXPECT().GetLetter("foo").Times(0)

				sndr.EXPECT().SendEmail("foo", "foo").Times(0)
			},
			expectedError: nil,
		},
		{
			name: "get letter error",
			configureDateIndex: func() map[time.Time]struct{} {
				m := make(map[time.Time]struct{})
				m[currentDateFormatted] = struct{}{}

				return m
			},
			configureMockBehavior: func(
				repo *mocks.MockRepository,
				sndr *mocks.MockEmailSender,
			) {
				repo.EXPECT().GetLetter(currentDateFormatted).
					Return([]bot_commander.Letter{}, errors.New("db err")).
					Times(1)

				sndr.EXPECT().SendEmail("foo", "foo").Times(0)
			},
			expectedError: fmt.Errorf("getting the letter with date: %w", errors.New("db err")),
		},
		{
			name: "send email error",
			configureDateIndex: func() map[time.Time]struct{} {
				m := make(map[time.Time]struct{})
				m[currentDateFormatted] = struct{}{}

				return m
			},
			configureMockBehavior: func(
				repo *mocks.MockRepository,
				sndr *mocks.MockEmailSender,
			) {
				repo.EXPECT().GetLetter(currentDateFormatted).Return(
					[]bot_commander.Letter{
						{
							Date:   currentDateFormatted,
							Email:  "foo@gmail.com",
							Letter: "hii",
						},
					},
					nil,
				).Times(1)

				sndr.EXPECT().SendEmail("foo@gmail.com", "hii").
					Return(errors.New("email sender error")).
					Times(1)
			},
			expectedError: fmt.Errorf("sending the letters to emails: %w", errors.New("email sender error")),
		},
	}

	a := assert.New(t)

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepository(ctrl)
			sndr := mocks.NewMockEmailSender(ctrl)

			test.configureMockBehavior(repo, sndr)

			b := bot_commander.BotCommander{
				Repo:        repo,
				EmailSender: sndr,
				DateIndex:   test.configureDateIndex(),
				Logger:      logrus.New(),
			}

			a.Equal(b.CheckForActualDate(), test.expectedError)
		})
	}
}
