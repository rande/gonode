.PHONY: test run explorer build

PID = .pid
GO_FILES = $(shell find . -type f -name "*.go")
GONODE_PLUGINS = $(shell ls -d ./plugins/* | grep -v go)

GO_PATH = $(shell go env GOPATH)
GO_BINDATA_PATHS = $(GO_PATH)/src/github.com/rande/gonode/plugins/... $(GO_PATH)/src/github.com/rande/gonode/explorer/dist/...
GO_BINDATA_IGNORE = "(.*)\.(go|DS_Store)"
GO_BINDATA_OUTPUT = $(GO_PATH)/src/github.com/rande/gonode/assets/bindata.go
GO_BINDATA_PACKAGE = assets

default: clean test build

clean:
	rm -rf explorer/dist/*

install-backend:
	go list -f '{{range .Imports}}{{.}} {{end}}' ./... | xargs go get -v
	go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs go get -v
	go build -v ./...

install-frontend:
	cd explorer && npm install

install: install-backend install-frontend

update:
	go get -u all
	cd explorer && npm update

run: bin
	cd commands && go run main.go server -config=../server.toml.dist

bin:
	cd $(GO_PATH)/src && go-bindata -dev -prefix $(GO_PATH)/src -o $(GO_BINDATA_OUTPUT) -pkg $(GO_BINDATA_PACKAGE) -ignore $(GO_BINDATA_IGNORE) $(GO_BINDATA_PATHS)

build:
	rm -rf dist && mkdir dist
	cd explorer && npm run-script build
	cd $(GO_PATH)/src && go-bindata -prefix $(GO_PATH)/src -o $(GO_BINDATA_OUTPUT) -pkg $(GO_BINDATA_PACKAGE) -ignore $(GO_BINDATA_IGNORE)  $(GO_BINDATA_PATHS)
	cd commands && go build -a -o ../dist/gonode

format:
	gofmt -l -w -s .
	go fix ./...

test-backend:
	go test $(GONODE_PLUGINS) ./test/api ./core ./core/config ./commands/server
	go vet ./...

test-frontend:
	cd explorer && npm test

test: test-backend test-frontend

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
