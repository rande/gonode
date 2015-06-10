.PHONY: test run

default: clean test build

install:
	go list -f '{{range .Imports}}{{.}} {{end}}' ./... | xargs go get -v
	go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs go get -v
	go build -v ./...
	#cd explorer && npm install && npm install react-admin

update:
	go get -u all
	cd explorer && npm update && npm update react-admin

run:
	cd explorer && go run main.go -bind :9090

format:
	gofmt -l -w -s .

test:
	go test -v ./handlers ./test/api ./core
	#cd explorer && npm test

clean:
	rm -rf explorer/dist

build:
	cd explorer && gulp build
	cd explorer && go build -o dist/explorer

