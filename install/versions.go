//go:build mage

package main

const (
	// https://golang.org/dl/
	goVersion = "1.17.3"

	// https://github.com/protocolbuffers/protobuf/releases
	protocVersion = "3.19.1"

	// https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/
	protocGenGRPCJavaVersion = "1.42.1"

	// https://github.com/nvm-sh/nvm/releases
	nvmVersion = "0.39.0"

	// https://www.azul.com/downloads/zulu-community/?version=java-11-lts&os=ubuntu&architecture=x86-64-bit&package=jdk
	zuluVersion    = "11.52.13"
	zuluJDKVersion = "11.0.13"

	// https://maven.apache.org/download.cgi
	mavenVersion = "3.8.4"

	// https://github.com/restic/restic/releases
	resticVersion = "0.12.1"
)

var nodeLTSNames = []string{"gallium", "dubnium"}
