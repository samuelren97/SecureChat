package models

import (
	"SecureChat/internal/datastructures"
	"SecureChat/internal/dto"
	"errors"

	"github.com/google/uuid"
)

const (
	MaxNbUsers int = 2
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
	if sm.UserCount() < MaxNbUsers {
		sm.users.Add(u)
		return nil
	}

	return ErrSessionFull
}

func (sm *SessionModel) UserCount() int {
	return sm.users.Count
}

func (sm *SessionModel) Close() {
	sm.users.ForEach(func(u *User) {
		u.Conn.Close()
	})
}

func (sm *SessionModel) SendKeys() {
	user1 := sm.users.Get(0)
	user2 := sm.users.Get(1)

	user1Message := dto.NewMessage(dto.UserKeyExchangeMessage, user2.PubKey)
	user2Message := dto.NewMessage(dto.UserKeyExchangeMessage, user1.PubKey)

	user1.Conn.Write(user1Message.Bytes())
	user2.Conn.Write(user2Message.Bytes())
}
