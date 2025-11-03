# URL Monitor

A tool to monitor a given webpage

## Requirements

- Go 1.23 or later

## Usage

### Config

Set the following environment variables (See Makefile):
- `URL`: The URL to monitor
- `INTERVAL`: Check interval (e.g., `5s`, `1m`)
- `FAILURE_THRESHOLD`: Number of consecutive failures before outage is reported

### Running with Makefile

1. Run the mock server (skip in case you want to monitor another url)
```bash
make run-test-server
```

2. Run the monitor:
```bash
make run
```

### Running directly

```bash
URL=http://example.com INTERVAL=10s FAILURE_THRESHOLD=5 go run main.go
```

### Grafana Dashboard

1. Run Prometheus so it scrapes the probe metrics (default endpoint: `http://localhost:2112/metrics`).
2. Open Grafana -> Dashboards -> Import.
3. Upload `grafana/page-monitor-dashboard.json` and select your Prometheus datasource when prompted.
4. The dashboard shows availability, response times, status codes, and check volume in real time.

### Building

```bash
make build
```

Then run the binary:
```bash
URL=http://example.com INTERVAL=10s FAILURE_THRESHOLD=5 ./bin/page-monitor
```
