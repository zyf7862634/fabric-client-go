## Install
- mkdir -p $GOPATH/src/github.com/commis
- cd $GOPATH/src/github.com/commis
- git clone https://github.com/commis/fabric-client-go.git

## Compile
- cd $GOPATH/src/github.com/commis/fabtic-client-go
- make build

## Running
- $ cd $GOPATH/src/github.com/commis/fabtic-client-go/run/server
- ./bin/start.sh or ./httpserver

`Notes:`
- Configure your fabric network:
  - run/server/etc/fabric.yaml

## Using
<pre>
Query:
    Restful send body:
    {"operator":{"func":"query"},"p":{"type":0,"args":"a"}}

Invoke:
    Restful send body:
    {"operator":{"func":"invoke"},"p":{"type":0,"args":"a,b,10"}}
    {"operator":{"func":"invoke"},"p":{"type":1,"args":"{\"name\":\"test\",\"value\":1000}"}}
</pre>

## Table of Features
- Restful server support, argument support JSON string or string array
- Connect to fabric network, support query and invoke operation
