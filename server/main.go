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
)

var (
	config   *models.ConfigModel = &models.ConfigModel{}
	sessions *datastructures.LinkList[*models.SessionModel]
)

func main() {
	sessions = datastructures.NewLinkList[*models.SessionModel]()
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
	defer conn.Close()

	menuSequence(conn)
}

func menuSequence(conn net.Conn) error {
	defer conn.Close()

	isRunning := true
	var answer *dto.Message

	for isRunning {
		message := dto.NewMessage(dto.ServerMessage, "1. Open session\n2. Join session\n")

		_, err := conn.Write(message.Bytes())
		if err != nil {
			log.Println("Error, could not write to connection: ", err.Error())
			return err
		}

		answer, err = listenForMessage(conn)
		if err != nil {
			log.Println("Error, could not listen for message: ", err.Error())
			return err
		}

		isRunning = answer.Body != "1" && answer.Body != "2"
	}

	if answer.Body == "1" {
		// Create session
		session := models.NewSession()
		user := models.NewUser(conn, session.Id)
		session.Users.Add(user)
		m := dto.NewMessage(dto.ServerMessage, session.Id.String())
		_, err := conn.Write([]byte(m.String()))
		if err != nil {
			log.Println("Error, could not write to connection: ", err.Error())
			return err
		}

		sessions.Add(session)

		go listenForMessages(user)
	}

	return nil
}

func listenForMessage(conn net.Conn) (*dto.Message, error) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		// Disconnect
		return nil, err
	}

	log.Println("Received: ", message)

	return dto.NewMessageFromString(message)
}

func listenForMessages(u *models.User) {
	defer u.Conn.Close()

	for {
		message, err := listenForMessage(u.Conn)
		if err != nil {
			log.Println("Error, could not read connection: ", err.Error())
			break
		}

		_, err = u.Conn.Write(message.Bytes())
		if err != nil {
			log.Println("Error, could not write to connection: ", err.Error())
			break
		}
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
		session.Users.ForEach(func(u *models.User) {
			u.Conn.Close()
		})
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
