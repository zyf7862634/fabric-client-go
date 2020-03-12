package common

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

func Index(ctx *fasthttp.RequestCtx) {
	WriteSuccessResponse(&ctx.Response, nil, "Welcome")
}

func Enroll(ctx *fasthttp.RequestCtx) {
	req := &HttpRequest{}
	if err := ParseRequestBody(ctx, req); err != nil {
		WriteFailedResponse(&ctx.Response, StatusExecuteFailed, err.Error())
		return
	}

	userName, err := GetFabricSetupIns().GetUserNameFromCert(req.Operate.User)
	if err != nil {
		WriteFailedResponse(&ctx.Response, StatusExecuteFailed, err.Error())
		return
	}

	userCert, result := GetFabricSetupIns().RegisteredAndEnrollUser(userName)
	if !result {
		WriteFailedResponse(&ctx.Response, StatusExecuteFailed, userCert)
	} else {
		WriteSuccessResponse(&ctx.Response, nil, userCert)
	}
}

func ParseRequestBody(ctx *fasthttp.RequestCtx, req *HttpRequest) error {
	body := ctx.PostBody()
	if err := json.Unmarshal(body, req); err != nil {
		logger.Errorf("Parse error: %v, %s", err, body)
		message := "Unprocessed Entity"
		WriteFailedResponse(&ctx.Response, http.StatusUnprocessableEntity, message)
		return err
	}

	svrCfg := GetSvrConfigIns()
	if req.ChainCode.Channel == "" || req.ChainCode.ChainCodeName == "" {
		req.ChainCode.Channel = Setup.ChannelID
		req.ChainCode.ChainCodeName = svrCfg.GetCfgString(FabricChainCodeName)
	}

	return nil
}

func ProcessOperateResult(w *fasthttp.Response, txId interface{}, payload interface{}, err error) {
	if err != nil {
		WriteFailedResponse(w, StatusExecuteFailed, err.Error())
	} else {
		WriteSuccessResponse(w, txId, payload)
	}
}

func WriteSuccessResponse(w *fasthttp.Response, txId interface{}, m interface{}) {
	jsonResponse := JsonResponse{Code: StatusExecuteSuccess, TxID: txId, Data: m, Message: nil}

	if err := json.NewEncoder(w.BodyWriter()).Encode(&jsonResponse); err != nil {
		WriteFailedResponse(w, StatusInternalServerError, "Internal Server Error")
	} else {
		setResponseHeader(w)
	}
}

func WriteFailedResponse(w *fasthttp.Response, errorCode int, errorMsg string) {
	jsonResponse := JsonResponse{Code: errorCode, TxID: nil, Data: nil, Message: errorMsg}

	if err := json.NewEncoder(w.BodyWriter()).Encode(&jsonResponse); err != nil {
		logger.Errorf("write error response failed, %v", jsonResponse)
	}

	logger.Errorf("Error: %s", w.Body())
	setResponseHeader(w)
}

func setResponseHeader(w *fasthttp.Response) {
	w.Header.SetContentType("application/json; charset=utf8")
	w.Header.SetStatusCode(http.StatusOK)
}
