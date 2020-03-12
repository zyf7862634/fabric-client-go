package block

import (
	"fmt"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"regexp"
	"strconv"
	"strings"

	cn "github.com/commis/fabric-client-go/core/common"
	"github.com/commis/fabric-client-go/utils"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
)

var logger = cn.SetupModuleLogger("core.block")

func QueryBlock(req *cn.HttpRequest) (interface{}, error) {
	logger.Debugf("request: %v", *req)

	dbClient, err := GetDbClientIns().GetClient(req.ChainCode.Channel)
	if err != nil {
		return nil, err
	}

	var target []string
	switch req.Operate.Method {
	case "height":
		return queryCurrentBlockHeight(dbClient, target)
	case "range":
		return queryBlockByRange(dbClient, target, req)
	default:
		return utils.NotSupportFunctionError(req.Operate.Method)
	}
}

func queryCurrentBlockHeight(cli *ledger.Client, target []string) (interface{}, error) {
	resp, err := cli.QueryInfo(ledger.WithTargetEndpoints(target...))
	if err != nil {
		message := fmt.Errorf("failed to query current block height: %s", err)
		logger.Errorf("%v", message)
		return nil, message
	}

	return *resp.BCI, nil
}

func queryBlockByNumber(cli *ledger.Client, blockNumber uint64, option ledger.RequestOption) *common.Block {
	block, err := cli.QueryBlock(blockNumber, option)
	if err == nil {
		if block.Data != nil {
			return block
		}
		logger.Errorf("failed to query block by height")
	}

	return nil
}

func queryBlockByRange(cli *ledger.Client, target []string, req *cn.HttpRequest) (interface{}, error) {
	resp, err := queryCurrentBlockHeight(cli, target)
	if err != nil {
		message := fmt.Errorf("failed to query block height: %s", err)
		logger.Errorf("%v", message)
		return nil, message
	}
	currHeight := resp.(common.BlockchainInfo).Height

	blockSlice := make([][][]byte, 0)
	if req.Parameter.Value != "" && !strings.Contains(req.Parameter.Value, "-") {
		requestOption := ledger.WithTargetEndpoints(target...)
		for _, id := range strings.Split(req.Parameter.Value, ",") {
			blockNumber, _ := strconv.ParseUint(id, 10, 64)
			if blockNumber >= currHeight {
				break
			}
			block := queryBlockByNumber(cli, blockNumber, requestOption)
			if block != nil {
				blockSlice = append(blockSlice, block.Data.GetData())
			}
		}
	} else {
		if min, max := getBlockRange(req.Parameter.Value); min != 0 && max != 0 {
			requestOption := ledger.WithTargetEndpoints(target...)
			for {
				if min > max || min >= currHeight {
					break
				} else {
					block := queryBlockByNumber(cli, min, requestOption)
					if block != nil {
						blockSlice = append(blockSlice, block.Data.GetData())
					}
				}
				min++
			}
		}
	}

	return blockSlice, nil
}

func QueryTransaction(req *cn.HttpRequest) (interface{}, error) {
	logger.Debugf("request: %v", *req)

	dbClient, err := GetDbClientIns().GetClient(req.ChainCode.Channel)
	if err != nil {
		return nil, err
	}

	var target []string
	resp, err := dbClient.QueryBlockByTxID(fab.TransactionID(req.Parameter.Value), ledger.WithTargetEndpoints(target...))
	if err != nil {
		logger.Errorf("failed to query block by txID : %s", err)
		return nil, err
	}

	return *resp, nil
}

func getBlockRange(jsonArgs string) (uint64, uint64) {
	regex := regexp.MustCompile(utils.RegexpJsonReplace)
	values := strings.Split(regex.ReplaceAllString(jsonArgs, ""), "-")
	if len(values) == 2 {
		min, _ := strconv.ParseUint(values[0], 10, 64)
		max, _ := strconv.ParseUint(values[1], 10, 64)
		return min, max
	}

	logger.Errorf("invalid args : %s", jsonArgs)
	return 0, 0
}
