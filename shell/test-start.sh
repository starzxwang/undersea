#!/bin/bash

# 此脚本做部署测试用

# 删除掉之前的老进程
pathSuffixList=('\/bin\/im-api' '\/bin\/im-balance' '\/bin\/im-manage')

for pathSuffix in ${pathSuffixList[@]}; do
  pid=`ps aux | awk '/'$pathSuffix'/ {print $2}'`
  if [ $pid > 0 ]
  then
    kill -9 $pid
    fi
done

# 准备建立新程序
go build -o ${GOPATH}/bin/im-api ./im-api/cmd
go build -o ${GOPATH}/bin/im-balance ./im-balance/cmd
# go build -o ${GOPATH}/bin/im-manage ./im-manage/cmd

nohup ${GOPATH}/bin/im-api > ${GOPATH}/log/im-api.log 2>&1 &
nohup ${GOPATH}/bin/im-balance > ${GOPATH}/log/im-balance.log 2>&1 &


