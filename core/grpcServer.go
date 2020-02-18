package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/commis/fabric-client-go/core/chaincode"
	cc "github.com/commis/fabric-client-go/core/common"
	gm "github.com/commis/fabric-client-go/grpc/message"
)

type GrpcServer struct {
	gm.ChaincodeServiceServer
}

func (t *GrpcServer) Send(ctx context.Context, in *gm.ChaincodeRequest) (*gm.ResultResponse, error) {
	var token = ""
	switch in.GetMethod() {
	case gm.ChaincodeRequest_QUERY:
		return t.queryCC(token, in)

	case gm.ChaincodeRequest_INVOKE:
		return t.invokeCC(token, in)
	default:
		return nil, fmt.Errorf("invalid grpc method: %v", in.GetMethod())
	}
}

func (t *GrpcServer) queryCC(token string, in *gm.ChaincodeRequest) (*gm.ResultResponse, error) {
	var oper cc.OperateInfo
	if err := json.Unmarshal([]byte(in.GetOperator()), &oper); err != nil {
		return t.setExecuteError(err.Error())
	}

	var blockChain cc.BlockChain
	if err := json.Unmarshal([]byte(in.GetChaincode()), &blockChain); err != nil {
		return t.setExecuteError(err.Error())
	}

	req := &cc.GrpcRequest{ChainCode: &blockChain, Operate: &oper, JsonArgs: in.GetArgs()}
	payload, err := chaincode.QueryChainCode(req)
	if err != nil {
		return t.setExecuteError(err.Error())
	}

	return t.setExecuteSuccess(payload)
}

func (t *GrpcServer) invokeCC(token string, in *gm.ChaincodeRequest) (*gm.ResultResponse, error) {
	var oper cc.OperateInfo

	if err := json.Unmarshal([]byte(in.GetOperator()), &oper); err != nil {
		return t.setExecuteError(err.Error())
	}

	var blockChain cc.BlockChain
	if err := json.Unmarshal([]byte(in.GetChaincode()), &blockChain); err != nil {
		return t.setExecuteError(err.Error())
	}

	req := &cc.GrpcRequest{ChainCode: &blockChain, Operate: &oper, JsonArgs: in.GetArgs()}
	payload, err := chaincode.InvokeChainCode(req)
	if err != nil {
		return t.setExecuteError(err.Error())
	}

	return t.setExecuteSuccess(payload)
}

func (t *GrpcServer) setExecuteError(message string) (*gm.ResultResponse, error) {
	err := fmt.Errorf("%s", message)
	return &gm.ResultResponse{Code: int32(cc.StatusExecuteFailed), Data: "", Message: message}, err
}

func (t *GrpcServer) setExecuteSuccess(payload interface{}) (*gm.ResultResponse, error) {
	return &gm.ResultResponse{Code: int32(cc.StatusExecuteSuccess), Data: payload.(string), Message: ""}, nil
}
