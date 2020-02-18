#!/bin/bash

curr_dir=`pwd`
PROJECT_DIR=$(cd `dirname $(readlink -f "$0")`/..; pwd)
BIN_EXE=${PROJECT_DIR}/test/server/bin/httpserver

echo "start build ..."
cd ${PROJECT_DIR}
go build -o ${BIN_EXE} .

echo "finished ..."
cd ${curr_dir}
