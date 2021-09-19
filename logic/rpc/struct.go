package rpc

import "time"

type getTransactionListRequest struct {
	Type string `schema:"type" validate:"required"`
	Page uint   `schema:"page" validate:"required"`
	Size uint   `schema:"size" validate:"required"`
}

type createPartnerRequest struct {
	Name                           string `json:"name" validate:"required,lte=100"`
	EthAddress                     string `json:"eth_address" validate:"required"`
	ChiaCustodianDepositoryAddress string `json:"chia_custodian_depository_address" validate:"required"`
	ChiaBrokerDepositAddress       string `json:"chia_broker_deposit_address" validate:"required"`
	BridgeUrl                      string `json:"bridge_url" validate:"required"`
	Token                          string `json:"token" validate:"required"`
}

type getPartnerListRequest struct {
	Page uint `schema:"page" validate:"required"`
	Size uint `schema:"size" validate:"required"`
}

// ============ Response ===============

type getTransactionListResponse struct {
	Total        uint                  `json:"total"`
	Transactions []transactionListItem `json:"transactions"`
}

type transactionListItem struct {
	Id               uint      `json:"id"`
	Type             string    `json:"type"`
	PartnerName      string    `json:"partner_name"`
	Amount           float64   `json:"amount"`
	FeeAmount        float64   `json:"fee_amount"`
	EthRequestTxHash string    `json:"eth_request_tx_hash"`
	EthReviewTxHash  string    `json:"eth_review_tx_hash"`
	ChiaSendTxHash   string    `json:"chia_send_tx_hash"`
	CreatedAt        time.Time `json:"created_at"`
}

type getPartnerListResponse struct {
	Total    uint              `json:"total"`
	Partners []partnerListItem `json:"partners"`
}

type partnerListItem struct {
	Name                           string `json:"name"`
	Role                           string `json:"role"`
	EthAddress                     string `json:"eth_address"`
	ChiaCustodianDepositoryAddress string `json:"chia_custodian_depository_address"`
	ChiaBrokerDepositAddress       string `json:"chia_broker_deposit_address"`
	BridgeUrl                      string `json:"bridge_url"`
}
