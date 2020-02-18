package chaincode

import (
	"github.com/commis/fabric-client-go/core/common"
	"github.com/commis/fabric-client-go/utils"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

var logger = common.SetupModuleLogger("core.chaincode")

func QueryChainCode(req *common.GrpcRequest) (interface{}, error) {
	logger.Debugf("request: %v", *req)

	hfClient := common.GetFabricSetupIns().Client

	args := utils.GetArgBytes(req.JsonArgs)
	request := channel.Request{
		ChaincodeID: req.ChainCode.ChainCodeName,
		Fcn:         req.Operate.String(),
		Args:        args}
	resp, err := hfClient.Query(request, channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		logger.Errorf("failed to query funds : %s", err)
		return nil, err
	}

	return string(resp.Payload), nil
}

func InvokeChainCode(req *common.GrpcRequest) (interface{}, error) {
	logger.Debugf("request: %v", *req)

	hfClient := common.GetFabricSetupIns().Client

	args := utils.GetArgBytes(req.JsonArgs)
	request := channel.Request{
		ChaincodeID: req.ChainCode.ChainCodeName,
		Fcn:         req.Operate.String(),
		Args:        args}
	response, err := hfClient.Execute(request, channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		logger.Errorf("failed to invoke funds : %s", err)
		return nil, err
	}

	return response.TransactionID, nil
}
