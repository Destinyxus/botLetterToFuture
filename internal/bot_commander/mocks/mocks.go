// Code generated by MockGen. DO NOT EDIT.
// Source: internal/bot_commander/bot-commander.go

// Package mock_bot_commander is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	bot_commander "github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	gomock "github.com/golang/mock/gomock"
)

// MockEmailSender is a mock of EmailSender interface.
type MockEmailSender struct {
	ctrl     *gomock.Controller
	recorder *MockEmailSenderMockRecorder
}

// MockEmailSenderMockRecorder is the mock recorder for MockEmailSender.
type MockEmailSenderMockRecorder struct {
	mock *MockEmailSender
}

// NewMockEmailSender creates a new mock instance.
func NewMockEmailSender(ctrl *gomock.Controller) *MockEmailSender {
	mock := &MockEmailSender{ctrl: ctrl}
	mock.recorder = &MockEmailSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailSender) EXPECT() *MockEmailSenderMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockEmailSender) SendEmail(email, letter string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", email, letter)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockEmailSenderMockRecorder) SendEmail(email, letter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockEmailSender)(nil).SendEmail), email, letter)
}

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetActualDates mocks base method.
func (m *MockRepository) GetActualDates() (map[time.Time]struct{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActualDates")
	ret0, _ := ret[0].(map[time.Time]struct{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActualDates indicates an expected call of GetActualDates.
func (mr *MockRepositoryMockRecorder) GetActualDates() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActualDates", reflect.TypeOf((*MockRepository)(nil).GetActualDates))
}

// GetLetter mocks base method.
func (m *MockRepository) GetLetter(date time.Time) ([]bot_commander.Letter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLetter", date)
	ret0, _ := ret[0].([]bot_commander.Letter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLetter indicates an expected call of GetLetter.
func (mr *MockRepositoryMockRecorder) GetLetter(date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLetter", reflect.TypeOf((*MockRepository)(nil).GetLetter), date)
}

// InsertLetter mocks base method.
func (m *MockRepository) InsertLetter(letter, email string, date time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertLetter", letter, email, date)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertLetter indicates an expected call of InsertLetter.
func (mr *MockRepositoryMockRecorder) InsertLetter(letter, email, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertLetter", reflect.TypeOf((*MockRepository)(nil).InsertLetter), letter, email, date)
}