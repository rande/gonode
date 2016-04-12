.PHONY: test run explorer build

PID = .pid
GO_FILES = $(shell find . -type f -name "*.go")
GONODE_MODULES = $(shell ls -d ./modules/* | grep -v go)
GONODE_CORE = $(shell ls -d ./core/* | grep -v go)
GO_PATH = $(shell go env GOPATH)
GO_BINDATA_PATHS = $(GO_PATH)/src/github.com/rande/gonode/modules/... $(GO_PATH)/src/github.com/rande/gonode/explorer/dist/...
GO_BINDATA_IGNORE = "(.*)\.(go|DS_Store)"
GO_BINDATA_OUTPUT = $(shell pwd)/assets/bindata.go
GO_BINDATA_PACKAGE = assets

default: test

format:
	goimports -w $(GO_FILES)
	gofmt -l -w -s .
	go fix ./...
	go vet ./...

bin:
	cd $(GO_PATH)/src && go-bindata -dev -prefix $(GO_PATH)/src -o $(GO_BINDATA_OUTPUT) -pkg $(GO_BINDATA_PACKAGE) -ignore $(GO_BINDATA_IGNORE) $(GO_BINDATA_PATHS)

run: bin
	cd ../gonode-skeleton && go run main.go server -config=./server.toml.dist

test: bin test-backend test-frontend

test-backend: bin
	go test $(GONODE_CORE) $(GONODE_MODULES) ./test/modules
	go vet ./...

test-frontend:
	cd explorer && npm test

install: install-backend install-frontend

install-backend:
	go get golang.org/x/tools/cmd/goimports
	go get -u github.com/jteeuwen/go-bindata/...
	go list -f '{{range .Imports}}{{.}} {{end}}' ./... | xargs go get -v
	go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs go get -v
	git clone https://github.com/rande/gonode-skeleton.git $(GO_PATH)/src/github.com/rande/gonode-skeleton

install-frontend:
	cd explorer && npm install

update:
	go get -u all
	cd explorer && npm update

load:
	curl -XPOST http://localhost:2508/setup/uninstall && exit 0
	curl -XPOST http://localhost:2508/setup/install
	curl -XPOST http://localhost:2508/setup/data/load

build:
	cd ../gonode-skeleton && make build

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
