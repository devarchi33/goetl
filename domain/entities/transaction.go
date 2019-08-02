package entities

import (
	"clearance-adapter/infra"
	"errors"
)

// Transaction 以StockTransaction为主的信息，会把id替换成code
type Transaction struct {
	BrandCode string
	ShopCode  string
	WaybillNo string
	BoxNo     string
	EmpID     string
	Items     []TransactionItem
}

// TransactionItem 以StockTransactionItem为主的信息，会把id替换成code
type TransactionItem struct {
	SkuCode string
	Qty     int
}

// Create 根据[]map类型的数据转换成Transaction
func (Transaction) Create(data []map[string]string) (Transaction, error) {
	transaction := Transaction{}
	if data == nil || len(data) == 0 {
		return transaction, errors.New("data is empty")
	}

	txnData := data[0]
	err := Transaction{}.checkRequirement(txnData, "brand_code", "shop_code", "waybill_no", "box_no", "emp_id")
	if err != nil {
		return transaction, err
	}

	transaction.BrandCode = txnData["brand_code"]
	transaction.ShopCode = txnData["shop_code"]
	transaction.WaybillNo = txnData["waybill_no"]
	transaction.BoxNo = txnData["box_no"]
	transaction.EmpID = txnData["emp_id"]
	transaction.Items = make([]TransactionItem, 0)

	for _, item := range data {
		err := Transaction{}.checkRequirement(txnData, "sku_code", "qty")
		if err != nil {
			return transaction, err
		}
		transaction.Items = append(transaction.Items, TransactionItem{
			SkuCode: item["sku_code"],
			Qty:     infra.ConvertStringToInt(item["qty"]),
		})
	}

	return transaction, nil
}

func (Transaction) checkRequirement(data map[string]string, requiredKeys ...string) error {
	for _, key := range requiredKeys {
		if _, ok := data[key]; !ok {
			return errors.New(key + " is required")
		}
	}
	return nil
}
