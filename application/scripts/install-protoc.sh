#!/bin/bash

PROTOC_ZIP=protoc-31.1-linux-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v31.1/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc 'include/*'
rm -f $PROTOC_ZIP