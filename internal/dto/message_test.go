package dto_test

import (
	"SecureChat/internal/dto"
	"fmt"
	"testing"
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

func TestMessageToString(t *testing.T) {
	// arrange
	message := dto.NewMessage(dto.Key, "hello")
	expected := "0::hello"

	// assert
	if message.String() != expected {
		t.Errorf("wanted: %s, got: %s", expected, message.String())
	}
}

func TestMessageFromString(t *testing.T) {
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
