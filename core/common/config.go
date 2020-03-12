package common

import (
	"fmt"
	"github.com/commis/fabric-client-go/utils"
	"os"
	"strings"
	"sync"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var (
	svrCfgMutex sync.Once
	svrCfg      *ServerConfigure
)

type ServerConfigure struct {
	cfgPath  string
	cfgViper *viper.Viper
}

func GetSvrConfigIns() *ServerConfigure {
	svrCfgMutex.Do(func() {
		svrCfg = &ServerConfigure{}
		svrCfg.cfgPath = utils.GetCurrentExeFileDir() + "/etc"
		_, err := os.Stat(svrCfg.cfgPath)
		if err != nil {
			if !os.IsExist(err) {
				fmt.Println("change dir for develop configure.")
				svrCfg.cfgPath = "/home/developWork/DevProjects/src/github.com/hyperledger/fabric-client-go/run/server/etc"
			}
		}
		svrCfg.initConfig()
	})
	return svrCfg
}

func SetupModuleLogger(module string) *logging.Logger {
	serverConfig := GetSvrConfigIns()
	logLevel := serverConfig.cfgViper.GetString(ServerLogLevel)
	level, err := logging.LogLevel(logLevel)
	if err != nil {
		level = DefaultLogLevel
	}
	moduleName := serverConfig.cfgViper.GetString(ServerName) + "." + module
	logging.SetLevel(level, moduleName)
	return logging.MustGetLogger(moduleName)
}

func (sf *ServerConfigure) initConfig() {
	sf.cfgViper = viper.New()
	sf.cfgViper.SetEnvPrefix("svr")
	sf.cfgViper.AutomaticEnv()
	sf.cfgViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	sf.cfgViper.SetConfigFile(sf.GetCfgFile(ServerConfig))

	if err := sf.cfgViper.ReadInConfig(); err != nil {
		fmt.Printf("loading config faile failed, %v", err)
		os.Exit(1)
	}
	sf.setLoggerFormat()
}

func (sf *ServerConfigure) setLoggerFormat() {
	logFormat := sf.cfgViper.GetString(ServerLogFormat)
	format := logging.MustStringFormatter(logFormat)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(logging.SetBackend(backendFormatter))
}

func (sf *ServerConfigure) GetCfgFile(fileName string) string {
	return sf.cfgPath + "/" + fileName
}

func (sf *ServerConfigure) GetCfgString(key string) string {
	return sf.cfgViper.GetString(key)
}

func (sf *ServerConfigure) GetCfgInt(key string) int {
	return sf.cfgViper.GetInt(key)
}

func (sf *ServerConfigure) GetCfgBool(key string) bool {
	return sf.cfgViper.GetBool(key)
}
