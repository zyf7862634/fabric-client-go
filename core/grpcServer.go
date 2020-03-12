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
	switch in.GetMethod() {
	case gm.ChaincodeRequest_QUERY:
		return t.queryCC(in)
	case gm.ChaincodeRequest_INVOKE:
		return t.invokeCC(in)
	default:
		return nil, fmt.Errorf("invalid grpc method: %v", in.GetMethod())
	}
}

func (t *GrpcServer) getRequest(in *gm.ChaincodeRequest) (*cc.HttpRequest, error) {
	var oper cc.OperateInfo
	if err := json.Unmarshal([]byte(in.GetOperator()), &oper); err != nil {
		return nil, err
	}

	var blockChain cc.BlockChain
	if err := json.Unmarshal([]byte(in.GetChaincode()), &blockChain); err != nil {
		return nil, err
	}

	var param cc.ParamInfo
	if err := json.Unmarshal([]byte(in.GetArgs()), &param); err != nil {
		return nil, err
	}

	req := &cc.HttpRequest{ChainCode: blockChain, Operate: oper, Parameter: param}
	return req, nil
}

func (t *GrpcServer) queryCC(in *gm.ChaincodeRequest) (*gm.ResultResponse, error) {
	req, err := t.getRequest(in)
	if err != nil {
		return t.setExecuteError(err.Error())
	}

	payload, err := chaincode.QueryChainCode(req)
	if err != nil {
		return t.setExecuteError(err.Error())
	}

	return t.setExecuteSuccess(payload)
}

func (t *GrpcServer) invokeCC(in *gm.ChaincodeRequest) (*gm.ResultResponse, error) {
	req, err := t.getRequest(in)
	if err != nil {
		return t.setExecuteError(err.Error())
	}

	_, payload, err := chaincode.InvokeChainCode(req)
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
