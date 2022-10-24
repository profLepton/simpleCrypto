package main

import (
	"log"
	"net"
)

var (
	ADDRESS = ":8888"
)

func main() {
	log.Print("Bob> Starting Bob.")
	log.Print("Listening for connections.")

	serverListener, err := net.Listen("tcp", ADDRESS)
	if err != nil {
		panic(err)
	}

	for {
		serverConnection, err := serverListener.Accept()
		if err != nil {
			panic(err)
		}

		go HandleConnections(serverConnection)

	}
}

func HandleConnections(conn net.Conn) {
	//pass
}
