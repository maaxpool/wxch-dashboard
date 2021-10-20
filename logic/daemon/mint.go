package daemon

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"wxch-dashboard/logic/db"
)

func mintConfirmedHandler(transLog types.Log) (err error) {
	// check transaction hash exist
	if isExist := db.CheckTransactionIsExistByTypeReviewHash("mint", transLog.TxHash.String()); isExist {
		return nil
	}

	// parse event data
	var mintConfirmedEvent LogEventMintConfirmed
	err = wxchBridgeContractAbi.UnpackIntoInterface(&mintConfirmedEvent, "MintConfirmed", transLog.Data)
	if err != nil {
		return
	}

	mintConfirmedEvent.Requester = common.HexToAddress(transLog.Topics[2].Hex())
	mintAmount := mintConfirmedEvent.Amount.Div(mintConfirmedEvent.Amount, big.NewInt(XchDecimalBase))
	requestHash := fmt.Sprintf("%x", mintConfirmedEvent.RequestHash)

	fmt.Printf("Log Block Number: %d\n", transLog.BlockNumber)
	fmt.Printf("Requester: %s\n", mintConfirmedEvent.Requester.Hex())
	fmt.Printf("Mint Amount: %d\n", mintConfirmedEvent.Amount.Int64())

	// create new transaction log
	amount := float64(mintAmount.Uint64())
	newTransactionLog := &db.Transaction{
		Type:             "mint",
		Amount:           amount,
		FeeAmount:        0,
		SenderAddress:    mintConfirmedEvent.DepositAddress,
		ReceiverAddress:  mintConfirmedEvent.Requester.Hex(),
		Status:           "mint_completed",
		EthRequestTxHash: requestHash,
		EthReviewTxHash:  transLog.TxHash.String(),
		ChiaSendTxHash:   mintConfirmedEvent.Txid,
	}

	// get partner
	partner, _ := db.FindPartnerByEthAddress(mintConfirmedEvent.Requester.Hex())
	if partner.ID > 0 {
		newTransactionLog.PartnerId = partner.ID
		newTransactionLog.PartnerName = partner.Name

		// update partner balance
		_ = db.UpdateBalanceById(partner.ID, partner.Balance+amount)
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

func mintRejectedHandler(transLog types.Log) (err error) {
	// check transaction hash exist
	if isExist := db.CheckTransactionIsExistByTypeReviewHash("mint", transLog.TxHash.String()); isExist {
		return nil
	}

	// parse event data
	var mintRejectedEvent LogEventMintRejected
	err = wxchBridgeContractAbi.UnpackIntoInterface(&mintRejectedEvent, "MintRejected", transLog.Data)
	if err != nil {
		return
	}

	mintRejectedEvent.Requester = common.HexToAddress(transLog.Topics[2].Hex())
	mintAmount := mintRejectedEvent.Amount.Div(mintRejectedEvent.Amount, big.NewInt(XchDecimalBase))
	requestHash := fmt.Sprintf("%x", mintRejectedEvent.RequestHash)

	fmt.Printf("Log Block Number: %d\n", transLog.BlockNumber)
	fmt.Printf("Requester: %s\n", mintRejectedEvent.Requester.Hex())
	fmt.Printf("Mint Amount: %d\n", mintRejectedEvent.Amount.Int64())

	// create new transaction log
	newTransactionLog := &db.Transaction{
		Type:             "mint",
		Amount:           float64(mintAmount.Uint64()),
		FeeAmount:        0,
		SenderAddress:    mintRejectedEvent.DepositAddress,
		ReceiverAddress:  mintRejectedEvent.Requester.Hex(),
		Status:           "mint_rejected",
		EthRequestTxHash: requestHash,
		EthReviewTxHash:  transLog.TxHash.String(),
		ChiaSendTxHash:   mintRejectedEvent.Txid,
	}

	// get partner
	partner, _ := db.FindPartnerByEthAddress(mintRejectedEvent.Requester.Hex())
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
