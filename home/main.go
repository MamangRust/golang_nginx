package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, HTTPS!"))
	})

	certPath := "/app/ssl/cert.pem"
	keyPath := "/app/ssl/key.pem"

	fmt.Printf("Server running on https://localhost\n")

	err := http.ListenAndServeTLS(":8082", certPath, keyPath, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
