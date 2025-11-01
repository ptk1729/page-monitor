package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	startTime := time.Now()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		elapsed := time.Since(startTime).Seconds()
		cycle := int(elapsed) % 60

		if cycle >= 50 {
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
