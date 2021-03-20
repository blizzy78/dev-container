#!/bin/bash -xe

PROTOC_VERSION="3.15.6"
PROTOC_JAVA_VERSION="1.36.0"

curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/protoc-$PROTOC_VERSION-linux-x86_64.zip
mkdir protoc-$PROTOC_VERSION
cd protoc-$PROTOC_VERSION && unzip ../protoc-$PROTOC_VERSION-linux-x86_64.zip && cd ..
rm protoc-$PROTOC_VERSION-linux-x86_64.zip
ln -s protoc-$PROTOC_VERSION protoc
sudo ln -s /home/$(whoami)/protoc/bin/protoc /usr/bin/protoc

export PATH="${PATH}:/go/bin"

go get \
	google.golang.org/protobuf/cmd/protoc-gen-go \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc \
	github.com/yoheimuta/protolint/cmd/protolint

curl -LO https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/$PROTOC_JAVA_VERSION/protoc-gen-grpc-java-$PROTOC_JAVA_VERSION-linux-x86_64.exe
mv protoc-gen-grpc-java-$PROTOC_JAVA_VERSION-linux-x86_64.exe protoc/bin/
ln -s protoc-gen-grpc-java-$PROTOC_JAVA_VERSION-linux-x86_64.exe protoc/bin/protoc-gen-grpc-java
sudo ln -s /home/$(whoami)/protoc/bin/protoc-gen-grpc-java /usr/bin/protoc-gen-grpc-java
