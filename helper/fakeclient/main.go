package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/model"
)

func main() {
	logger.InitLogger()

	conn, err := net.Dial("tcp", model.DefaultFrontdoor.Address())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to server: %v\n", err)
		os.Exit(1)
	}

	go readServerMessage(conn)

	input := bufio.NewScanner(os.Stdin)
	fmt.Println("Connected, please input your message:")
	for input.Scan() {
		fmt.Fprintf(conn, "%s\n", input.Text())
	}
}

func readServerMessage(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
