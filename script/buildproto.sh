#!/bin/bash
set -xe
DIR="$( cd "$( dirname "$0"  )" && pwd  )"
export GOPATH="$( cd ${DIR}/.. && pwd )"
cd ${GOPATH}/src
find . -name "*.proto" | xargs -t -I{} protoc -I. -I${GOPATH}/src --go_out=plugins=micro:. {}
cd ${GOPATH}