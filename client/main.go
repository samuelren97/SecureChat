package main

import (
	"SecureChat/internal/dto"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	privKey string
	pubKey  string
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
		}
	}
}

func main() {
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

		message := dto.NewMessage(dto.ClientMessage, s)

		_, err = conn.Write([]byte(message.String()))
		if err != nil {
			log.Println("Error writing: ", err.Error())
			break
		}
	}
}
