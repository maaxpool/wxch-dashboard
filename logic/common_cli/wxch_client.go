package common_cli

import (
	"sync"
	"wxch-dashboard/config"
	"wxch-dashboard/logic/wxch"
)

var (
	wxchClientInit = sync.Once{}
	wxchClient     *wxch.Wxch
)

func GetWXCHClientClient() (cli *wxch.Wxch) {
	wxchClientInit.Do(func() {
		var err error

		wxchClient, err = GetEthClient().GetWxch(config.Get().EthNode.WxchBridgeContractAddress)
		if err != nil {
			panic(err)
		}
	})

	return wxchClient
}
