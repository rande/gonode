.PHONY: test run

default: clean test build

install:
	cd explorer && npm install && npm install react-admin

run:
	cd explorer && go run main.go -bind :9090

test:
	go test -v ./core ./handlers
	cd explorer && npm test

clean:
	rm -rf explorer/dist

build:
	cd explorer && gulp build
	cd explorer && go build -o dist/explorer

