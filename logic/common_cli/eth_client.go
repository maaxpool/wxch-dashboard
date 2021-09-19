package common_cli

import (
	"sync"
	"wxch-dashboard/config"
	"wxch-dashboard/logic/eth"
)

var (
	ethClientInit = sync.Once{}
	ethClient     *eth.Client
)

func GetEthClient() (cli *eth.Client) {
	ethClientInit.Do(func() {
		var err error

		// init eth client
		ethClient, err = eth.NewClient(config.Get().EthNode.InfuraHost)
		if err != nil {
			panic(err)
		}
	})

	return ethClient
}
