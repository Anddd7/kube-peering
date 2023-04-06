package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	example "github.com/kube-peering/example"
	"golang.org/x/net/http2"
)

const (
	concurrency = 10
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(concurrency)

	client := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig:    example.ClientTlsConfig,
			DisableCompression: true,
			AllowHTTP:          false,
		},
	}
	url := fmt.Sprintf("https://localhost:%d/ping", port())
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

func port() int {
	if len(os.Args) > 1 {
		if os.Args[1] == "proxy" {
			return example.ProxyPort
		}
		if os.Args[1] == "vpn" {
			return example.VPNPort
		}
	}
	return example.AppHttpsPort
}
