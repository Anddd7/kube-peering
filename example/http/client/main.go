package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	resp, err := http.Post("http://localhost:10021/ping", "text/plain", strings.NewReader("PING"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(body))
}
