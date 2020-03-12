package chaincode

import (
	"github.com/commis/fabric-client-go/core/common"
	"github.com/valyala/fasthttp"
)

func Query(ctx *fasthttp.RequestCtx) {
	req := common.HttpRequest{}
	if err := common.ParseRequestBody(ctx, &req); err != nil {
		return
	}

	payload, err := QueryChainCode(&req)
	common.ProcessOperateResult(&ctx.Response, nil, payload, err)
}

func Invoke(ctx *fasthttp.RequestCtx) {
	req := common.HttpRequest{}
	if err := common.ParseRequestBody(ctx, &req); err != nil {
		return
	}

	txId, payload, err := InvokeChainCode(&req)
	common.ProcessOperateResult(&ctx.Response, txId, payload, err)
}
