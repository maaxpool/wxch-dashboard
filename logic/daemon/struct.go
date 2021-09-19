package daemon

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type LogEventMintConfirmed struct {
	Nonce          *big.Int
	Requester      common.Address
	Amount         *big.Int
	DepositAddress string
	Txid           string
	Timestamp      *big.Int
	RequestHash    [32]byte
}

type LogEventMintRejected struct {
	Nonce          *big.Int
	Requester      common.Address
	Amount         *big.Int
	DepositAddress string
	Txid           string
	Timestamp      *big.Int
	RequestHash    [32]byte
}

type LogEventBurnConfirmed struct {
	Nonce            *big.Int
	Requester        common.Address
	Amount           *big.Int
	DepositAddress   string
	Txid             string
	Timestamp        *big.Int
	InputRequestHash [32]byte
}
