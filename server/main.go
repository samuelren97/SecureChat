package main

import (
	"SecureChat/internal/models"
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
)

var config *models.ConfigModel = &models.ConfigModel{}

func main() {
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

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Warning, could not read message: ", err.Error(), "\n this can be caused by client connection lost")
			break
		}

		log.Println("Received: ", message)

		_, err = conn.Write([]byte("Echo: " + message))
		if err != nil {
			log.Println("Error, could not write to connection: ", err.Error())
			break
		}
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
