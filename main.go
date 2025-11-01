package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ptk1729/page-monitor/probe"
)

func main() {
	url := os.Getenv("URL")
	interval := os.Getenv("INTERVAL")
	failureThreshold := os.Getenv("FAILURE_THRESHOLD")

	if url == "" || interval == "" || failureThreshold == "" {
		log.Fatal("URL and INTERVAL environment variables are required")
	}

	failureThresholdInt, err := strconv.Atoi(failureThreshold)
	if err != nil {
		log.Fatalf("Invalid FAILURE_THRESHOLD: %v", err)
	}

	intervalDuration, err := time.ParseDuration(interval)
	if err != nil {
		log.Fatalf("Invalid INTERVAL: %v", err)
	}

	probe.RunProbe(url, intervalDuration, failureThresholdInt)
}
