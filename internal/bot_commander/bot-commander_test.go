package bot_commander_test

import (
	"testing"
	"time"

	"github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	"github.com/Destinyxus/botLetterToFuture/internal/bot_commander/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBotCommander_CheckForActualDate(t *testing.T) {
	dateExample := time.Now()

	testTable := []struct {
		name                  string
		configureDateIndex    func() map[time.Time]struct{}
		configureMockBehavior func(repo *mocks.MockRepository, sndr *mocks.MockEmailSender)
		expectedError         error
	}{
		{
			name: "everything ok",
			configureDateIndex: func() map[time.Time]struct{} {
				m := make(map[time.Time]struct{})
				m[dateExample] = struct{}{}

				return m
			},
			configureMockBehavior: func(
				repo *mocks.MockRepository,
				sndr *mocks.MockEmailSender,
			) {
				repo.EXPECT().GetLetter(dateExample).Return(
					[]bot_commander.Letter{
						{
							Date:   dateExample,
							Email:  "foo@gmail.com",
							Letter: "hii",
						},
					},
					nil,
				).Times(1)

				sndr.EXPECT().SendEmail("foo@gmail.com", "hii").Return(nil).Times(1)
			},
			expectedError: nil,
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
			}

			a.Equal(b.CheckForActualDate(), test.expectedError)
		})
	}
}
