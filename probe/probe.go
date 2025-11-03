package probe

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const TIMEOUT = 5 * time.Second

var (
	totalChecks = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "webpage_total_checks",
		Help: "Total number of checks performed",
	})
	successChecks = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "webpage_success_checks",
		Help: "Number of successful checks",
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
	prometheus.MustRegister(totalChecks, successChecks, latencyHist, statusCodeGauge)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics available at :2112/metrics")
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()
}

// isSuccessCode defines which HTTP codes count as succes
func isSuccessCode(code int) bool {
	if (code >= 200 && code <= 299) || code == 401 {
		return true
	}
	return false
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

		if err != nil {
			log.Printf("Check failed for %s: %v", url, err)
			continue
		}

		statusCodeGauge.Set(float64(resp.StatusCode))
		resp.Body.Close()

		if isSuccessCode(resp.StatusCode) {
			successChecks.Inc()
		}

		log.Printf("Checked %s -> %d, took %.2fs", url, resp.StatusCode, duration)
	}
}
