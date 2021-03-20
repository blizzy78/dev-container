#!/bin/bash -xe

for i in .bash_history .gitconfig ; do
	mkdir ${i}_dir
	ln -s ${i}_dir/${i} ${i}
done

echo "Europe/Berlin" |sudo tee /etc/timezone
sudo env DEBIAN_FRONTEND=noninteractive apt install tzdata

echo "export PATH=\"/home/$(whoami)/bin:\$PATH\"" >>~/.bashrc
