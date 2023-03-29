package main

import (
	"fmt"
	"io"
	"net/http"
)

func pingHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Recived connection from client: %s", r.RemoteAddr)

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

	if msg != "PING" {
		fmt.Printf("Recived message from client: %s", msg)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Recived message from client: %s", msg)
	fmt.Fprintf(w, "PONG: %s", msg)
}

func main() {
	http.HandleFunc("/ping", pingHandlerFunc)
	fmt.Println("Server listening at :8080 ...")
	http.ListenAndServe(":8080", nil)
}