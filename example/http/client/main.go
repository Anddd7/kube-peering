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
	concurrency = 11
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	client := &http.Client{}
	url := fmt.Sprintf("http://localhost%s/ping", port())
	for i := 0; i < concurrency; i++ {
		go func(index int) {
			defer wg.Done()

			resp, err := client.Post(url, "text/plain", strings.NewReader(fmt.Sprintf("PING %d", index)))
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
		}(i)
	}

	wg.Wait()
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
