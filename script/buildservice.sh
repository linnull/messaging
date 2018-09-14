#!/bin/bash
set -xe
export CGO_ENABLED=0
export GOOS=linux
export GOOS=windows
export GOARCH=amd64

DIR="$( cd "$( dirname "$0"  )" && pwd  )"
export GOPATH="$( cd ${DIR}/.. && pwd )"

cd ${DIR}
rm -rf msgbin
mkdir msgbin
cd msgbin
now=`date "+%Y-%m-%d/%H-%M-%S"`
commit=`git rev-parse HEAD`
branch=`git symbolic-ref --short -q HEAD`

if [[ ${GOOS} == 'windows' ]]; then
	EXE=.exe
fi

go build -o msggateway${EXE} $GOPATH/src/gateway/main.go