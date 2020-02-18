package main

import (
	"fmt"
	"github.com/commis/fabric-client-go/core"
	"github.com/commis/fabric-client-go/core/common"
	gm "github.com/commis/fabric-client-go/grpc/message"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
)

var logger = common.SetupModuleLogger("main")

func main() {
	common.GetFabricSetupIns()

	go startGrpsServer()
	go startHttpServer()

	startPerformanceDebug()
}

func startGrpsServer() {
	grpsPort := common.GetSvrConfigIns().GetCfgInt(common.ServerGrpsPort)
	if grpsPort != 0 {
		logger.Infof("start server grps port %d", grpsPort)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpsPort))
		if err != nil {
			logger.Fatalf("failed to grps port listen: %v", err)
		}

		s := grpc.NewServer()
		gm.RegisterChaincodeServiceServer(s, &core.GrpcServer{})
		reflection.Register(s)
		if err := s.Serve(lis); err != nil {
			logger.Fatalf("failed to start server: %v", err)
		}
	}
}

func startHttpServer() {
	httpPort := common.GetSvrConfigIns().GetCfgInt(common.ServerHttpPort)
	if httpPort != 0 {
		logger.Infof("start server http port %d", httpPort)
		router := core.NewFastHttpRouter()
		err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", httpPort), router.Handler)
		if err != nil {
			logger.Fatalf("failed to http port listen: %v", err)
		}
	}
}

func startPerformanceDebug() {
	debugPort := common.GetSvrConfigIns().GetCfgInt(common.ServerDebugPort)
	if debugPort != 0 {
		logger.Infof("start server debug pprof port %d", debugPort)
		err := http.ListenAndServe(fmt.Sprintf(":%d", debugPort), nil)
		if err != nil {
			logger.Fatalf("failed to debug pprof port listen: %v", err)
		}
	}
}
