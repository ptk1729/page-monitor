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

Metrics are at: [http://localhost:2112/metrics](http://localhost:2112/metrics)

### 4. Grafana

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