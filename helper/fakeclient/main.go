package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/kube-peering/internal/config"
	"github.com/kube-peering/internal/logger"
)

/*
simulate the requester in cluster, send the request to pod which is hosting the kpeering
*/
func main() {
	logger.InitLogger()

	conn, err := net.Dial("tcp", config.KpeeringForwardAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to server: %v\n", err)
		os.Exit(1)
	}

	go readServerMessage(conn) // 启动协程读取服务器消息

	input := bufio.NewScanner(os.Stdin)
	fmt.Println("connected, please input your message:")
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
