package dto_test

import (
	"SecureChat/internal/dto"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// Message structure: type::body

func TestMessage(t *testing.T) {
	//arrange
	message := dto.NewMessage(dto.ServerMessage, "hello")

	//assert
	if message.Type != dto.ServerMessage || message.Body != "hello" {
		t.Errorf("Contructor failed test")
	}
}

func TestJoinSessionMessage(t *testing.T) {
	//arrange
	expectedType := dto.JoinSessionMessage
	expectedId, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	message := dto.NewJoinSessionMessage(expectedId)

	//assert
	if message.Type != expectedType {
		t.Errorf("wanted: Type => %d, got: %d", expectedType, message.Type)
	}

	if message.SessionId != expectedId {
		t.Errorf("wanted: SessionId => %s, got: %s", expectedId.String(), message.SessionId.String())
	}
}

func TestServerMessageToString(t *testing.T) {
	// arrange
	message := dto.NewMessage(dto.ServerMessage, "hello")
	expected := "1::hello"

	// assert
	if message.String() != expected {
		t.Errorf("wanted: %s, got: %s", expected, message.String())
	}
}

func TestServerMessageFromString(t *testing.T) {
	// arrange
	s := fmt.Sprintf("%d::%s", dto.ServerMessage, "hello from goland")

	// act && assert
	message, err := dto.NewMessageFromString(s)
	if err != nil {
		t.Errorf("Failed parsing string: %s", err.Error())
	}

	if message.Type != dto.ServerMessage || message.Body != "hello from goland" {
		t.Errorf(
			"wanted: type=%d, body=%s; got: type=%d, body=%s;",
			dto.ServerMessage,
			"hello from goland",
			message.Type,
			message.Body,
		)
	}
}

func TestJoinSessionMessageToString(t *testing.T) {
	// arrange
	expectedId, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	message := dto.NewJoinSessionMessage(expectedId)
	expected := fmt.Sprintf("%d::%s", dto.JoinSessionMessage, expectedId.String())

	// assert
	if message.String() != expected {
		t.Errorf("wanted: %s, got: %s", expected, message.String())
	}
}

func TestJoinSessionMessageFromString(t *testing.T) {
	// arrange
	expectedId, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	expectedType := dto.JoinSessionMessage
	s := fmt.Sprintf("%d::%s", expectedType, expectedId.String())

	// act && assert
	message, err := dto.NewMessageFromString(s)
	if err != nil {
		t.Errorf("Failed parsing string: %s", err.Error())
	}

	if message.Type != expectedType ||
		message.SessionId != expectedId {
		t.Errorf(
			"wanted: type=%d, key=%s; got: type=%d, key=%s;",
			expectedType,
			expectedId.String(),
			message.Type,
			message.SessionId.String(),
		)
	}
}
