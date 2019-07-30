package main

import (
	"clearance-adapter/domain/entities"
	"clearance-adapter/repositories"
	"context"
	"errors"
	"log"
	"time"

	"github.com/pangpanglabs/goetl"
)

// InStorageETL 入库 p2-brand -> CSL
type InStorageETL struct{}

// New 创建InStorageETL对象，从Clearance到CSL
func (InStorageETL) New() *goetl.ETL {
	etl := goetl.New(InStorageETL{})
	etl.Before(InStorageETL{}.buildTransactions)
	etl.Before(InStorageETL{}.filterStorableTransactions)
	return etl
}

// Extract ...
func (etl InStorageETL) Extract(ctx context.Context) (interface{}, error) {
	local, _ := time.LoadLocation("Local")
	start, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-07-01 00:00:00", local)
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-07-31 23:59:59", local)
	result, err := repositories.StockTransactionRepository{}.GetInStorageByCreateAt(start, end)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 相同运单号合并为一个Transaction
func (etl InStorageETL) buildTransactions(ctx context.Context, source interface{}) (interface{}, error) {
	data, ok := source.([]map[string]string)
	if !ok {
		return nil, errors.New("Convert Failed")
	}

	items := make(map[string][]map[string]string, 0)
	for _, item := range data {
		key := item["brand_coce"] + "-" + item["shop_coce"] + "-" + item["waybill_no"] + "-" + item["box_no"]
		if _, ok := items[key]; ok {
			items[key] = append(items[key], item)
		} else {
			items[key] = make([]map[string]string, 0)
			items[key] = append(items[key], item)
		}
	}

	transactions := make([]entities.Transaction, 0)
	for _, v := range items {
		txn, err := entities.Transaction{}.Create(v)
		if err != nil {
			continue
		}
		transactions = append(transactions, txn)
	}

	return transactions, nil
}

// validateTransactions 验证transaction是否可以入库
func (etl InStorageETL) validateTransaction(transaction entities.Transaction) (bool, error) {
	recvSupp, err := repositories.RecvSuppRepository{}.GetByWaybillNo(transaction.BrandCode, transaction.ShopCode, transaction.WaybillNo)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	if len(recvSupp) == 0 {
		return false, errors.New("there is no outbound order which waybill no is: " + transaction.WaybillNo)
	}

	ok := true
	for _, v := range recvSupp {
		ok = ok && (v.RecvSuppStatusCode == "R")
	}

	if !ok {
		return false, errors.New("some sku already in storage")
	}

	return true, nil
}

// filterStorableTransactions 过滤出可以入库的运单
func (etl InStorageETL) filterStorableTransactions(ctx context.Context, source interface{}) (interface{}, error) {
	transactions, ok := source.([]entities.Transaction)
	if !ok {
		return nil, errors.New("Convert Failed")
	}

	storableTransactions := make([]entities.Transaction, 0)
	for _, txn := range transactions {
		ok, err := InStorageETL{}.validateTransaction(txn)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if ok {
			storableTransactions = append(storableTransactions, txn)
		}
	}
	return storableTransactions, nil
}

// Transform ...
func (etl InStorageETL) Transform(ctx context.Context, source interface{}) (interface{}, error) {
	transactions, ok := source.([]entities.Transaction)
	if !ok {
		return nil, errors.New("Convert Failed")
	}

	return transactions, nil
}

// Load ...
func (etl InStorageETL) Load(ctx context.Context, source interface{}) error {
	transactions, ok := source.([]entities.Transaction)
	if !ok {
		return errors.New("Convert Failed")
	}

	for _, txn := range transactions {
		err := repositories.RecvSuppRepository{}.PutInStorage(txn.BrandCode, txn.ShopCode, txn.WaybillNo, txn.UserID)
		if err != nil {
			log.Printf(err.Error())
		}
	}

	return nil
}
