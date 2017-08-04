.PHONY: test run explorer build help install

PID = .pid
GO_FILES = $(shell find . -type f -name "*.go")
GONODE_MODULES = $(shell ls -d ./modules/* | grep -v go)
GONODE_CORE = $(shell ls -d ./core/* | grep -v go)
GO_PATH = $(shell go env GOPATH)

help:     ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

define back
    $(call docker,back,$(1))
endef

define docker
    if [ ! -f /.dockerenv ]; then \
        docker-compose run $(1) /bin/bash -c "$(2)"; \
    else \
        /bin/bash -c "$(2)"; \
    fi
endef

shell:
	docker-compose run back /bin/bash

install:
	mkdir -p runtime/src
	$(call back,glide install)
	$(call back,go get github.com/wadey/gocovmerge && go get golang.org/x/tools/cmd/cover && go get golang.org/x/tools/cmd/goimports && go get -u github.com/jteeuwen/go-bindata/...)
	$(call back,cp -rvf /usr/local/go/src/* ./runtime/src/ )

test:
	$(call back,./app/assets/bindata.sh && go test -v $(GONODE_CORE) $(GONODE_MODULES) ./test/modules)
	$(call back,go vet $(GONODE_CORE) $(GONODE_MODULES) ./test/modules/)

format:
	$(call back,gofmt -w $(GONODE_CORE) $(GONODE_MODULES) ./test/modules)
	$(call back,go fix $(GONODE_CORE) $(GONODE_MODULES) ./test/modules)
	$(call back,go vet $(GONODE_CORE) $(GONODE_MODULES) ./test/modules)

run:
	docker-compose kill
	docker-compose up

load:    ## Load fixtures
	curl -XPOST http://localhost:2508/setup/uninstall && exit 0
	curl -XPOST http://localhost:2508/setup/install
	curl -XPOST http://localhost:2508/setup/data/load