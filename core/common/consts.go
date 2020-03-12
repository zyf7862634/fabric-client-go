package common

import "github.com/op/go-logging"

const (
	FabricCliConfig     = "fabric.yaml"
	FabricChainCodeName = "fabric.ccName"
	FabricOrgUser       = "fabric.orgUser"
	FabricOrgName       = "fabric.orgName"
	FabricOrgChannel    = "fabric.channelName"
	FabricAffiliation   = "fabric.affiliation"
	FabricCertCaName    = "fabric.caName"
	FabricIdentityType  = "user"

	ServerConfig          = "server.yaml"
	ServerName            = "server.name"
	ServerHttpPort        = "server.listen.http"
	ServerGrpsPort        = "server.listen.grps"
	ServerDebugPort       = "server.listen.debug"
	ServerLogLevel        = "server.logging.level"
	ServerLogFormat       = "server.logging.format"
	ServerCacheExpiration = "server.cache.expiration"
	ServerCacheGCInterval = "server.cache.gcInterval"
	ServerUserCertOrName  = "server.user.cert"

	DefaultLogLevel = logging.INFO
)

const (
	StatusExecuteSuccess      = 200
	StatusInternalServerError = 500
	StatusExecuteFailed       = 1001 // Execute failed
)
