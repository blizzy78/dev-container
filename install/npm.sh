#!/bin/bash -xe

NVM_VERSION="0.37.2"

curl -Lo- https://raw.githubusercontent.com/nvm-sh/nvm/v${NVM_VERSION}/install.sh | bash
export NVM_DIR="/home/$(whoami)/.nvm"
. ${NVM_DIR}/nvm.sh
nvm install node
curl -Lo- https://www.npmjs.com/install.sh |bash
for i in npm npx ; do
	sudo ln -s $(which ${i}) /usr/bin/${i}
done
