package chaincode

import (
	"fmt"
	"github.com/commis/fabric-client-go/core/common"
	"github.com/commis/fabric-client-go/utils"
	"github.com/commis/fabric-sdk-go/pkg/client/channel"
)

const (
	RequestArgsTypeArray int = 0
	RequestArgsTypeJson  int = 1
)

var logger = common.SetupModuleLogger("core.chaincode")

func getArgBytes(param *common.ParamInfo) [][]byte {
	switch param.Type {
	case RequestArgsTypeArray:
		return utils.GetArrayArgBytes(param.Value)
	case RequestArgsTypeJson:
		return utils.GetJsonArgBytes(param.Value)
	default:
	}
	return nil
}

func QueryChainCode(req *common.HttpRequest) (interface{}, error) {
	logger.Debugf("request: %v", *req)

	hfClient := common.GetFabricSetupIns().Client
	args := getArgBytes(&req.Parameter)
	if args == nil {
		return nil, fmt.Errorf("invalid argument type %d when query chaincode", req.Parameter.Type)
	}

	request := channel.Request{
		ChaincodeID: req.ChainCode.ChainCodeName,
		Fcn:         req.Operate.String(),
		Args:        args}
	resp, err := hfClient.Query(request)
	if err != nil {
		return nil, err
	}

	return string(resp.Payload), nil
}

func InvokeChainCode(req *common.HttpRequest) (interface{}, interface{}, error) {
	logger.Debugf("request: %v", *req)

	hfClient := common.GetFabricSetupIns().Client
	args := getArgBytes(&req.Parameter)
	if args == nil {
		return nil, nil, fmt.Errorf("invalid argument type %d when invoke chaincode", req.Parameter.Type)
	}

	request := channel.Request{
		ChaincodeID: req.ChainCode.ChainCodeName,
		Fcn:         req.Operate.String(),
		Args:        args}
	response, err := hfClient.Execute(request)
	if err != nil {
		return nil, nil, err
	}

	return response.TransactionID, string(response.Payload), nil
}
