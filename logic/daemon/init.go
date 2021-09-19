package daemon

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/getsentry/sentry-go"
	"github.com/robfig/cron/v3"
	"math/big"
	"strconv"
	"strings"
	"wxch-dashboard/config"
	"wxch-dashboard/logic/common_cli"
	"wxch-dashboard/logic/db"
	"wxch-dashboard/logic/log"
	"wxch-dashboard/logic/wxch"
)

var (
	isEventCheckDoing     = false
	cronClient            *cron.Cron
	wxchBridgeContractAbi abi.ABI
	mintConfirmedSigHash  = crypto.Keccak256Hash([]byte("MintConfirmed(uint256,address,uint256,string,string,uint256,bytes32)")).Hex()
	mintRejectedSigHash   = crypto.Keccak256Hash([]byte("MintRejected(uint256,address,uint256,string,string,uint256,bytes32)")).Hex()
	burnConfirmedSigHash  = crypto.Keccak256Hash([]byte("BurnConfirmed(uint256,address,uint256,string,string,uint256,bytes32)")).Hex()
)

const (
	SecurityIntervalBlock = 5
	XchDecimalBase        = 1e12
)

func InitCron() {
	cronClient = cron.New(cron.WithSeconds())

	_, _ = cronClient.AddFunc("*/10 * * * * *", wxchBridgeEventHandler)

	cronClient.Start()
}

func StopCron(ctx context.Context) {
	if cronClient == nil {
		return
	}

	select {
	case <-cronClient.Stop().Done():
		return
	case <-ctx.Done():
		return
	}
}

func wxchBridgeEventHandler() {
	if isEventCheckDoing {
		return
	}
	isEventCheckDoing = true
	defer func() {
		isEventCheckDoing = false
	}()

	// get latest block
	latestHeader, err := common_cli.GetEthClient().HeaderByNumber(context.Background(), nil)
	if err != nil {
		sentry.CaptureException(err)
		log.GetLogger().Error(err.Error())
		return
	}
	needCheckBlockNum := latestHeader.Number.Int64() - SecurityIntervalBlock

	// get last check block num
	configVal, err := db.FindConfigByKeyName("wxch_bridge_last_check")
	if err != nil {
		sentry.CaptureException(err)
		log.GetLogger().Error(err.Error())
		return
	}
	lastCheckBlockNum, _ := strconv.Atoi(configVal.KeyValue)

	// no checking if new block less than 10
	if needCheckBlockNum-int64(lastCheckBlockNum) < 10 {
		return
	}

	// get latest transfer event
	contractAddress := common.HexToAddress(config.Get().EthNode.WxchBridgeContractAddress)
	wxchBridgeContractAbi, _ = abi.JSON(strings.NewReader(wxch.WxchMetaData.ABI))

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(lastCheckBlockNum)),
		ToBlock:   big.NewInt(needCheckBlockNum - 1),
		// FromBlock: big.NewInt(int64(11057670)),
		// ToBlock:   big.NewInt(int64(11057679)),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	transLogs, err := common_cli.GetEthClient().FilterLogs(context.Background(), query)
	if err != nil {
		sentry.CaptureException(err)
		log.GetLogger().Error(err.Error())
		return
	}

	if len(transLogs) > 0 {
		var handleErr error
		for _, transLog := range transLogs {
			switch transLog.Topics[0].Hex() {
			case mintConfirmedSigHash:
				handleErr = mintConfirmedHandler(transLog)
			case mintRejectedSigHash:
				handleErr = mintRejectedHandler(transLog)
			case burnConfirmedSigHash:
				handleErr = burnConfirmedHandler(transLog)
			}

			if handleErr != nil {
				sentry.CaptureException(err)
				log.GetLogger().Error(err.Error())
			}
		}
	}

	err = db.UpdateValueByKeyName("wxch_bridge_last_check", strconv.Itoa(int(needCheckBlockNum)))
	if err != nil {
		sentry.CaptureException(err)
		log.GetLogger().Error(err.Error())
	}
}
