.PHONY: test run

run:
	go run gnode/main.go

test:
	go test -v ./core ./handlers

