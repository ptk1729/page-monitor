# URL Monitor

A tool to monitor a given webpage

## Requirements

* Go 1.23+
* Prometheus
* Grafana

## Setup

### 1. Configure

Set the following environment variables (See Makefile):
- `URL`: The URL to monitor
- `INTERVAL`: Check interval (e.g., `5s`, `1m`)

Use the Makefile:

```bash
make run
```

Or run directly:

```bash
URL=http://example.com INTERVAL=10s go run main.go
```

For testing, the small server can be used

Run (using make file)
```bash
make run-test-server
```

or
```bash
go run test_server/main.go
```
### 3. Prometheus

Start Prometheus with:

```bash
--config.file=prometheus.yml --web.enable-lifecycle
```
example:
```bash
docker run -d \
  --name prometheus \
  -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  -v $(pwd)/alert_rules.yml:/etc/prometheus/alert_rules.yml \
  prom/prometheus \
  --config.file=/etc/prometheus/prometheus.yml \
  --web.enable-lifecycle
```


Metrics are at: [http://localhost:2112/metrics](http://localhost:2112/metrics)

### 4. Grafana

Run locally
```bash
grafana server path/to/conf
```

or with docker
```bash
docker run -d -p 3000:3000 grafana/grafana 
```

1. Open Grafana, import dashboard json
2. Upload `grafana/page-monitor-dashboard.json`
3. Select your Prometheus datasource

You’ll see:

* Availability (1m, 5m, 1h)
* Response time (avg & p95)
* HTTP status
* Page loads over time

### 5. Build

```bash
make build
./bin/page-monitor
```

## Design

Given a monitoring target, a webpage, this tool probes it and pushes metrics, then scraped by prometheus.

A status code <500 is considered success (e.g., 401, 404 are treated as available).

Requests timing out are considered to be an error.

Availability is computed from success vs total checks within a time window (e.g., 2 m, 5 m).

Alerts fire when availability < 95 % over the last 2 minutes.

Grafana connects to the same Prometheus datasource for visualization.