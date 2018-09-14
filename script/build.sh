#!/bin/bash
set -xe

DIR="$( cd "$( dirname "$0"  )" && pwd  )"

cd ${DIR}
source ./buildproto.sh

cd ${DIR}
source buildservice.sh
