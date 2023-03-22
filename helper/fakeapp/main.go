package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/kube-peering/internal/pkg/io"
	"github.com/kube-peering/internal/pkg/logger"
)

func main() {
	logger.InitSimpleLogger()

	connChan := make(chan net.Conn)

	go io.AcceptConnections("test_app", "tcp", ":8080", connChan)

	for conn := range connChan {
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			logger.Z.Error(err)
			return
		}

		reply := fmt.Sprintf("replay: %s", msg)

		conn.Write([]byte(reply))

		logger.Z.Infof("Recived message from client: %s, replied: %s", msg, reply)
	}
}
