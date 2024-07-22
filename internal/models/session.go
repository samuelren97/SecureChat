package models

import (
	"SecureChat/internal/datastructures"

	"github.com/google/uuid"
)

type SessionModel struct {
	Id    uuid.UUID
	Users *datastructures.LinkList[*User]
}

func NewSession() *SessionModel {
	return &SessionModel{
		Id:    uuid.New(),
		Users: datastructures.NewLinkList[*User](),
	}
}
