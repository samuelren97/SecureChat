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
			ss, err := security.ComputeSharedSecret(PrivKey, PeerPubKey)
			if err != nil {
				log.Println("Error computing shared secret: ", err.Error())
				break
			}
			decryptedBody, err := security.DecryptAESGCM(message.Body, ss[:32])
			if err != nil {
				log.Println("Error decrypting message: ", err.Error())
				break
			}

			fmt.Println(decryptedBody)

		} else if message.Type == dto.AskSessionIdMessage {
			fmt.Println("Peer: " + message.Body)
			IsAskingSessionId = true

		} else if message.Type == dto.AskKeyMessage {
			message := dto.NewMessage(dto.KeyExchangeMessage, PubKey+"\n")
			_, err := conn.Write(message.Bytes())
			if err != nil {
				log.Fatal(err)
			}

		} else if message.Type == dto.UserKeyExchangeMessage {
			PeerPubKey = message.Body
			IsInChat = true
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
			ss, err := security.ComputeSharedSecret(PrivKey, PeerPubKey)
			if err != nil {
				log.Println("Error computing shared secret: ", err.Error())
				break
			}
			encryptedBody, err := security.EncryptAESGCM(s, ss[32:])
			if err != nil {
				log.Println("Error encrypting message: ", err.Error())
				break
			}

			message = dto.NewMessage(dto.ChatMessage, encryptedBody)
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
