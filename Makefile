.PHONY: test run

run:
	cd explorer && go run main.go -bind :9090

test:
	go test -v ./core ./handlers

build:
	cd explorer && gulp build
	cd explorer && go build -o dist/explorer