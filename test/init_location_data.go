package test

import (
	"clearance-adapter/factory"
	"clearance-adapter/models"
	"fmt"
	"log"
	"path/filepath"
	"time"
)

func initLocation() {
	createLocationDB()
	createStockTransactionTable()
	setStockTransactionData()
	createStockTransactionItemTable()
	setStockTransactionItemData()
}

func createLocationDB() {
	if _, err := factory.GetP2BrandEngine().Exec("DROP DATABASE IF EXISTS pangpang_brand_sku_location;"); err != nil {
		log.Printf("createLocationDB error: %v", err.Error())
		log.Println()
	}
	if _, err := factory.GetP2BrandEngine().Exec("CREATE DATABASE pangpang_brand_sku_location;"); err != nil {
		log.Printf("createLocationDB error: %v", err.Error())
		log.Println()
	}
}

func createStockTransactionTable() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_sku_location;"); err != nil {
		log.Printf("createLocationDB error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE stock_transaction
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			tenant_code VARCHAR(255),
			type VARCHAR(255),
			waybill_no VARCHAR(255),
			created_at DATETIME,
			colleague_id BIGINT(20),
			box_no VARCHAR(255),
			shipment_location_id BIGINT(20),
			receipt_location_id BIGINT(20),
			brand_code VARCHAR(255)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockTransactionTable error: %v", err.Error())
		log.Println()
	}
}

func setStockTransactionData() {
	filename, err := filepath.Abs("test/data/test_in_storage_etl_stock_transaction_data.csv")
	if err != nil {
		panic(err.Error())
	}
	headers, data := readDataFromCSV(filename)
	transactions := buildStockTransaction(headers, data)

	loadStockTransactionData(transactions)
}

func buildStockTransaction(headers map[int]string, data [][]string) []models.StockTransaction {
	transactions := make([]models.StockTransaction, 0)
	for _, row := range data {
		txn := new(models.StockTransaction)
		setObjectValue(headers, row, txn)
		local, _ := time.LoadLocation("Local")
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", row[4], local)
		txn.CreatedAt = t
		transactions = append(transactions, *txn)
	}

	return transactions
}

func loadStockTransactionData(transactions []models.StockTransaction) {
	for _, txn := range transactions {
		if affected, err := factory.GetP2BrandEngine().Insert(&txn); err != nil {
			fmt.Printf("loadStockTransactionData error: %v", err.Error())
			fmt.Println()
			fmt.Printf("affected: %v", affected)
			fmt.Println()
		}
	}
}

func createStockTransactionItemTable() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_sku_location;"); err != nil {
		log.Printf("createLocationDB error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE stock_transaction_item
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			stock_transaction_id BIGINT(20),
			product_id BIGINT(20),
			sku_id BIGINT(20),
			barcode VARCHAR(255),
			quantity BIGINT(20),
			created_at DATETIME,
			brand_code VARCHAR(255)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockTransactionItemTable error: %v", err.Error())
		log.Println()
	}
}

func setStockTransactionItemData() {
	filename, err := filepath.Abs("test/data/test_in_storage_etl_stock_transaction_item_data.csv")
	if err != nil {
		panic(err.Error())
	}
	headers, data := readDataFromCSV(filename)
	transactionItems := buildStockTransactionItem(headers, data)

	loadStockTransactionItemData(transactionItems)
}

func buildStockTransactionItem(headers map[int]string, data [][]string) []models.StockTransactionItem {
	transactionItems := make([]models.StockTransactionItem, 0)
	for _, row := range data {
		txn := new(models.StockTransactionItem)
		setObjectValue(headers, row, txn)
		local, _ := time.LoadLocation("Local")
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", row[6], local)
		txn.CreatedAt = t
		transactionItems = append(transactionItems, *txn)
	}

	return transactionItems
}

func loadStockTransactionItemData(transactionItems []models.StockTransactionItem) {
	for _, txn := range transactionItems {
		if affected, err := factory.GetP2BrandEngine().Insert(&txn); err != nil {
			fmt.Printf("loadStockTransactionItemData error: %v", err.Error())
			fmt.Println()
			fmt.Printf("affected: %v", affected)
			fmt.Println()
		}
	}
}
