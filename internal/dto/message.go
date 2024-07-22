package dto

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type MessageType uint8

const (
	ChatMessage            MessageType = iota //Structure: type::body
	ServerMessage                             //Structure: type::body
	AskSessionIdMessage                       //Structure: type::body
	AskKeyMessage                             //Structure: type::body
	KeyExchangeMessage                        //Structure: type::body
	UserKeyExchangeMessage                    //Structure: type::body
	ClientMessage                             //Structure: type::body
	JoinSessionMessage                        //Structure: type::sessionId
)

var (
	LastType       MessageType = JoinSessionMessage
	MinMessageLen  int         = 2
	BaseMessageLen int         = 2
	JoinMessageLen int         = 2
	ChatMessageLen int         = 3
	MaxMessageLen  int         = 3
	TypePos        int         = 0
	BodyPos        int         = 1
	PubKeyPos      int         = 2
	SessionIdPos   int         = 1
)

var (
	ErrMessageLength      error = errors.New("incorrect message length")
	ErrMessageType        error = errors.New("incorrect message type")
	ErrMessageSessionId   error = errors.New("incorrect session ID format")
	ErrExpectedKeyExhange error = errors.New("message type should be of key exhange")
)

type Message struct {
	Type      MessageType
	Body      string
	SessionId uuid.UUID
}

func NewMessage(t MessageType, b string) *Message {
	return &Message{
		Type: t,
		Body: b,
	}
}

func NewJoinSessionMessage(id uuid.UUID) *Message {
	return &Message{
		Type:      JoinSessionMessage,
		SessionId: id,
	}
}

func NewMessageFromString(s string) (*Message, error) {
	messageParts := strings.Split(s, "::")
	l := len(messageParts)
	if l < MinMessageLen || l > MaxMessageLen {
		return nil, ErrMessageLength
	}

	t, err := strconv.Atoi(messageParts[TypePos])
	if err != nil {
		return nil, ErrMessageType
	}

	if t > int(LastType) || t < 0 {
		return nil, ErrMessageType
	}

	// JOIN SESSION MESSAGE
	if t == int(JoinSessionMessage) {
		if l != JoinMessageLen {
			return nil, ErrMessageType
		}

		id, err := uuid.Parse(messageParts[SessionIdPos])
		if err != nil {
			log.Println("Error, could not parse uuid: ", err.Error())
			return nil, ErrMessageSessionId
		}

		return NewJoinSessionMessage(
			id,
		), nil
	}

	// BASE MESSAGE TYPES
	return NewMessage(MessageType(t), messageParts[1]), nil
}
func (m *Message) String() string {
	if m.Type == JoinSessionMessage {
		return fmt.Sprintf("%d::%s", m.Type, m.SessionId.String())
	}

	return fmt.Sprintf("%d::%s", m.Type, m.Body)
}

func (m *Message) Bytes() []byte {
	return []byte(m.String())
}
