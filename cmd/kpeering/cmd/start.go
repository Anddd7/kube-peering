package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/kube-peering/internal/logger"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitLogger()
		start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

type client chan<- string // 对于其他客户端的输出通道

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // 所有接收到的客户端消息
)

func broadcaster() {
	clients := make(map[client]bool) // 所有连接的客户端

	for {
		select {
		case msg := <-messages:
			// 将消息广播到所有客户端
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			// 新客户端连接
			clients[cli] = true
		case cli := <-leaving:
			// 客户端断开连接
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // 发送给客户端的消息通道
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // 发送给客户端
	}
}

func start() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to listen: %v\n", err)
		os.Exit(1)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to accept incoming connection: %v\n", err)
			continue
		}

		go handleConn(conn)
	}
}
