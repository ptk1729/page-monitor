URL ?= https://google.com
INTERVAL ?= 1s
FAILURE_THRESHOLD ?= 3

build:
	go build -o bin/page-monitor .

run: 
	URL=$(URL) INTERVAL=$(INTERVAL) FAILURE_THRESHOLD=$(FAILURE_THRESHOLD) go run main.go

run-test-server:
	go run test_server/main.go