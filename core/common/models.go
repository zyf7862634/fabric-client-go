package common

type OperateInfo struct {
	User   string `json:"user,omitempty"`
	Module string `json:"module,omitempty"`
	Method string `json:"func"`
}

type JsonResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message interface{} `json:"message"`
}

type GrpcRequest struct {
	ChainCode *BlockChain
	Operate   *OperateInfo
	JsonArgs  string
}

type HttpRequest struct {
	Token     string      `json:"token,omitempty"`
	ChainCode BlockChain  `json:"chaincode,omitempty"`
	Operate   OperateInfo `json:"operator"`
	JsonArgs  string      `json:"args,omitempty"`
}

type BlockChain struct {
	Channel       string `json:"channel,omitempty"`
	ChainCodeName string `json:"ccName,omitempty"`
	ChainCodePath string `json:"ccPath,omitempty"`
	ChainCodeVer  string `json:"version,omitempty"`
}

func (t *OperateInfo) String() string {
	//todo::后续正式版本恢复
	/*data := []string{t.User, t.Module, t.Method}
	return strings.Join(data, "#")*/
	return t.Method
}
