package main

import (
	"SecureChat/internal/dto"
	"SecureChat/internal/security"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/google/uuid"
)

var (
	PrivKey           string
	PubKey            string
	PeerPubKey        string
	IsAskingSessionId bool
	IsInChat          bool
)

func listenServerMessages(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		len, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading: " + err.Error())
			}
			break
		}

		message, err := dto.NewMessageFromString(string(buffer[:len]))
		if err != nil {
			log.Println("Error parsing message: ", err.Error())
			break
		}

		if message.Type == dto.ServerMessage {
			fmt.Println(message.Body)

		} else if message.Type == dto.ChatMessage {
			// TODO: Compute shared secret here
			fmt.Println("")

		} else if message.Type == dto.AskSessionIdMessage {
			fmt.Println(message.Body)
			IsAskingSessionId = true

		} else if message.Type == dto.AskKeyMessage {
			message := dto.NewMessage(dto.KeyExchangeMessage, PubKey+"\n")
			_, err := conn.Write(message.Bytes())
			if err != nil {
				log.Fatal(err)
			}

		} else if message.Type == dto.UserKeyExchangeMessage {
			PeerPubKey = message.Body
			log.Println("DEBUG: PeerPubKey => ", PeerPubKey)
			IsInChat = true

			ss, _ := security.ComputeSharedSecret(PrivKey, PeerPubKey)
			log.Println("DEBUG: Shared secret => ", ss)
		}
	}
}

func main() {
	IsAskingSessionId = false
	IsInChat = false

	var err error
	PrivKey, PubKey, err = security.GenerateKeyPair()
	if err != nil {
		panic(err)
	}
	log.Println("DEBUG: PubKey => ", PubKey)

	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		panic(err)
	}

	go listenServerMessages(conn)

	reader := bufio.NewReader(os.Stdin)

	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading line: ", err.Error())
			break
		}

		s = strings.ReplaceAll(s, "\r", "")

		var message *dto.Message
		if IsAskingSessionId {
			id, err := uuid.Parse(strings.ReplaceAll(s, "\n", ""))
			if err != nil {
				panic(err)
			}
			message = dto.NewJoinSessionMessage(id)
			IsAskingSessionId = false
		} else if IsInChat {
			// TODO: Send chat messages
		} else {
			message = dto.NewMessage(dto.ClientMessage, s)
		}

		_, err = conn.Write([]byte(message.String() + "\n"))
		if err != nil {
			log.Println("Error writing: ", err.Error())
			break
		}
	}
}
