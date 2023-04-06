package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	example "github.com/kube-peering/example"
)

func pingHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Recived connection from client: %s\n", r.RemoteAddr)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := string(body)

	if !strings.Contains(msg, "PING") {
		fmt.Printf("Recived message from client: %s\n", msg)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Recived message from client: %s\n", msg)
	fmt.Fprintf(w, "PONG: %s", msg)
}

func main() {
	server := http.Server{
		Addr:      example.AppHttpsAddr,
		TLSConfig: example.AppTlsConfig,
	}

	http.HandleFunc("/ping", pingHandlerFunc)

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}
