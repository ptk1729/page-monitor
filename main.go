package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	url := os.Getenv("URL")
	if url == "" {
		fmt.Println("URL environment variable is required")
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s is DOWN: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("%s is UP\n", url)
	} else {
		fmt.Printf("%s returned status %d\n", url, resp.StatusCode)
	}
}
