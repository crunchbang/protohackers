package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
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

type Request struct {
	Method *string  `json:"method"`
	Number *float64 `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

const ERR_MSG = "invalid payload"

func handler(c net.Conn) {
	defer func() {
		log.Println("Closing connection")
		c.Close()
	}()

	buf := bufio.NewReader(c)
	for {
		payload, err := buf.ReadBytes('\n')
		if err != nil {
			log.Printf("could not read bytes from connection: %v", err)
			return
		}

		result, err := process(payload)
		if err != nil {
			log.Printf("unable to process payload: %v", err)
			c.Write([]byte(ERR_MSG))
			return
		}

		log.Println("#debug request ", string(payload), "response ", string(result))
		_, err = c.Write(result)
		if err != nil {
			log.Printf("could not write bytes to connection: %v", err)
			return
		}
		c.Write([]byte("\n"))
	}
}

func process(payload []byte) ([]byte, error) {
	req := Request{}
	err := json.Unmarshal(payload, &req)
	if err != nil {
		return nil, err
	}

	if req.Method == nil || req.Number == nil {
		return nil, fmt.Errorf("required field missing")
	}

	if *req.Method != "isPrime" {
		return nil, fmt.Errorf("method value mismatch")
	}

	resp := Response{
		Method: "isPrime",
		Prime:  isPrime(*req.Number),
	}

	return json.Marshal(resp)
}

func isPrime(num float64) bool {
	if math.Floor(num) != math.Ceil(num) || num <= 1 {
		return false
	}

	for i := float64(2); i <= math.Sqrt(num); i++ {
		if math.Mod(num, i) == 0 {
			return false
		}
	}
	return true
}
