#!/usr/bin/env bash

GOPATH=`go env GOPATH`
#GO_BINDATA_PATHS="${GOPATH}/src/github.com/rande/gonode/modules/..."
if [ -d ${GOPATH}/src/github.com/rande/gonode/modules ]; then
    GO_BINDATA_PATHS="${GOPATH}/src/github.com/rande/gonode/modules/..."
    GO_BINDATA_OUTPUT="${GOPATH}/src/github.com/rande/gonode/app/assets/bindata.go"
    GO_BINDATA_PREFIX="${GOPATH}/src/github.com/rande/gonode"
else
    GO_BINDATA_PATHS="./modules/..."
    GO_BINDATA_OUTPUT="./app/assets/bindata.go"
    GO_BINDATA_PREFIX=`pwd`
fi

GO_BINDATA_IGNORE="(.*)\.(go|DS_Store|jpg)"
GO_BINDATA_PACKAGE="assets"


echo "Generating bindata file..."
echo "GO_BINDATA_PATHS=${GO_BINDATA_PATHS}"
echo "GO_BINDATA_IGNORE=${GO_BINDATA_IGNORE}"
echo "GO_BINDATA_OUTPUT=${GO_BINDATA_OUTPUT}"
echo "GO_BINDATA_PACKAGE=${GO_BINDATA_PACKAGE}"
echo "GO_BINDATA_PREFIX=${GO_BINDATA_PREFIX}"

${GOPATH}/bin/go-bindata \
    -debug \
    -prefix ${GO_BINDATA_PREFIX}/ \
    -o ${GO_BINDATA_OUTPUT} \
    -pkg ${GO_BINDATA_PACKAGE} \
    -ignore ${GO_BINDATA_IGNORE} \
    ${GO_BINDATA_PATHS}

echo "Done!"