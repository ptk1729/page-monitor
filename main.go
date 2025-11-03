package main

import (
	"log"
	"os"
	"time"

	"github.com/ptk1729/page-monitor/probe"
)

func main() {
	url := os.Getenv("URL")
	interval := os.Getenv("INTERVAL")

	if url == "" || interval == "" {
		log.Fatal("URL and INTERVAL environment variables are required")
	}

	intervalDuration, err := time.ParseDuration(interval)
	if err != nil {
		log.Fatalf("Invalid INTERVAL: %v", err)
	}

	probe.RunProbe(url, intervalDuration)
}
