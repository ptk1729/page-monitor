package probe

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		Help: "Availability of the webpage in percent",
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

func notifyOutage(url, reason string) {
	log.Printf("OUTAGE ALERT: %s - %s", url, reason)
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

func RunProbe(url string, interval time.Duration, failureThreshold int) {
	client := &http.Client{Timeout: 10 * time.Second}
	var total, success int

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		start := time.Now()
		total++
		totalChecks.Inc()

		resp, err := client.Get(url)
		duration := time.Since(start).Seconds()
		latencyHist.Observe(duration)

		data := ProbeData{Timestamp: time.Now(), DurationMS: int64(duration * 1000)}

		if err != nil {
			data.ErrorType = classifyError(err)
			consecutiveFails++
			log.Printf("Check failed for %s: %v", url, err)
		} else {
			data.StatusCode = resp.StatusCode
			statusCodeGauge.Set(float64(resp.StatusCode))
			resp.Body.Close()

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				success++
				successChecks.Inc()
				consecutiveFails = 0
			} else {
				consecutiveFails++
			}
		}

		avail := (float64(success) / float64(total)) * 100
		data.Availability = avail
		results = append(results, data)
		availabilityGauge.Set(avail)

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
