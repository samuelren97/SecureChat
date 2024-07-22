package models

import (
	"net"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	Conn      net.Conn
	SessionId uuid.UUID
}

func NewUser(conn net.Conn, sessionId uuid.UUID) *User {
	return &User{
		Id:        uuid.New(),
		Conn:      conn,
		SessionId: sessionId,
	}
}
