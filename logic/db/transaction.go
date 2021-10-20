package db

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	PartnerId        uint    `gorm:"not null;default:0" json:"partner_id"`
	PartnerName      string  `gorm:"size:100;not null;default:Unknown Partner" json:"partner_name"`
	Type             string  `gorm:"size:30;not null" json:"type"` // 'mint', 'burn'
	Amount           float64 `gorm:"not null" json:"amount"`
	FeeAmount        float64 `gorm:"not null" json:"fee_amount"`
	SenderAddress    string  `gorm:"size:500;not null" json:"sender_address"`
	ReceiverAddress  string  `gorm:"size:500;not null" json:"receiver_address"`
	Status           string  `gorm:"size:30;not null" json:"status"` // 'mint_completed','mint_rejected','burn_completed'
	EthRequestTxHash string  `gorm:"size:500;not null" json:"eth_request_tx_hash"`
	EthReviewTxHash  string  `gorm:"size:500;not null" json:"eth_review_tx_hash"`
	ChiaSendTxHash   string  `gorm:"size:500;not null" json:"chia_send_tx_hash"`
}

func GetTransactionById(transactionId uint) (transaction Transaction, err error) {
	err = db.Where("id = ?", transactionId).Find(&transaction).Error
	return
}

func CheckTransactionIsExistByTypeReviewHash(transactionType string, hash string) bool {
	var count int64
	db.Table("transactions").Where("type = ? AND eth_review_tx_hash = ?", transactionType, hash).Count(&count)
	return count > 0
}

func GetTransactionByUserIdTypeEthHash(transactionType string, userId uint, ethHash string) (transaction Transaction, err error) {
	err = db.Where("user_Id = ? AND type = ? AND eth_transaction_hash = ?", userId, transactionType, ethHash).Find(&transaction).Error
	return
}

func GetTransactionCountByType(transactionType string) uint {
	var count int64
	query := db.Table("transactions")

	if transactionType != "all" {
		query = query.Where("type = ?", transactionType)
	}

	query.Where("deleted_at IS NULL").Count(&count)
	return uint(count)
}

func FindPaginationTransactionsByType(transactionType string, offset int, size int) (transactions []Transaction, err error) {
	query := db.Table("transactions")

	if transactionType != "all" {
		query = query.Where("type = ?", transactionType)
	}

	err = query.Limit(size).Offset(offset).Order("id desc").Find(&transactions).Error
	return
}

func FindTransactionsByTypeStatus(transactionType string, status string) (transactions []Transaction, err error) {
	err = db.Where("type = ? AND status = ?", transactionType, status).Find(&transactions).Error
	return
}

func UpdateTransactionStatusById(transactionId uint, status string) error {
	ret := db.Table("transactions").Where("id = ?", transactionId).UpdateColumn("status", status)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

func UpdateTransactionChiaTransactionHashById(transactionId uint, chiaTransactionHash string) error {
	ret := db.Table("transactions").Where("id = ?", transactionId).UpdateColumn("chia_transaction_hash", chiaTransactionHash)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

func UpdateTransactionEthTransactionHashById(transactionId uint, ethTransactionHash string) error {
	ret := db.Table("transactions").Where("id = ?", transactionId).UpdateColumn("eth_transaction_hash", ethTransactionHash)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

func GetTotalTransactionAmountByTypeUserId(transactionType string, userId uint) (float64, error) {
	type TransactionAmountStatResult struct {
		Total float64
	}
	result := new(TransactionAmountStatResult)
	err := db.Table("transactions").Select("sum(amount) as total").Where("type = ? AND user_id = ? AND status = 'finished'", transactionType, userId).Find(result).Error
	return result.Total, err
}

func GetTotalMintAmount() float64 {
	var total float64
	db.Table("transactions").Where("status = ? AND deleted_at IS NULL", "mint_completed").Pluck("COALESCE(SUM(amount), 0) as total", &total)
	return total
}

func GetTotalBurnAmount() float64 {
	var total float64
	db.Table("transactions").Where("status = ? AND deleted_at IS NULL", "burn_completed").Pluck("COALESCE(SUM(amount), 0) as total", &total)
	return total
}

func SaveTransaction(transaction *Transaction) (err error) {
	return db.Save(transaction).Error
}
