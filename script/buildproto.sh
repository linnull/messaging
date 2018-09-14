#!/bin/bash
set -xe
DIR="$( cd "$( dirname "$0"  )" && pwd  )"
export GOPATH="$( cd ${DIR}/.. && pwd )"
cd ${GOPATH}/src
find . ! -path "./vendor/*" ! -path "./template/*" -name "*.proto" | xargs -t -I{} protoc -I. -I${GOPATH}/src --micro_out=. --go_out=. {}