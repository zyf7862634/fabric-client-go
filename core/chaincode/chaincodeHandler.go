package chaincode

import (
	"github.com/commis/fabric-client-go/core/common"
	"github.com/valyala/fasthttp"
)

func Query(ctx *fasthttp.RequestCtx) {
	req := &common.HttpRequest{}
	if err := common.ParseRequestBody(ctx, req); err != nil {
		return
	}

	payload, err := QueryChainCode(
		&common.GrpcRequest{
			ChainCode: &req.ChainCode,
			Operate:   &req.Operate,
			JsonArgs:  req.JsonArgs})
	common.ProcessOperateResult(&ctx.Response, payload, err)
}

func Invoke(ctx *fasthttp.RequestCtx) {
	req := &common.HttpRequest{}
	if err := common.ParseRequestBody(ctx, req); err != nil {
		return
	}

	payload, err := InvokeChainCode(
		&common.GrpcRequest{
			ChainCode: &req.ChainCode,
			Operate:   &req.Operate,
			JsonArgs:  req.JsonArgs})
	common.ProcessOperateResult(&ctx.Response, payload, err)
}
