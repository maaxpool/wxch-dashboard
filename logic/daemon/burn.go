package daemon

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"wxch-dashboard/logic/db"
)

func burnConfirmedHandler(transLog types.Log) (err error) {
	// check transaction hash exist
	if isExist := db.CheckTransactionIsExistByTypeReviewHash("burn", transLog.TxHash.String()); isExist {
		return nil
	}

	// parse event data
	var burnConfirmedEvent LogEventBurnConfirmed
	err = wxchBridgeContractAbi.UnpackIntoInterface(&burnConfirmedEvent, "BurnConfirmed", transLog.Data)
	if err != nil {
		return
	}

	burnConfirmedEvent.Requester = common.HexToAddress(transLog.Topics[2].Hex())
	burnAmount := burnConfirmedEvent.Amount.Div(burnConfirmedEvent.Amount, big.NewInt(XchDecimalBase))
	requestHash := fmt.Sprintf("%x", burnConfirmedEvent.InputRequestHash)

	fmt.Printf("Log Block Number: %d\n", transLog.BlockNumber)
	fmt.Printf("Requester: %s\n", burnConfirmedEvent.Requester.Hex())
	fmt.Printf("Burn Amount: %d\n", burnConfirmedEvent.Amount.Int64())

	// create new transaction log
	newTransactionLog := &db.Transaction{
		Type:             "burn",
		Amount:           float64(burnAmount.Uint64()),
		FeeAmount:        0,
		SenderAddress:    burnConfirmedEvent.Requester.Hex(),
		ReceiverAddress:  burnConfirmedEvent.DepositAddress,
		Status:           "burn_completed",
		EthRequestTxHash: requestHash,
		EthReviewTxHash:  transLog.TxHash.String(),
		ChiaSendTxHash:   burnConfirmedEvent.Txid,
	}

	// get partner
	partner, _ := db.FindPartnerByEthAddress(burnConfirmedEvent.Requester.Hex())
	if partner.ID > 0 {
		newTransactionLog.PartnerId = partner.ID
		newTransactionLog.PartnerName = partner.Name
	} else {
		newTransactionLog.PartnerId = 0
		newTransactionLog.PartnerName = "Unknown Partner"
	}

	err = db.SaveTransaction(newTransactionLog)
	if err != nil {
		return
	}

	return
}
