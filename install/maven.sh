#!/bin/bash -xe

MAVEN_VERSION="3.6.3"

curl -LO https://mirror.netcologne.de/apache.org/maven/maven-3/${MAVEN_VERSION}/binaries/apache-maven-${MAVEN_VERSION}-bin.tar.gz
tar xzf apache-maven-${MAVEN_VERSION}-bin.tar.gz
rm apache-maven-${MAVEN_VERSION}-bin.tar.gz
ln -s apache-maven-${MAVEN_VERSION} maven
sudo ln -s /home/$(whoami)/maven/bin/mvn /usr/bin/mvn
