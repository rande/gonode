#!/usr/bin/env bash

GO_BINDATA_PATHS="${GOPATH}/src/github.com/rande/gonode/modules/... ${GOPATH}/src/github.com/rande/gonode/explorer/dist/..."
GO_BINDATA_IGNORE="(.*)\.(go|DS_Store)"
GO_BINDATA_OUTPUT="${GOPATH}/src/github.com/rande/gonode/assets/bindata.go"
GO_BINDATA_PACKAGE="assets"

cd ${GOPATH}/src && go-bindata -dev -prefix ${GOPATH}/src -o ${GO_BINDATA_OUTPUT} -pkg ${GO_BINDATA_PACKAGE} -ignore ${GO_BINDATA_IGNORE} ${GO_BINDATA_PATHS}