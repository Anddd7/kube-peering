package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/kube-peering/internal/pkg/io"
	"github.com/kube-peering/internal/pkg/logger"
)

func main() {
	io.StartTCPServer(":8080", func(s string) {}, handleConnection)
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
