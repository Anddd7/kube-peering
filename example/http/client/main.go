package main

import (
	"fmt"
	"io"
	"net/http"
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
	resp, err := http.Post("http://localhost:10021/ping", "text/plain", strings.NewReader(fmt.Sprintf("PING %d", index)))
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
