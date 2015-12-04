.PHONY: test run explorer build

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
	cd commands && go run main.go server -config=../server.toml.dist

build:
	rm -rf dist && mkdir dist
	#cd explorer && webpack --progress --color
	#cd commands && go build -o dist/gonode
	cd commands && go build -o ../dist/gonode

format:
	gofmt -l -w -s .
	go fix ./...

test:
	go test -v ./handlers ./test/api ./core ./vault ./commands/server
	go vet ./...
	#cd explorer && npm test

kill:
	kill `cat $(PID)` || true

serve: clean
	make restart
	#cd explorer && node webpack.config.js $$! > $(PID)_wds &
	fswatch $(GO_FILES) | xargs -n1 -I{} make restart || make kill
	kill `cat $(PID)_wds` || true

restart:
	make kill
	rm -rf dist/gonode
	cd commands && go build -o ../dist/gonode
	cp server.toml.dist dist/config.toml
	cd dist && (./gonode server & echo $$! > ../../$(PID))
