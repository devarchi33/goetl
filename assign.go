package main

import (
	"context"
	"errors"

	"clearance-adapter/factory"
	"clearance-adapter/models"

	"github.com/pangpanglabs/goetl"
)

// ConvertByteResult ...
func ConvertByteResult(source []map[string][]byte) []map[string]string {
	result := make([]map[string]string, 0)
	if source == nil || len(source) == 0 {
		return result
	}
	for _, sourceItem := range source {
		item := make(map[string]string, 0)
		if sourceItem == nil || len(sourceItem) == 0 {
			continue
		}
		for key, value := range sourceItem {
			item[key] = string(value)
		}
		result = append(result, item)
	}
	return result
}

// AssignETL 分配
type AssignETL struct{}

// New 创建AssignETL对象
func (AssignETL) New() *goetl.ETL {
	etl := goetl.New(AssignETL{})
	etl.After(AssignETL{}.ReadyToLoad)
	return etl
}

// Extract ...
func (etl AssignETL) Extract(ctx context.Context) (interface{}, error) {
	engine := factory.GetCSLEngine()
	details := make([]models.RecvSupp, 0)
	engine.Join("INNER", "RecvSuppMst", "RecvSuppMst.RecvSuppNo = RecvSuppDtl.RecvSuppNo").
		Where("RecvSuppMst.BrandCode = ? AND RecvSuppMst.ShopCode = ?", "SA", "CFW5").
		Find(&details)

	return details, nil
}

// Transform ...
func (etl AssignETL) Transform(ctx context.Context, source interface{}) (interface{}, error) {
	recvSuppList, ok := source.([]models.RecvSupp)
	if !ok {
		return nil, errors.New("Convert Failed")
	}
	transactions := make([]models.Transaction, 0)
	for _, recvSupp := range recvSuppList {
		transactions = append(transactions, models.Transaction{
			TransactionID: recvSupp.RecvSuppMst.RecvSuppNo,
			WaybillNo:     recvSupp.RecvSuppMst.WayBillNo,
			BoxNo:         recvSupp.RecvSuppMst.BoxNo,
			SkuCode:       recvSupp.RecvSuppDtl.ProdCode,
			Qty:           recvSupp.RecvSuppDtl.RecvSuppQty,
		})
	}

	return transactions, nil
}

// ReadyToLoad ...
func (etl AssignETL) ReadyToLoad(ctx context.Context, source interface{}) error {
	masters, ok := source.([]models.Transaction)
	if !ok {
		return errors.New("Convert Failed")
	}
	savedMasters := make([]models.Transaction, 0)
	for _, recvSupp := range masters {
		sql := `SELECT id
		FROM transactions
		WHERE transaction_id = ?
		`

		engine := factory.GetClrEngine()
		result, err := engine.Query(sql, recvSupp.TransactionID)
		if err != nil {
			return err
		}

		if len(result) == 0 {
			savedMasters = append(savedMasters, recvSupp)
		}
	}

	return nil
}

// Load ...
func (etl AssignETL) Load(ctx context.Context, source interface{}) error {
	if source == nil {
		return errors.New("source is nil")
	}
	transactions, ok := source.([]models.Transaction)
	if !ok {
		return errors.New("Convert Failed")
	}
	engine := factory.GetClrEngine()

	if _, err := engine.Insert(&transactions); err != nil {
		return err
	}
	return nil
}
