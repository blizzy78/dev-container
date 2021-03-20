#!/bin/bash -xe

ZULU_VERSION="11.45.27"
JDK_VERSION="11.0.10"

curl -LO https://cdn.azul.com/zulu/bin/zulu${ZULU_VERSION}-ca-jdk${JDK_VERSION}-linux_x64.tar.gz
tar xzf zulu${ZULU_VERSION}-ca-jdk${JDK_VERSION}-linux_x64.tar.gz
rm zulu${ZULU_VERSION}-ca-jdk${JDK_VERSION}-linux_x64.tar.gz
ln -s zulu${ZULU_VERSION}-ca-jdk${JDK_VERSION}-linux_x64 jdk

for i in $(ls -1 jdk/bin) ; do
	sudo ln -s /home/$(whoami)/jdk/bin/${i} /usr/bin/${i}
done
