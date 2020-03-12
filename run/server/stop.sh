#!/bin/bash

PID=$(ps -ef |grep httpserver |grep -v grep |awk '{print $2}')
if [ ! -z "${PID}" ]; then
    kill -9 ${PID}
fi
echo $(date "+%Y-%m-%d %H:%M:%S")" stop server success."
