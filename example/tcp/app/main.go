// nolint:G102
package main

import (
	"bufio"
	"fmt"
	"net"

	example "github.com/kube-peering/example"
)

func main() {
	listener, err := net.Listen("tcp", example.AppAddr)
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

	fmt.Printf("Recived connection from client: %s\n", conn.RemoteAddr().String())

	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("got an err %v\n", err)
			return
		}

		reply := fmt.Sprintf("PONG: %s", msg)

		_, err = conn.Write([]byte(reply))
		if err != nil {
			fmt.Printf("got an err %v", err)
			return
		}
		fmt.Printf("Recived message from client, replied: %s", reply)
	}
}
