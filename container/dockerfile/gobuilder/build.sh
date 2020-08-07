#!/bin/bash

# 输出编译变量
echo "源码目录 --> /go/src/$GO_PKG"
echo "编译包名 --> $GO_BUILD_PKG"

cd /go/src/$GO_PKG

out=`CGO_ENABLED=0 go build -o "binary" $GO_BUILD_PKG 2>&1 >/dev/null`

if [ $? -eq 0 ];then
  echo  -e  "\033[32m程序编译成功\033[0m"
  exit 1
else
  echo  -e  "\033[31m程序编译出错,请检查代码哦\033[0m"
  echo "$out"
  exit 2
fi

# 使用方法:
# gobuilder=gobuilder:1.8.3
# gopkg=app
#
# docker run --rm -v "$(pwd):/go/src/${gopkg}" -e GO_PKG=$gopkg -e GO_BUILD_PKG=$gopkg/cmd/xxx $gobuilder
# mv -f ./binary ./cmd/xxx