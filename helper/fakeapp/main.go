package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/kube-peering/internal/pkg/logger"
)

func main() {
	logger.InitSimpleLogger()

	// 1. 监听8080端口
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Z.Error(err)
		return
	}
	defer ln.Close()

	// 2. 循环接受客户端的连接
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Z.Error(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// 在函数退出时关闭连接
	defer conn.Close()

	for {
		// 3. 读取客户端发来的消息
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			logger.Z.Error(err)
			return
		}

		// 4. 给消息加前缀
		reply := fmt.Sprintf("已读：%s", msg)

		// 5. 将加了前缀的消息回复给客户端
		conn.Write([]byte(reply))

		logger.Z.Infof("Recived message from client: %s, replied: %s", msg, reply)
	}
}
