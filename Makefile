.PHONY: test run explorer

PID = .pid
GO_FILES = $(shell find . -type f -name "*.go")

default: clean test build

clean:
	rm -rf explorer/dist/*

install:
	go list -f '{{range .Imports}}{{.}} {{end}}' ./... | xargs go get -v
	go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs go get -v
	go build -v ./...
	#cd explorer && npm install && npm install react-admin

update:
	go get -u all
	cd explorer && npm update && npm update react-admin

run:
	cd cli && go run main.go server -config=../server.toml.dist

format:
	gofmt -l -w -s .
	go fix ./...

test:
	go test ./handlers ./test/api ./core ./vault
	go vet ./...
	#cd explorer && npm test

kill:
	kill `cat $(PID)` || true

build:
	cd explorer && webpack --progress --color
	cd explorer && go build -o dist/explorer

serve: clean
	make restart
	cd explorer && node webpack.config.js $$! > $(PID)_wds &
	fswatch $(GO_FILES) | xargs -n1 -I{} make restart || make kill
	kill `cat $(PID)_wds` || true

restart:
	make kill
	cd explorer && rm -rf dist/explorer
	cd explorer && go build -o dist/explorer
	cd explorer && cp config.toml dist/config.toml
	cd explorer/dist && (./explorer -bind :9090 & echo $$! > ../../$(PID))
