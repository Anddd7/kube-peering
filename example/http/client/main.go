package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	concurrency = 10
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			defer wg.Done()
			post(index)
		}(i)
	}

	wg.Wait()
}

func post(index int) {
	url := fmt.Sprintf("http://localhost%s/ping", port())
	resp, err := http.Post(url, "text/plain", strings.NewReader(fmt.Sprintf("PING %d", index)))
	if err != nil {
		fmt.Printf("Error-%d: %v\n", index, err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error-%d: %v\n", index, err)
		return
	}

	fmt.Printf("Response-%d: %s\n", index, string(body))
}

func port() string {
	if len(os.Args) > 1 && os.Args[1] == "proxy" {
		return ":10021"
	}
	return ":8080"
}
