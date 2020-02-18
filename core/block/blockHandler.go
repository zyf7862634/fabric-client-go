package block

import (
	"github.com/commis/fabric-client-go/core/common"
	"github.com/valyala/fasthttp"
)

func BlockHandler(ctx *fasthttp.RequestCtx) {
	req := &common.HttpRequest{}
	if err := common.ParseRequestBody(ctx, req); err != nil {
		return
	}

	payload, err := QueryBlock(req)
	common.ProcessOperateResult(&ctx.Response, payload, err)
}

func BlockTxHandler(ctx *fasthttp.RequestCtx) {
	req := &common.HttpRequest{}
	if err := common.ParseRequestBody(ctx, req); err != nil {
		return
	}

	payload, err := QueryTransaction(req)
	common.ProcessOperateResult(&ctx.Response, payload, err)
}
