package common

import "strings"

type OperateInfo struct {
	User   string `json:"user,omitempty"`
	Module string `json:"module,omitempty"`
	Method string `json:"func"`
}

type ParamInfo struct {
	Type  int    `json:"type"` // 参数类别：1 数组、2 JSON
	Value string `json:"args"` // 实际参数直
}

type JsonResponse struct {
	Code    int         `json:"code"` /*返回错误码*/
	TxID    interface{} `json:"txID"` /*交易Hash*/
	Data    interface{} `json:"data"`
	Message interface{} `json:"message"`
}

type HttpRequest struct {
	Token     string      `json:"token,omitempty"` // 用来Restful服务的鉴权
	ChainCode BlockChain  `json:"chaincode,omitempty"`
	Operate   OperateInfo `json:"operator"`
	Parameter ParamInfo   `json:"p,omitempty"`
}

type BlockChain struct {
	Channel       string `json:"channel,omitempty"`
	ChainCodeName string `json:"ccName,omitempty"`
	ChainCodePath string `json:"ccPath,omitempty"`
	ChainCodeVer  string `json:"version,omitempty"`
}

func (t *OperateInfo) String() string {
	var data []string
	if t.User != "" {
		data = append(data, t.User)
	}
	if t.Module != "" {
		data = append(data, t.Module)
	}
	if t.Method != "" {
		data = append(data, t.Method)
	}
	return strings.Join(data, ":")
}
