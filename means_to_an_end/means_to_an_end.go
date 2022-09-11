package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
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

		session := Session{CmdChannel: conn}
		go session.process()
	}
}

const CMD_LEN = 9

type Session struct {
	CmdChannel io.ReadWriteCloser
	store      map[int32]int32
}

func (s *Session) process() {
	log.Println("Session initalized")
	s.store = map[int32]int32{}

	defer func() {
		s.CmdChannel.Close()
		log.Println("closing connection")
	}()
	for {
		buf := make([]byte, 9)
		n, err := io.ReadAtLeast(s.CmdChannel, buf, CMD_LEN)
		if err != nil {
			if err != io.EOF {
				log.Println("error while reading from cmd channel: ", err, buf, n)
			}
			return
		}

		log.Println("#debug read", n, " bytes ", buf)
		err = s.processCmd(buf)
		if err != nil {
			log.Println("unable to process command: ", err)
			return
		}
	}
}

func (s *Session) processCmd(buf []byte) error {
	switch buf[0] {
	case []byte("I")[0]:
		timestamp := int32(binary.BigEndian.Uint32(buf[1:5]))
		price := int32(binary.BigEndian.Uint32(buf[5:9]))
		return s.processInsert(timestamp, price)
	case []byte("Q")[0]:
		minTime := int32(binary.BigEndian.Uint32(buf[1:5]))
		maxTime := int32(binary.BigEndian.Uint32(buf[5:9]))
		return s.processQuery(minTime, maxTime)
	}

	return fmt.Errorf("unknown command")
}

func (s *Session) processInsert(timestamp, price int32) error {
	log.Printf("#debug got insert %d %d\n", timestamp, price)
	s.store[timestamp] = price
	return nil
}

func (s *Session) processQuery(minTime, maxTime int32) error {
	count := int64(0)
	sum := int64(0)
	for time, price := range s.store {
		if minTime <= time && time <= maxTime {
			count += 1
			sum += int64(price)
		}
	}

	mean := int32(0)
	if count > 0 {
		mean = int32(sum / count)
	}

	log.Printf("#debug got query %d %d sum %d count %d mean %d\n", minTime, maxTime, sum, count, mean)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(mean))
	_, err := s.CmdChannel.Write(buf)
	return err
}
