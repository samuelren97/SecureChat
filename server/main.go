package main

import (
	"SecureChat/internal/datastructures"
	"SecureChat/internal/dto"
	"SecureChat/internal/models"
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"

	"github.com/google/uuid"
)

var (
	config   *models.ConfigModel = &models.ConfigModel{}
	sessions *datastructures.List[*models.SessionModel]
)

func main() {
	sessions = datastructures.NewList[*models.SessionModel]()
	loadConfig()

	listener, err := net.Listen("tcp", config.Server.Port)
	if err != nil {
		log.Fatal("Error, could not start TCP server: ", err.Error())
	}

	log.Println("Server is listening on port: ", config.Server.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error, could not accept incoming connection: ", err.Error())
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	user, err := menuSequence(conn)
	if err != nil {
		return
	}
	go listenForMessages(user)
}

func menuSequence(conn net.Conn) (*models.User, error) {
	isRunning := true
	var answer *dto.Message
	var user *models.User

	for isRunning {
		message := dto.NewMessage(dto.ServerMessage, "1. Open session\n2. Join session\n")

		_, err := conn.Write(message.Bytes())
		if err != nil {
			log.Println("Error, could not write to connection: ", err.Error())
			return nil, err
		}

		answer, err = listenForMessage(conn)
		if err != nil {
			log.Println("Error, could not listen for message: ", err.Error())
			return nil, err
		}

		isRunning = answer.Body != "1" && answer.Body != "2"
		print('e')
	}

	if answer.Body == "1" {
		// Create session
		var m *dto.Message
		session := models.NewSession()
		user = models.NewUser(conn, session.Id)

		err := session.AddUser(user)
		if err != nil {
			conn.Write(dto.NewMessage(dto.ServerMessage, err.Error()).Bytes())
			return nil, err
		}
		m = dto.NewMessage(dto.ServerMessage, "Session id: "+session.Id.String())

		_, err = conn.Write([]byte(m.String()))
		if err != nil {
			log.Println("Error, could not write to connection: ", err.Error())
			return nil, err
		}

		sessions.Add(session)

		askKeyMessage := dto.NewMessage(dto.AskKeyMessage, "")
		_, err = conn.Write(askKeyMessage.Bytes())
		if err != nil {
			log.Println("Error, could not write to connection: ", err.Error())
			return nil, err
		}

		responseMessage, err := listenForMessage(conn)
		if err != nil {
			log.Println("Error, could not read message: ", err.Error())
			return nil, err
		}

		if responseMessage.Type != dto.KeyExchangeMessage {
			err = dto.ErrExpectedKeyExhange
			log.Println("Error, could not process message: ", err.Error())
			return nil, err
		}

		user.PubKey = responseMessage.Body

	} else if answer.Body == "2" {
		message := dto.NewMessage(dto.AskSessionIdMessage, "Enter session Id:")
		_, err := conn.Write(message.Bytes())
		if err != nil {
			log.Println("Error, could not write to connection: ", err.Error())
			return nil, err
		}

		responseMessage, err := listenForMessage(conn)
		if err != nil {
			log.Println("Error, could not read message: ", err.Error())
			return nil, err
		}

		for i := 0; i < sessions.Count; i++ {
			session := sessions.Get(i)
			if session.Id == responseMessage.SessionId {
				id, err := uuid.NewUUID()
				if err != nil {
					log.Println("Error, could not create UUID: ", err.Error())
					return nil, err
				}
				user = models.NewUser(conn, id)
				if err = session.AddUser(user); err != nil {
					conn.Write([]byte("Can't join the session: " + err.Error()))
					return nil, err
				}

				askKeyMessage := dto.NewMessage(dto.AskKeyMessage, "")
				_, err = conn.Write(askKeyMessage.Bytes())
				if err != nil {
					log.Println("Error, could not write to connection: ", err.Error())
					return nil, err
				}

				responseMessage, err := listenForMessage(conn)
				if err != nil {
					log.Println("Error, could not read message: ", err.Error())
					return nil, err
				}

				if responseMessage.Type != dto.KeyExchangeMessage {
					err = dto.ErrExpectedKeyExhange
					log.Println("Error, could not process message: ", err.Error())
					return nil, err
				}

				user.PubKey = responseMessage.Body

				session.SendKeys()
				break
			}
		}
	}

	return user, nil
}

func listenForMessage(conn net.Conn) (*dto.Message, error) {
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		// Disconnect
		return nil, err
	}

	message = strings.ReplaceAll(message, "\n", "")
	log.Println("Received: ", message)

	return dto.NewMessageFromString(message)
}

func listenForMessages(u *models.User) {
	for {
		message, err := listenForMessage(u.Conn)
		if err != nil {
			log.Println("Error, could not read connection: ", err.Error())
			break
		}

		log.Println("Received: ", message.Body)
	}

	isFound := false
	var sessionIndex int = 0
	sessions.ForEach(func(session *models.SessionModel) {
		isFound = session.Id.String() == u.SessionId.String()

		if !isFound {
			sessionIndex++
		}
	})

	if isFound {
		session := sessions.Get(sessionIndex)
		session.Close()
		sessions.Remove(sessionIndex)
	}
}

func loadConfig() {
	file, err := os.OpenFile("appsettings.json", os.O_RDONLY, 0444)
	if err != nil {
		log.Fatal("Error, could not open appsettings file: ", err.Error())
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		log.Fatal("Error, could not decode appsettings to config obj: ", err.Error())
	}
}
