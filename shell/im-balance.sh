#!/bin/bash

# 删除掉之前的老进程
if [ $1 == "update" ]
then
  pid=`ps aux | awk '/\/bin\/im-balance/ {print $2}'`
  if [ $pid > 0 ]
  then
    kill -9 $pid
    fi
fi

if [ $1 == "kill" ]
then
  exit
  fi

# 准备建立新程序
go build -o ${GOPATH}/bin/im-balance ./im-balance/cmd

nohup ${GOPATH}/bin/im-balance > ${GOPATH}/log/im-balance.log 2>&1 &


