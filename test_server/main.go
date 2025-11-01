package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sec := time.Now().Second()

		if sec >= 50 && sec <= 59 {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "404 - Not Found (temporary)")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "200 - OK")
	})

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
