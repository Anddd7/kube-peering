package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", port())
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

func readServerMessage(conn io.Reader) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func port() string {
	if len(os.Args) > 1 {
		if os.Args[1] == "proxy" {
			return ":10021"
		}
		if os.Args[1] == "vpn" {
			return ":10022"
		}
	}
	return ":8080"
}
