URL ?= http://localhost:8080
INTERVAL ?= 1s

build:
	go build -o bin/page-monitor .

run: 
	URL=$(URL) INTERVAL=$(INTERVAL) go run main.go

run-test-server:
	go run test_server/main.go