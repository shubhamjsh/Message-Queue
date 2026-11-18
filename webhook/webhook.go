package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "could not read request", http.StatusBadRequest)
			return
		}
		log.Printf("Webhook received: %s", body)
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Webhook listening at http://localhost:4001/webhook")
	log.Fatal(http.ListenAndServe(":4001", nil))
}
