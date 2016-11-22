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
GO_PKG = ./core/bindata,./core/commands,./core/config,./core/guard,./core/helper,./core/logger,./core/router,./core/security,./core/squirrel,./core/vault,./modules/api,./modules/base,./modules/blog,./modules/debug,./modules/feed,./modules/guard,./modules/media,./modules/prism,./modules/raw,./modules/search,./modules/setup,./modules/user

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
	go get -u github.com/jteeuwen/go-bindata/...
	cd ../gonode-skeleton && make run

test: bin test-backend test-frontend  ## Run tests

test-backend: bin     ## Run backend tests
	go test -v $(GONODE_CORE) $(GONODE_MODULES) ./test/modules
	go vet ./...

test-frontend:        ## Run frontend tests
	cd explorer && npm test

install: install-backend install-frontend ## Install dependencies

coverage-backend: bin ## run coverage tests
	mkdir -p build/coverage && rm -rf build/coverage/*.cov
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/core_bindata.cov ./core/bindata
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/core_commands.cov ./core/commands
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/core_config.cov ./core/config
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/core_guard.cov ./core/guard
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/core_helper.cov ./core/helper
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/core_logger.cov ./core/logger
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/core_router.cov ./core/router
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/core_security.cov ./core/security
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/core_squirrel.cov ./core/squirrel
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/core_vault.cov ./core/vault
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/modules_api.cov ./modules/api
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/modules_base.cov ./modules/base
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/modules_blog.cov ./modules/blog
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/modules_debug.cov ./modules/debug
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/modules_feed.cov ./modules/feed
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/modules_guard.cov ./modules/guard
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/modules_media.cov ./modules/media
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/modules_prism.cov ./modules/prism
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/modules_raw.cov ./modules/raw
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/modules_search.cov ./modules/search
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/modules_setup.cov ./modules/setup
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/modules_user.cov ./modules/user
	go test -v -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/functionnals.cov ./test/modules
	gocovmerge build/coverage/* > build/gonode.coverage
	go tool cover -html=./build/gonode.coverage -o build/gonode.html

install-backend: ## Install backend dependencies
	go get github.com/wadey/gocovmerge
	go get golang.org/x/tools/cmd/cover
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


dkr-test:
	docker-compose exec back make test

dkr-run:
	docker-compose exec back make run

dkr-back:
	docker-compose exec back /bin/bash

dkr-front:
	docker-compose exec front /bin/bash

dkr-watch:
	docker-compose exec front ./node_modules/.bin/webpack-dev-server --config webpack-dev-server.config.js --progress --inline --colors