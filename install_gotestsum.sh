#!/bin/bash
os=$(uname)
arch=$(uname -m)

if [[ $(which gotestsum) == *"not found"* ]]
then
	if [ $os = Darwin ]
	then
		if [ $arch = "arm64" ]
		then
			curl -sSL "https://github.com/gotestyourself/gotestsum/releases/download/v1.8.1/gotestsum_1.8.1_darwin_arm64.tar.gz" | sudo tar -xz -C /usr/local/bin gotestsum 
		else
			curl -sSL "https://github.com/gotestyourself/gotestsum/releases/download/v1.8.1/gotestsum_1.8.1_darwin_amd64.tar.gz" | sudo tar -xz -C /usr/local/bin gotestsum 
		fi
	else
		curl -sSL "https://github.com/gotestyourself/gotestsum/releases/download/v1.8.1/gotestsum_1.8.1_linux_amd64.tar.gz" | sudo tar -xz -C /usr/local/bin gotestsum 
	fi
else
	echo "skipping install"
fi