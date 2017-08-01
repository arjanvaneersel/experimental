package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":50000")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		connection, err := l.Accept()
		if err != nil {
			panic(err)
		}

		io.WriteString(connection, fmt.Sprint("Hi there!\n", time.Now(), "\n"))

		connection.Close()
	}
}
