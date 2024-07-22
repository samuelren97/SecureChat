package servertools

import (
	"SecureChat/internal/datastructures"
	"SecureChat/internal/dto"
	"SecureChat/internal/models"
	"log"
	"net"

	"github.com/google/uuid"
)

func HandleOpenSession(conn net.Conn,
	user *models.User,
	sessions *datastructures.List[*models.SessionModel],
	askForPublicKey func(*models.User) error) error {
	// Create session
	var m *dto.Message
	session := models.NewSession()
	*user = *models.NewUser(conn, session.Id)

	err := session.AddUser(user)
	if err != nil {
		conn.Write(dto.NewMessage(dto.ServerMessage, err.Error()).Bytes())
		return err
	}
	m = dto.NewMessage(dto.ServerMessage, "Session id: "+session.Id.String())

	_, err = conn.Write([]byte(m.String()))
	if err != nil {
		log.Println("Error, could not write to connection: ", err.Error())
		return err
	}

	sessions.Add(session)
	askForPublicKey(user)
	return nil
}

func HandleJoinSession(conn net.Conn,
	user *models.User,
	sessions *datastructures.List[*models.SessionModel],
	messageListener func(net.Conn) (*dto.Message, error),
	askForPublicKey func(*models.User) error) error {
	message := dto.NewMessage(dto.AskSessionIdMessage, "Enter session Id:")
	_, err := conn.Write(message.Bytes())
	if err != nil {
		log.Println("Error, could not write to connection: ", err.Error())
		return err
	}

	responseMessage, err := messageListener(conn)
	if err != nil {
		log.Println("Error, could not read message: ", err.Error())
		return err
	}

	isFound := false
	for i := 0; !isFound && i < sessions.Count; i++ {
		session := sessions.Get(i)
		if isFound = session.Id == responseMessage.SessionId; isFound {
			id, err := uuid.NewUUID()
			if err != nil {
				log.Println("Error, could not create UUID: ", err.Error())
				return err
			}
			*user = *models.NewUser(conn, id)
			if err = session.AddUser(user); err != nil {
				conn.Write([]byte("Can't join the session: " + err.Error()))
				return err
			}

			askForPublicKey(user)
			session.SendKeys()
			break
		}
	}

	if !isFound {
		conn.Write([]byte("session id was not found"))
		// TODO:
	}
	return nil
}
