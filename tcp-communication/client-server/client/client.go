package main

import (
	"fmt"
	"io/ioutil"
	"net"
)

func main() {
	connection, err := net.Dial("tcp", "localhost:50000")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	msg, _ := ioutil.ReadAll(connection)
	fmt.Println(string(msg))
}
