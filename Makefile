.PHONY: test run explorer build help

PID = .pid
GO_FILES = $(shell find . -type f -name "*.go")
GONODE_MODULES = $(shell ls -d ./modules/* | grep -v go)
GONODE_CORE = $(shell ls -d ./core/* | grep -v go)
GO_PATH = $(shell go env GOPATH)
GO_BINDATA_PATHS = $(GO_PATH)/src/github.com/rande/gonode/modules/... $(GO_PATH)/src/github.com/rande/gonode/explorer/dist/...
GO_BINDATA_IGNORE = "(.*)\.(go|DS_Store)"
GO_BINDATA_OUTPUT = $(shell pwd)/assets/bindata.go
GO_BINDATA_PACKAGE = assets

help:     ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

default: test

format:  ## Format code to respect CS
	goimports -w $(GO_FILES)
	gofmt -l -w -s .
	go fix ./...
	go vet ./...

bin:                 ## Generate bin assets file
	go get -u github.com/jteeuwen/go-bindata/...
	cd $(GO_PATH)/src && go-bindata -dev -prefix $(GO_PATH)/src -o $(GO_BINDATA_OUTPUT) -pkg $(GO_BINDATA_PACKAGE) -ignore $(GO_BINDATA_IGNORE) $(GO_BINDATA_PATHS)

run:               ## Run server
	cd ../gonode-skeleton && make run

test: bin test-backend test-frontend  ## Run tests

test-backend: bin     ## Run backend tests
	go test $(GONODE_CORE) $(GONODE_MODULES) ./test/modules
	go vet ./...

test-frontend:        ## Run frontend tests
	cd explorer && npm test

install: install-backend install-frontend ## Install dependencies

install-backend:  ## Install backend dependencies
	go get golang.org/x/tools/cmd/goimports
	go get -u github.com/jteeuwen/go-bindata/...
	go list -f '{{range .Imports}}{{.}} {{end}}' ./... | xargs go get -v
	go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs go get -v
	git clone https://github.com/rande/gonode-skeleton.git $(GO_PATH)/src/github.com/rande/gonode-skeleton || exit 0

install-frontend: ## Install frontend dependencies
	cd explorer && npm install

update:  ## Update dependencies
	go get -u all
	cd explorer && npm update

load:    ## Load fixtures
	curl -XPOST http://localhost:2508/setup/uninstall && exit 0
	curl -XPOST http://localhost:2508/setup/install
	curl -XPOST http://localhost:2508/setup/data/load

build:   ## Build final binary
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

build-assets:
	cd ../gonode/explorer && npm run-script build
