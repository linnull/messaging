#!/bin/bash
set -xe
source ./buildproto.sh

now=`date "+%Y/%m/%d_%H:%M:%S"`
commit=`git rev-parse HEAD`
branch=`git symbolic-ref --short -q HEAD`

rm -rf msgbin
mkdir msgbin
cd msgbin

go build -o msggateway ${GOPATH}/src/gateway/main.go