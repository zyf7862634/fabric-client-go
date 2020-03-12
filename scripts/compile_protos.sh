#!/bin/bash

#set -eux

#JAVA=$GOPATH/src/github.com/hyperledger/fabric-client-go/grpc

PROTO_ROOT_FILES=$(find . -path $PWD/vendor -prune -o -name ".protoroot" -exec readlink -f {} \;)
PROTO_ROOT_DIRS=$(dirname $PROTO_ROOT_FILES)

for dir in $PROTO_ROOT_DIRS; do
    echo "Working on dir $dir"
    # protoc --proto_path="$dir" --java_out=$JAVA --go_out=plugins=grpc:$GOPATH/src "$protos"/*.proto >& /dev/null
    protoc --proto_path="$dir" --go_out=plugins=grpc:$GOPATH/src "$dir"/*.proto >& /dev/null
done
