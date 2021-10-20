package db

import (
	"gorm.io/gorm"
)

type Partner struct {
	gorm.Model
	Name                           string  `gorm:"size:100;not null" json:"name"`
	Role                           string  `gorm:"size:30;not null" json:"role"` // broker
	EthAddress                     string  `gorm:"size:500;not null" json:"eth_address"`
	ChiaCustodianDepositoryAddress string  `gorm:"size:500;not null" json:"chia_custodian_depository_address"`
	ChiaBrokerDepositAddress       string  `gorm:"size:500;not null" json:"chia_broker_deposit_address"`
	Balance                        float64 `gorm:"not null" json:"balance"`
	BridgeUrl                      string  `gorm:"size:500;not null" json:"bridge_url"`
	Status                         string  `gorm:"size:30;not null" json:"status"` // available, available
}

func FindPartnerByEthAddress(ethAddress string) (partner Partner, err error) {
	err = db.Where("eth_address = ?", ethAddress).Find(&partner).Error
	return
}

func GetPartnerCountByStatusRole(status string, role string) uint {
	var count int64
	db.Table("partners").Where("status = ? AND role = ? AND deleted_at IS NULL", status, role).Count(&count)
	return uint(count)
}

func GetPartnerTotalBalance() float64 {
	var total float64
	db.Table("partners").Where("status = ? AND role = ? AND deleted_at IS NULL", "available", "broker").Pluck("COALESCE(SUM(balance), 0) as total", &total)
	return total
}

func FindPaginationPartnersByStatusRole(status string, role string, offset int, size int) (partners []Partner, err error) {
	err = db.Where("status = ? AND role = ?", status, role).Limit(size).Offset(offset).Order("id desc").Find(&partners).Error
	return
}

func UpdateBalanceById(partnerId uint, balance float64) error {
	err := db.Where("id = ?", partnerId).Update("balance", balance).Error
	return err
}

func SavePartner(partner *Partner) (err error) {
	return db.Save(partner).Error
}
