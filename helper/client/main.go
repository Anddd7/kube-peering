package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to server: %v\n", err)
		os.Exit(1)
	}

	go readServerMessage(conn) // 启动协程读取服务器消息

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fmt.Fprintf(conn, "%s\n", input.Text()) // 将输入发送给服务器
	}
}

func readServerMessage(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() { // 循环读取来自服务器的消息
		fmt.Println(scanner.Text()) // 打印消息
	}
}
