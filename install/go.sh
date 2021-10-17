#!/bin/bash -e

# https://golang.org/dl/
GO_VERSION="1.17.2"

curl -LO https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz
tar xzf go${GO_VERSION}.linux-amd64.tar.gz
rm go${GO_VERSION}.linux-amd64.tar.gz
mv go go-${GO_VERSION}
ln -s go-${GO_VERSION} go
sudo ln -s /home/$(whoami)/go /go
sudo ln -s /go/bin/go /usr/bin/go
echo "export PATH=\"\$PATH:/go/bin\"" >>~/.bashrc
