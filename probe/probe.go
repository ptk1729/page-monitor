package probe

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ptk1729/page-monitor/notify"
)

type ProbeData struct {
	Timestamp  time.Time
	StatusCode int
	DurationMS int64
	ErrorType  string
	Success    bool
}

const TIMEOUT = 5 * time.Second
const WindowDuration = 2 * time.Minute
const AvailabilityThreshold = 95.0

var (
	results      []ProbeData
	outageActive bool

	totalChecks = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "webpage_total_checks",
		Help: "Total number of checks performed",
	})
	successChecks = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "webpage_success_checks",
		Help: "Number of successful checks",
	})
	availabilityGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "webpage_availability_percent",
		Help: "Availability of the webpage in percent (last 2 minutes)",
	})
	latencyHist = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "webpage_response_time_seconds",
		Help:    "Response time of the webpage in seconds",
		Buckets: prometheus.DefBuckets,
	})
	statusCodeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "webpage_last_status_code",
		Help: "Last HTTP status code returned by the webpage",
	})
)

func init() {
	prometheus.MustRegister(totalChecks, successChecks, availabilityGauge, latencyHist, statusCodeGauge)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics available at :2112/metrics")
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()
}

// notifyOutage just logs. Replace with a call to another module later.
func notifyOutage(url, reason string) {
	log.Printf("OUTAGE ALERT: %s - %s", url, reason)
	notify.Send("slack", fmt.Sprintf("%s is down: %s", url, reason))
}

func classifyError(err error) string {
	if err == nil {
		return ""
	}
	if os.IsTimeout(err) {
		return "timeout"
	}
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return "timeout"
	}
	return "network"
}

// HTTP 5xx are treated as failure.
func isSuccessCode(code int) bool {
	if code >= 500 {
		return false
	}
	return true
}

// computeAvailability calculates success percentage within the last 2 minutes.
func computeAvailability() float64 {
	cutoff := time.Now().Add(-WindowDuration)
	var total, success int
	for _, r := range results {
		if r.Timestamp.After(cutoff) {
			total++
			if r.Success {
				success++
			}
		}
	}
	if total == 0 {
		return 100.0
	}
	return (float64(success) / float64(total)) * 100
}

func RunProbe(url string, interval time.Duration) {
	client := &http.Client{Timeout: TIMEOUT}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		start := time.Now()
		totalChecks.Inc()

		resp, err := client.Get(url)
		duration := time.Since(start).Seconds()
		latencyHist.Observe(duration)

		data := ProbeData{Timestamp: time.Now(), DurationMS: int64(duration * 1000)}

		if err != nil {
			data.ErrorType = classifyError(err)
			data.Success = false
			log.Printf("Check failed for %s: %v", url, err)
		} else {
			data.StatusCode = resp.StatusCode
			statusCodeGauge.Set(float64(resp.StatusCode))
			resp.Body.Close()
			data.Success = isSuccessCode(resp.StatusCode)
			if data.Success {
				successChecks.Inc()
			}
		}

		results = append(results, data)
		avail := computeAvailability()
		availabilityGauge.Set(avail)

		log.Printf("2-min Availability: %.2f%%, took %.2f seconds", avail, duration)

		if avail < AvailabilityThreshold && !outageActive {
			outageActive = true
			notifyOutage(url, fmt.Sprintf("availability dropped to %.2f%%", avail))
		} else if avail >= AvailabilityThreshold && outageActive {
			outageActive = false
			log.Printf("Service recovered for %s (availability %.2f%%)", url, avail)
		}
	}
}
