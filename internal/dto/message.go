package dto

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

type MessageType uint8

const (
	Key MessageType = iota
	ServerMessage
)

type Message struct {
	Type MessageType
	Body string
}

func NewMessageFromString(s string) (*Message, error) {
	log.Println("DEBUG: message => ", s)

	messageParts := strings.Split(s, "::")
	if len(messageParts) != 2 {
		return nil, errors.New("incorrect message format")
	}

	log.Println("DEBUG: ", messageParts[0])
	t, err := strconv.Atoi(messageParts[0])
	if err != nil {
		return nil, errors.New("incorrect message type")
	}

	return &Message{
		Type: MessageType(t),
		Body: messageParts[0],
	}, nil
}

func NewMessage(t MessageType, b string) *Message {
	return &Message{
		Type: t,
		Body: b,
	}
}

func (m *Message) String() string {
	log.Println("DEBUG: tostring => ", m.Type)
	// FIXME: Problem here with type
	return string(m.Type) + "::" + m.Body
}

func (m *Message) Bytes() []byte {
	return []byte(m.String())
}
