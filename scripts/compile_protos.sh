#!/bin/bash

#set -eux

JAVA=$GOPATH/src/github.com/commis/fabric-client-go/grpc

# Find explicit proto roots
PROTO_ROOT_FILES="$(find . -name ".protoroot" -exec readlink -f {} \;)"
PROTO_ROOT_DIRS="$(dirname $PROTO_ROOT_FILES)"

# Find all proto files to be compiled, excluding any which are in a proto root or in the vendor folder, as well as the gotools dir
ROOTLESS_PROTO_FILES="$(find $PWD \
                            $(for dir in $PROTO_ROOT_DIRS ; do echo "-path $dir -prune -o " ; done) \
                            -path $PWD/vendor -prune -o \
                            -name "*.proto" -exec readlink -f {} \;)"
ROOTLESS_PROTO_DIRS="$(dirname $ROOTLESS_PROTO_FILES | sort | uniq)"

for dir in $ROOTLESS_PROTO_DIRS; do
echo Working on dir $dir
	for protos in $(find "$dir" -name '*.proto' -exec dirname {} \; | sort | uniq) ; do
        protoc --proto_path="$dir" --java_out=$JAVA --go_out=plugins=grpc:$GOPATH/src "$protos"/*.proto >& /dev/null
	done
done
