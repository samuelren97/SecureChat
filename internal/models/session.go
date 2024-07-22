package models

import (
	"SecureChat/internal/datastructures"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrSessionFull error = errors.New("session room is full")
)

type SessionModel struct {
	Id    uuid.UUID
	users *datastructures.LinkList[*User]
}

func NewSession() *SessionModel {
	return &SessionModel{
		Id:    uuid.New(),
		users: datastructures.NewLinkList[*User](),
	}
}

func (sm *SessionModel) AddUser(u *User) error {
	if sm.UserCount() < 2 {
		sm.users.Add(u)
		return nil
	}

	return ErrSessionFull
}

func (sm *SessionModel) UserCount() int {
	return sm.users.Count
}
