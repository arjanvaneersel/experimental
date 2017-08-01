package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func handler(c net.Conn) {
	defer c.Close()
	s := bufio.NewScanner(c)
	i := 0
	var method string
	var route string

	for s.Scan() {
		item := s.Text()
		fmt.Println(item)

		if i == 0 {
			method = strings.Fields(item)[0]
			route = strings.Fields(item)[1]
			fmt.Printf("METHOD: %s\n", method)
			fmt.Printf("ROUTE: %s\n", route)
		} else {
			if item == "" {
				break
			}
		}
		i++
	}

	response := `
    <!DOCTYPE html>
    <html lang="en"
      <head>
        <meta charset="UTF-8">
        <title>Bare HTTP server version 0.1</title>
      </head>
      <body>
        <form method="POST">
          <input type="text" name="key" value="">
          <input type="submit">
        </form>
      </body>
    </html>
  `
	io.WriteString(c, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(c, "Content-Length: %d\r\n", len(response))
	io.WriteString(c, "\r\n")
	io.WriteString(c, response)
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
