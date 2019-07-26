package main

import (
	"context"
	"errors"
	"strconv"

	"clearance-adapter/factory"
	"clearance-adapter/models"
)

//AssignMslETL ...
type AssignMslETL struct{}

// Extract ...
func (etl AssignMslETL) Extract(ctx context.Context) (interface{}, error) {
	engine := factory.GetMslEngine()
	sourceData := make([]models.StockTrans, 0)
	var tranDtl models.StockTrans
	//engine.Join("INNER", "stock_transaction", "stock_transaction.id = stock_transaction_item.stock_transaction_id").Where(" 1=1 ").Find(&sourceData)
	sql := `select m.id,m.waybill_no,d.sku_id,d.quantity
		from stock_transaction m
		join stock_transaction_item d on m.id = d.stock_transaction_id 
		where type = 'IN' `
	res, err := engine.Query(sql)
	if err != nil {
		return nil, err
	}
	for _, data := range res {
		tranDtl.StockTransaction.WaybillNo = string(data["waybill_no"])
		tranDtl.StockTransactionItem.SkuID, _ = strconv.ParseInt(string(data["sku_id"]), 10, 32)
		tranDtl.StockTransactionItem.Quantity, _ = strconv.Atoi(string(data["quantity"]))
		sourceData = append(sourceData, tranDtl)
	}
	return sourceData, nil
}

//Transform ...
func (etl AssignMslETL) Transform(ctx context.Context, target interface{}) (interface{}, error) {
	stockTrans, ok := target.([]models.StockTrans)
	if !ok {
		return nil, errors.New("Convert Failed")
	}
	transactions := make([]models.Transaction, 0)
	for _, stockTran := range stockTrans {
		transactions = append(transactions, models.Transaction{
			TransactionID: string(stockTran.StockTransaction.ID),
			WaybillNo:     stockTran.StockTransaction.WaybillNo,
			BoxNo:         stockTran.StockTransaction.BoxNo,
			SkuCode:       string(stockTran.StockTransactionItem.SkuID),
			Qty:           stockTran.StockTransactionItem.Quantity,
		})
	}

	return transactions, nil
}

// Load ...
func (etl AssignMslETL) Load(ctx context.Context, target interface{}) error {
	if target == nil {
		return errors.New("target is nil")
	}
	transactions, ok := target.([]models.Transaction)
	if !ok {
		return errors.New("Convert Failed")
	}
	engine := factory.GetClrEngine()

	if _, err := engine.Insert(&transactions); err != nil {
		return err
	}
	return nil
}
