package block

import (
	"fmt"
	"github.com/commis/fabric-client-go/core/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	dbMutex sync.Once
	dbSetup *DbClientSetup
)

type DbClientSetup struct {
	DbCache *cache.Cache
}

func GetDbClientIns() *DbClientSetup {
	dbMutex.Do(func() {
		dbSetup = &DbClientSetup{}
		dbSetup.setupClientCache()
	})
	return dbSetup
}

func (dc *DbClientSetup) GetClient(channelName string) (*ledger.Client, error) {
	value, found := dbSetup.DbCache.Get(channelName)
	if !found {
		fs := common.GetFabricSetupIns()
		options := []fabsdk.ContextOption{fabsdk.WithUser(fs.OrgAdmin), fabsdk.WithOrg(fs.OrgName)}
		provider := fs.Sdk.ChannelContext(channelName, options...)

		dbCli, err := ledger.New(provider)
		if err != nil {
			return nil, fmt.Errorf("failed to create ledger client: %s", err)
		}
		value = dbCli
		dc.DbCache.Set(channelName, value, cache.DefaultExpiration)
	}
	return value.(*ledger.Client), nil
}

func (dc *DbClientSetup) setupClientCache() {
	svrCfg := common.GetSvrConfigIns()
	defaultExpiration, _ := time.ParseDuration(svrCfg.GetCfgString(common.ServerCacheExpiration))
	gcInterval, _ := time.ParseDuration(svrCfg.GetCfgString(common.ServerCacheGCInterval))

	dc.DbCache = cache.New(defaultExpiration, gcInterval)
	if _, err := dc.GetClient(svrCfg.GetCfgString(common.FabricOrgChannel)); err != nil {
		logger.Errorf("setup ledger client error: %v", err)
	}
}
