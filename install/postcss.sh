#!/bin/bash -xe

export NVM_DIR="/home/$(whoami)/.nvm"
. ${NVM_DIR}/nvm.sh
npm install -g postcss@latest postcss-cli@latest
sudo ln -s $(which postcss) /usr/bin/postcss
