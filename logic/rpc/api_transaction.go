package rpc

import (
	"net/http"
	"strconv"
	"wxch-dashboard/logic/db"
)

func getTransactionListHandler(typ interface{}, r *http.Request) (resp interface{}, err error) {
	req := r.URL.Query()

	// stat total transaction amount
	total := db.GetTransactionCountByType(req.Get("type"))
	if total == 0 {
		return &getPartnerListResponse{
			Total:    0,
			Partners: []partnerListItem{},
		}, nil
	}

	// calc pagination
	page, err := strconv.Atoi(req.Get("page"))
	if err != nil {
		return nil, err
	}
	size, err := strconv.Atoi(req.Get("size"))
	if err != nil {
		return nil, err
	}
	offset := (page - 1) * size

	// get transaction list
	transactions, err := db.FindPaginationTransactionsByType(req.Get("type"), offset, size)
	if err != nil {
		return nil, err
	}

	transactionItem := make([]transactionListItem, 0)
	for _, transaction := range transactions {
		transactionItem = append(transactionItem, transactionListItem{
			Id:               transaction.ID,
			Type:             transaction.Type,
			PartnerName:      transaction.PartnerName,
			Amount:           transaction.Amount,
			FeeAmount:        transaction.FeeAmount,
			EthRequestTxHash: transaction.EthRequestTxHash,
			EthReviewTxHash:  transaction.EthReviewTxHash,
			ChiaSendTxHash:   transaction.ChiaSendTxHash,
			CreatedAt:        transaction.CreatedAt,
		})
	}
	return &getTransactionListResponse{
		Total:        total,
		Transactions: transactionItem,
	}, nil
}
