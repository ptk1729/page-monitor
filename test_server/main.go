package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"time"
)

var statusPage []byte

func init() {
	data, err := os.ReadFile("test_server/status.html")
	if err != nil {
		log.Fatalf("failed to load status.html: %v", err)
	}
	statusPage = data
}

func renderStatus(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if _, err := w.Write(statusPage); err != nil {
		log.Println("failed to write response:", err)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sec := time.Now().Second()
		randomNumber := rand.IntN(100) + 1
		if randomNumber < 1 {
			renderStatus(w, http.StatusInternalServerError)
			return
		}

		if sec == 10 && sec <= 20 {
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(200 * time.Millisecond)
		}

		renderStatus(w, http.StatusOK)
	})

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
