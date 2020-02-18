package core

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/commis/fabric-client-go/core/chaincode"

	"github.com/commis/fabric-client-go/core/block"
	"github.com/commis/fabric-client-go/core/common"
)

var logger = common.SetupModuleLogger("core")

func NewFastHttpRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.GET("/", common.Index)
	router.POST("/enroll", common.Enroll)

	router.POST("/query", chaincode.Query)
	router.POST("/invoke", chaincode.Invoke)

	router.POST("/block", block.BlockHandler)
	router.POST("/tx", block.BlockTxHandler)

	return router
}
