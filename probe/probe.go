package probe

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ProbeData struct {
	Timestamp    time.Time
	StatusCode   int
	DurationMS   int64
	ErrorType    string
	Availability float64
}

var (
	results          []ProbeData
	consecutiveFails int
	outageActive     bool
)

func notifyOutage(url, reason string) {
	log.Printf("OUTAGE: %s - %s", url, reason)
}

func classifyError(err error) string {
	switch {
	case err == nil:
		return ""
	case os.IsTimeout(err):
		return "timeout"
	default:
		return "network"
	}
}

// RunProbe runs periodic checks against a URL at a given interval.
func RunProbe(url string, interval time.Duration, failureThreshold int) {
	client := &http.Client{Timeout: 10 * time.Second}
	var total, success int

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		start := time.Now()
		total++

		resp, err := client.Get(url)
		duration := time.Since(start).Seconds()
		data := ProbeData{Timestamp: time.Now(), DurationMS: int64(duration * 1000)}

		if err != nil {
			data.ErrorType = classifyError(err)
			consecutiveFails++
			log.Printf("Check failed for %s: %v", url, err)
		} else {
			data.StatusCode = resp.StatusCode
			resp.Body.Close()

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				success++
				consecutiveFails = 0
			} else {
				consecutiveFails++
			}
		}

		avail := (float64(success) / float64(total)) * 100
		data.Availability = avail
		results = append(results, data)
		log.Printf("Availability: %.2f%% (failures: %d)", avail, consecutiveFails)

		if consecutiveFails >= failureThreshold && !outageActive {
			outageActive = true
			notifyOutage(url, fmt.Sprintf("%d consecutive failures", consecutiveFails))
		} else if consecutiveFails == 0 && outageActive {
			outageActive = false
			log.Printf("Service is back online for %s", url)
		}
	}
}
