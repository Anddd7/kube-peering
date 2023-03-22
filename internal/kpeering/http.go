package kpeering

import (
	"bufio"
	"log"
	"net"
	"net/http"
)

func main() {
	// 监听代理服务器的端口
	proxy, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	// 接受连接请求并设置代理目标
	for {
		conn, err := proxy.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}

		// 解析 HTTP 请求
		_, port, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			log.Println("Error parsing address: ", err)
			continue
		}

		req, err := http.ReadRequest(bufio.NewReader(conn))
		if err != nil {
			log.Println("Error reading request: ", err)
			continue
		}

		// 连接到目标服务器并转发请求
		client, err := net.Dial("tcp", req.Host+":"+port)
		if err != nil {
			log.Println("Error connecting to target: ", err)
			continue
		}

		req.RequestURI = ""
		err = req.Write(client)
		if err != nil {
			log.Println("Error writing request: ", err)
			continue
		}

		// 从目标服务器读取响应并转发到客户端
		resp, err := http.ReadResponse(bufio.NewReader(client), req)
		if err != nil {
			log.Println("Error reading response: ", err)
			continue
		}

		resp.Request = nil
		err = resp.Write(conn)
		if err != nil {
			log.Println("Error writing response: ", err)
			continue
		}
	}
}
