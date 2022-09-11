package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
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
	Method *string         `json:"method"`
	Number json.RawMessage `json:"number"`
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

	if req.Method == nil || len(req.Number) == 0 {
		return nil, fmt.Errorf("required field missing")
	}

	if *req.Method != "isPrime" {
		return nil, fmt.Errorf("method value mismatch")
	}

	if req.Number[0] == '"' {
		return nil, fmt.Errorf("incorrect number value")

	}

	num, _, err := big.ParseFloat(string(req.Number), 10, 127, big.ToNearestAway)
	if err != nil {
		return nil, fmt.Errorf("incorrect number value")
	}

	resp := Response{
		Method: "isPrime",
		Prime:  isPrimeBigInt(num),
	}

	return json.Marshal(resp)
}

func isPrimeBigInt(num *big.Float) bool {
	if !num.IsInt() {
		return false
	}

	n := new(big.Int)
	num.Int(n)
	if n.Cmp(big.NewInt(1)) < 1 {
		return false
	}

	sqrtN := new(big.Int)
	sqrtN.Sqrt(n)

	for i := big.NewInt(2); i.Cmp(sqrtN) < 1; i.Add(i, big.NewInt(1)) {
		mod := new(big.Int)
		if mod.Mod(n, i).Cmp(big.NewInt(0)) == 0 {
			return false
		}
	}
	return true
}
