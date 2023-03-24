package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/kube-peering/internal/pkg/config"
)

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", config.DefautlFrontdoorPort))
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
