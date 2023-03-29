package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("got an err %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Recived connection from client: %s", conn.RemoteAddr().String())

	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("got an err %v", err)
			return
		}

		reply := fmt.Sprintf("PONG: %s", msg)

		conn.Write([]byte(reply))

		fmt.Printf("Recived message from client: %s, replied: %s", msg, reply)
	}
}
