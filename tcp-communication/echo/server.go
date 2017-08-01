package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func handler(c net.Conn) {
	defer c.Close()
	s := bufio.NewScanner(c)
	for s.Scan() {
		fmt.Printf("Received: \"%s\" from %s\n", s.Text(), c.RemoteAddr().String())
	}
}

func main() {
	s, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer s.Close()

	for {
		c, err := s.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handler(c)
	}
}
