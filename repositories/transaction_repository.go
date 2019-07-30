package repositories

import (
	"clearance-adapter/factory"
	"clearance-adapter/models"
	"log"
)

// TransactionRepository Transaction仓库
type TransactionRepository struct{}

// Save 保存Transactions
func (TransactionRepository) Save(transactions []models.Transaction) {
	for _, transaction := range transactions {
		if (TransactionRepository{}.validateTransaction(transaction)) {
			TransactionRepository{}.save(transaction)
		}
	}
}

func (TransactionRepository) validateTransaction(transaction models.Transaction) bool {
	sql := `SELECT id
				FROM transactions
				WHERE waybill_no = ? AND box_no = ? AND sku_code = ?
			`
	result, err := factory.GetClrEngine().Query(sql, transaction.WaybillNo, transaction.BoxNo, transaction.SkuCode)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if len(result) == 0 {
		return true
	}
	return false
}

func (TransactionRepository) save(transaction models.Transaction) error {
	_, err := factory.GetClrEngine().Insert(&transaction)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
