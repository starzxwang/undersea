#!/bin/bash

# 此脚本做部署测试用

# 删除掉之前的老进程
if [ $1 == "update" ]
then
  pids=`ps aux | awk '/\/bin\/im-api/ {print $2}'`
    for pid in $pids; do
      if [ $pid > 0 ]
      then
        kill -9 $pid
        fi
    done
  fi

if [ $1 == "kill" ]
then
  exit
  fi

# 准备建立新程序
go build -o ${GOPATH}/bin/im-api ./im-api/cmd

nohup ${GOPATH}/bin/im-api > ${GOPATH}/log/im-api.log 2>&1 &


