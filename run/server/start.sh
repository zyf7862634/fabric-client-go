#!/bin/bash

curr_dir=`pwd`
PROJECT_DIR=$(cd `dirname $(readlink -f "$0")`; pwd)
LOG_DIR=${PROJECT_DIR}/logs
BIN_EXE=${PROJECT_DIR}/bin/httpserver

if [ ! -d "${LOG_DIR}" ]; then
    mkdir -p ${LOG_DIR}
fi

if [ -d "${LOG_DIR}" ]; then
    rm -rf ${LOG_DIR}/*
fi

if [ -f "${BIN_EXE}" ]; then
    ${BIN_EXE} > ${LOG_DIR}/server.log 2>&1 &
    echo $(date "+%Y-%m-%d %H:%M:%S")" start server success."
fi

cd ${curr_dir}
