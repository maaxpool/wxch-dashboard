package rpc

import (
	"net/http"
	"wxch-dashboard/logic/db"
)

func getGlobalStatHandler(r *http.Request) (resp interface{}, err error) {
	// get total mint amount
	totalMint := db.GetTotalMintAmount()

	// get total burn amount
	totalBurn := db.GetTotalBurnAmount()

	// stat total partner balance
	totalPartnerBalance := db.GetPartnerTotalBalance()

	return &getGlobalStatResponse{
		NetworkAmount: totalMint - totalBurn,
		CustodyAmount: totalPartnerBalance,
	}, nil
}
