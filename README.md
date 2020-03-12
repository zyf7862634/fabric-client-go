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
