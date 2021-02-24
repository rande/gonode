.PHONY: test run explorer build help install

PID = .pid
GO_FILES = $(shell find . -type f -name "*.go")
GONODE_MODULES = $(shell ls -d ./modules/* | grep -v go)
GONODE_CORE = $(shell ls -d ./core/* | grep -v go)
GOPATH = $(shell go env GOPATH)

install:
	$(call back,glide install)
	$(call back,go get github.com/wadey/gocovmerge && go get golang.org/x/tools/cmd/cover && go get golang.org/x/tools/cmd/goimports && go get -u github.com/jteeuwen/go-bindata/...)

test:
	./app/assets/bindata.sh
	mkdir -p data
	echo "mode: atomic" > data/coverage.out

	GONODE_TEST_OFFLINE=true GOPATH=${GOPATH} go test -v -failfast -covermode=atomic -coverprofile=data/coverage_core.out $(GONODE_CORE)
	GOPATH=${GOPATH} go test -v -failfast -covermode=atomic -coverprofile=data/coverage_modules.out $(GONODE_MODULES)
	GOPATH=${GOPATH} go test -v -failfast -covermode=atomic -coverpkg ./... -coverprofile=data/coverage_integration.out ./test/modules
	go vet $(GONODE_CORE) $(GONODE_MODULES) ./test/modules/

	tail -n +2 data/coverage_core.out >> data/coverage.out
	tail -n +2 data/coverage_modules.out >> data/coverage.out
	tail -n +2 data/coverage_integration.out >> data/coverage.out
	go tool cover -html data/coverage.out -o data/coverage.html

run:
	GOPATH=${GOPATH} `go env GOPATH`/bin/modd

format:
	gofmt -w $(GONODE_CORE) $(GONODE_MODULES) ./test/modules
	go fix $(GONODE_CORE) $(GONODE_MODULES) ./test/modules
	go vet $(GONODE_CORE) $(GONODE_MODULES) ./test/modules

load:    ## Load fixtures
	curl -XPOST http://localhost:2508/setup/uninstall && exit 0
	curl -XPOST http://localhost:2508/setup/install
	curl -XPOST http://localhost:2508/setup/data/load
