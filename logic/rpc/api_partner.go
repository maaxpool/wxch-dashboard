package rpc

import (
	"github.com/ethereum/go-ethereum/common"
	"net/http"
	"strconv"
	"wxch-dashboard/config"
	"wxch-dashboard/logic/db"
)

func getPartnerListHandler(typ interface{}, r *http.Request) (resp interface{}, err error) {
	req := r.URL.Query()

	// stat total transaction amount
	total := db.GetPartnerCountByStatusRole("available", "broker")
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

	// get partner list
	partners, err := db.FindPaginationPartnersByStatusRole("available", "broker", offset, size)
	if err != nil {
		return nil, err
	}

	partnerItem := make([]partnerListItem, 0)
	for _, partner := range partners {
		partnerItem = append(partnerItem, partnerListItem{
			Name:                           partner.Name,
			Role:                           partner.Role,
			EthAddress:                     partner.EthAddress,
			ChiaCustodianDepositoryAddress: partner.ChiaCustodianDepositoryAddress,
			ChiaBrokerDepositAddress:       partner.ChiaBrokerDepositAddress,
			BridgeUrl:                      partner.BridgeUrl,
		})
	}
	return &getPartnerListResponse{
		Total:    total,
		Partners: partnerItem,
	}, nil
}

func getPartnerAssetsListHandler(typ interface{}, r *http.Request) (resp interface{}, err error) {
	req := r.URL.Query()

	// stat total transaction amount
	total := db.GetPartnerCountByStatusRole("available", "broker")
	if total == 0 {
		return &getPartnerListResponse{
			Total:    0,
			Partners: []partnerListItem{},
		}, nil
	}

	// stat total balance
	totalBalance := db.GetPartnerTotalBalance()

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

	// get partner list
	partners, err := db.FindPaginationPartnersByStatusRole("available", "broker", offset, size)
	if err != nil {
		return nil, err
	}

	partnerItem := make([]partnerAssetsListItem, 0)
	for _, partner := range partners {
		partnerItem = append(partnerItem, partnerAssetsListItem{
			ChiaCustodianDepositoryAddress: partner.ChiaCustodianDepositoryAddress,
			Balance:                        partner.Balance,
		})
	}
	return &getPartnerAssetsListResponse{
		Total:        total,
		TotalBalance: totalBalance,
		Partners:     partnerItem,
	}, nil
}

func createPartnerHandler(obj interface{}, r *http.Request) (resp interface{}, err error) {
	req := obj.(*createPartnerRequest)

	if req.Token != config.Get().RPC.AdminToken {
		return nil, NewHttpError(0xE001001, "invalid admin access token")
	}

	// check eth address
	ethAddress := common.HexToAddress(req.EthAddress).Hex()
	partner, err := db.FindPartnerByEthAddress(ethAddress)
	if err != nil {
		return nil, err
	}

	if partner.ID > 0 {
		return nil, NewHttpError(0xE001002, "partner already exist")
	}

	newPartner := &db.Partner{
		Name:                           req.Name,
		Role:                           "broker",
		EthAddress:                     req.EthAddress,
		ChiaCustodianDepositoryAddress: req.ChiaCustodianDepositoryAddress,
		ChiaBrokerDepositAddress:       req.ChiaBrokerDepositAddress,
		Balance:                        0,
		BridgeUrl:                      req.BridgeUrl,
		Status:                         "available",
	}

	err = db.SavePartner(newPartner)
	if err != nil {
		return nil, err
	}

	return true, nil
}
