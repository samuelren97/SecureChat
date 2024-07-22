package main

import (
	"SecureChat/internal/datastructures"
	"SecureChat/internal/dto"
	"SecureChat/internal/models"
	"SecureChat/server/servertools"
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
)

var (
	config   *models.ConfigModel = &models.ConfigModel{}
	sessions *datastructures.List[*models.SessionModel]
)

func main() {
	sessions = datastructures.NewList[*models.SessionModel]()
	loadConfig()

	listener, err := net.Listen("tcp", config.Server.Address)
	if err != nil {
		log.Fatal("Error, could not start TCP server: ", err.Error())
	}

	log.Println("Server is listening on address: ", config.Server.Address)

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
	var user *models.User = &models.User{}

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
		servertools.HandleOpenSession(
			conn,
			user,
			sessions,
			askForPublicKey,
		)
	} else if answer.Body == "2" {
		servertools.HandleJoinSession(
			conn,
			user,
			sessions,
			listenForMessage,
			askForPublicKey,
		)
	}

	return user, nil
}

func askForPublicKey(user *models.User) error {
	askKeyMessage := dto.NewMessage(dto.AskKeyMessage, "")
	_, err := user.Conn.Write(askKeyMessage.Bytes())
	if err != nil {
		log.Println("Error, could not write to connection: ", err.Error())
		return err
	}

	responseMessage, err := listenForMessage(user.Conn)
	if err != nil {
		log.Println("Error, could not read message: ", err.Error())
		return err
	}

	if responseMessage.Type != dto.KeyExchangeMessage {
		err = dto.ErrExpectedKeyExhange
		log.Println("Error, could not process message: ", err.Error())
		return err
	}

	user.PubKey = responseMessage.Body
	return nil
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
