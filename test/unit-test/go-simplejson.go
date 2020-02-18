package main

import (
	"encoding/json"
	"fmt"
	"github.com/commis/fabric-client-go/core/common"
)

func main() {
	var req common.HttpRequest
	req.Operate.Method = "query"
	req.JsonArgs = "{\"a\"}"

	//struct 到json str
	if b, err := json.Marshal(req); err == nil {
		fmt.Println("================struct 到json str==")
		fmt.Println(string(b))
	}
}
