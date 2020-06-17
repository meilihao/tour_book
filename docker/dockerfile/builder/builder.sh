#!/usr/bin/env bash
[ "`docker version 2>&1 | grep 'Go version'`" = "" ] && echo -e "\n\n no docker.\n\n" && exit 1

name="builder-xxx"
tag=`sha256sum build/Dockerfile | cut -c 1-10`

if [ "`docker images |grep \"${name}\" | grep \"{tag}\"`" = "" ]; then
	echo -e "\n\nPlease create builder image with the following command.\n\n"
	echo -e "\tdocker build -f `pwd`/build/Dockerfile -t ${name}:${tag} `pwd`/build\n\n"
	exit 1
if

docker run -v `pwd`:/app -it ${name}:${tag} xxx