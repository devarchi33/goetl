package main

import (
	"clearance-adapter/factory"
	"clearance-adapter/models"
	"clearance-adapter/repositories"
	"context"
	"errors"

	"github.com/pangpanglabs/goetl"
)

// InStorageETL 入库
type InStorageETL struct{}

// New 创建InStorageETL对象，从Clearance到CSL
func (InStorageETL) New() *goetl.ETL {
	etl := goetl.New(InStorageETL{})
	return etl
}

// Extract ...
func (etl InStorageETL) Extract(ctx context.Context) (interface{}, error) {
	engine := factory.GetClrEngine()
	transactions := make([]models.Transaction, 0)
	engine.Find(&transactions)

	return transactions, nil
}

// Transform ...
func (etl InStorageETL) Transform(ctx context.Context, source interface{}) (interface{}, error) {
	transactions, ok := source.([]models.Transaction)
	if !ok {
		return nil, errors.New("Convert Failed")
	}

	mastersMap := make(map[string]models.RecvSuppMst, 0)
	details := make([]models.RecvSuppDtl, 0)
	for _, txn := range transactions {
		details = append(details, models.RecvSuppDtl{
			BrandCode:        "SA",
			ShopCode:         "CFW5",
			RecvSuppNo:       txn.TransactionID,
			ProdCode:         txn.SkuCode,
			RecvSuppQty:      txn.Qty,
			RecvSuppFixedQty: txn.Qty,
		})
		key := txn.TransactionID + "-" + txn.WaybillNo + "-" + txn.BoxNo
		mastersMap[key] = models.RecvSuppMst{
			RecvSuppNo: txn.TransactionID,
			BrandCode:  "SA",
			ShopCode:   "CFW5",
			WayBillNo:  txn.WaybillNo,
			BoxNo:      txn.BoxNo,
		}
	}

	masters := make([]models.RecvSuppMst, 0)
	for _, mst := range mastersMap {
		masters = append(masters, mst)
	}

	return map[string]interface{}{
		"RecvSuppMst": masters,
		"RecvSuppDtl": details,
	}, nil
}

// Load ...
func (etl InStorageETL) Load(ctx context.Context, source interface{}) error {
	mstDtlMap, ok := source.(map[string]interface{})
	if !ok {
		return errors.New("Convert Failed")
	}
	masters, ok := mstDtlMap["RecvSuppMst"].([]models.RecvSuppMst)
	if !ok {
		return errors.New("Convert Failed")
	}
	details, ok := mstDtlMap["RecvSuppDtl"].([]models.RecvSuppDtl)
	if !ok {
		return errors.New("Convert Failed")
	}

	repositories.RecvSuppRepository{}.SaveMasters(masters)
	repositories.RecvSuppRepository{}.SaveDetails(details)

	return nil
}
