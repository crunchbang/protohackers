package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

const (
	DEFAULT_PORT = 8899
)

func main() {
	port := flag.Int("p", DEFAULT_PORT, "port to listen for tcp connections")
	flag.Parse()

	log.Printf("starting tcp server on %d\n", *port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("unable to listen on %d: %v", *port, err)
	}
	defer listener.Close()

	log.Println("waiting for connections...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("unable to accept connection: %v", err)
		}

		log.Println("accepted a connection")
		go handler(conn)
	}
}

func handler(c net.Conn) {
}
