package models

import (
	"net"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	Conn      net.Conn
	SessionId uuid.UUID
	PubKey    string
}

func NewUser(conn net.Conn, sessionId uuid.UUID) *User {
	return &User{
		Id:        uuid.New(),
		Conn:      conn,
		SessionId: sessionId,
	}
}

func (u *User) AddPubKey(pubKey string) {
	u.PubKey = pubKey
}
