# Install

mkdir -p $GOPATH/src/github.com/hyperledger

cd $GOPATH/src/github.com/hyperledger

git clone http://192.168.1.232/blockchain/fabric-sdk-go.git

mkdir -p $GOPATH/src/github.com/commis/fabtic-client-go

git clone http://192.168.1.232/blockchain/fabtic-client-go.git

# Compile

cd $GOPATH/src/github.com/commis/fabtic-client-go

make build

# Running

cd $GOPATH/src/github.com/commis/fabtic-client-go/test/server

./bin/start.sh

`Notes:`

    1.Configure your fabric network:
        test/server/etc/fabric.yaml

