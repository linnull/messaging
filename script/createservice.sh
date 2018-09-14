#!/bin/bash
set -xe
DIR="$( cd "$( dirname "$0"  )" && pwd  )"
SrcDir=${DIR}/../src
cd ${SrcDir}

ServiceName=$1
ServiceNameLower=`tr '[A-Z]' '[a-z]' <<< ${ServiceName}`
ServiceDir=${SrcDir}/${ServiceNameLower}

echo ${ServiceName}
echo ${ServiceNameLower}
echo ${ServiceDir}

#if [ -d ${ServiceDir} ]; then
#    echo "[Fail] ${ServiceDir} already exist!"
#    exit
#fi

cp ${SrcDir}/template ${ServiceDir} -r
cd ${ServiceDir}
find . -type f | xargs sed -i "s/Template/${ServiceName}/g"
find . -type f | xargs sed -i "s/template/${ServiceNameLower}/g"
find . -type f | xargs rename "s/template/${ServiceNameLower} /"