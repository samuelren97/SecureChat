package main

import (
	"fmt"
	"log"
	"net"
)

const DEBUG bool = true

func main() {
	conn, err := net.Dial("tcp", "192.168.2.40:5000")
	if err != nil {
		panic(err)
	}

	for {
		var message string
		n, err := fmt.Scanln(&message)
		if err != nil {
			log.Println(err.Error())
			break
		}
		if DEBUG {
			log.Printf("DEBUG: Nb of bytes read => %d\n", n)
		}

		n, err = conn.Write([]byte(message))
		if err != nil {
			log.Println(err.Error())
			break
		}
		if DEBUG {
			log.Printf("DEBUG: Nb of bytes written => %d\n", n)
		}
	}
}
