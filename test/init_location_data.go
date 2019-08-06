package test

import (
	"clearance-adapter/factory"
	"log"
)

func initLocation() {
	createLocationDB()
	// createStockDistributeTable()
	setStockDistributeData()
	// createStockDistributeItemTable()
	setStockDistributeItemData()
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

func setStockDistributeData() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_sku_location;"); err != nil {
		log.Printf("setStockDistributeData error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE stock_distribute
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			tenant_code VARCHAR(255),
			brand_code VARCHAR(255),
			type VARCHAR(255),
			box_no VARCHAR(255),
			waybill_no VARCHAR(255),
			shipment_location_id BIGINT(20),
			receipt_location_id BIGINT(20),
			created_at DATETIME,
			colleague_id BIGINT(20)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setStockDistributeData error: %v", err.Error())
		log.Println()
	}

	sql = `
	INSERT INTO pangpang_brand_sku_location.stock_distribute (tenant_code, brand_code, type, box_no, waybill_no, shipment_location_id, receipt_location_id, created_at, colleague_id) 
	VALUES 
	('demo', 'SA', 'IN', '1010590009008', '1010590009008', 1, 2, '2019-07-30 10:56:43', 1),
	('demo', 'SA', 'IN', '1010590009014', '1010590009014', 1, 3, '2019-07-30 10:56:43', 1);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setStockDistributeData error: %v", err.Error())
		log.Println()
	}
}

// func setStockDistributeData() {
// filename, err := filepath.Abs("test/data/test_in_storage_etl_stock_distribute_data.csv")
// if err != nil {
// 	panic(err.Error())
// }
// headers, data := readDataFromCSV(filename)
// transactions := buildStockDistribute(headers, data)

// loadStockDistributeData(transactions)

// }

// func buildStockDistribute(headers map[int]string, data [][]string) []models.StockDistribute {
// 	transactions := make([]models.StockDistribute, 0)
// 	for _, row := range data {
// 		txn := new(models.StockDistribute)
// 		setObjectValue(headers, row, txn)
// 		local, _ := time.LoadLocation("Local")
// 		t, _ := time.ParseInLocation("2006-01-02 15:04:05", row[4], local)
// 		txn.CreatedAt = t
// 		transactions = append(transactions, *txn)
// 	}

// 	return transactions
// }

// func loadStockDistributeData(transactions []models.StockDistribute) {
// 	for _, txn := range transactions {
// 		if affected, err := factory.GetP2BrandEngine().Insert(&txn); err != nil {
// 			fmt.Printf("loadStockDistributeData error: %v", err.Error())
// 			fmt.Println()
// 			fmt.Printf("affected: %v", affected)
// 			fmt.Println()
// 		}
// 	}
// }

func setStockDistributeItemData() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_sku_location;"); err != nil {
		log.Printf("setStockDistributeItemData error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE stock_distribute_item
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			stock_distribute_id BIGINT(20),
			product_id BIGINT(20),
			brand_code VARCHAR(255),
			sku_id BIGINT(20),
			barcode VARCHAR(255),
			quantity BIGINT(20),
			created_at DATETIME
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setStockDistributeItemData error: %v", err.Error())
		log.Println()
	}

	sql = `
		INSERT INTO pangpang_brand_sku_location.stock_distribute_item (stock_distribute_id, product_id, brand_code, sku_id, barcode, quantity, created_at)
		VALUES
		(1, 2, 'SA', 8, 'SA0001', 2, '2019-07-30 10:58:14'),
		(1, 2, 'SA', 9, 'SA0001', 2, '2019-07-30 10:58:14'),
		(1, 2, 'SA', 10, 'SA0001', 3, '2019-07-30 10:58:14'),
		(1, 2, 'SA', 1, 'SA0001', 3, '2019-07-30 10:58:14'),
		(1, 2, 'SA', 12, 'SA0001', 3, '2019-07-30 10:58:14'),
		(1, 3, 'SA', 15, 'SA0001', 3, '2019-07-30 10:58:14'),
		(1, 3, 'SA', 17, 'SA0001', 3, '2019-07-30 10:58:14'),
		(1, 3, 'SA', 18, 'SA0001', 3, '2019-07-30 10:58:14'),
		(1, 3, 'SA', 19, 'SA0001', 4, '2019-07-30 10:58:14'),
		(1, 4, 'SA', 22, 'SA0001', 4, '2019-07-30 10:58:14'),
		(1, 4, 'SA', 23, 'SA0001', 4, '2019-07-30 10:58:14'),
		(1, 4, 'SA', 24, 'SA0001', 4, '2019-07-30 10:58:14'),
		(1, 4, 'SA', 29, 'SA0001', 4, '2019-07-30 10:58:14'),
		(1, 4, 'SA', 30, 'SA0001', 4, '2019-07-30 10:58:14'),
		(2, 1, 'SA', 1, 'SA0001', 4, '2019-07-30 10:58:14'),
		(2, 1, 'SA', 2, 'SA0001', 4, '2019-07-30 10:58:14'),
		(2, 1, 'SA', 3, 'SA0001', 4, '2019-07-30 10:58:14'),
		(2, 1, 'SA', 4, 'SA0001', 4, '2019-07-30 10:58:14'),
		(2, 1, 'SA', 5, 'SA0001', 4, '2019-07-30 10:58:14'),
		(2, 1, 'SA', 6, 'SA0001', 4, '2019-07-30 10:58:14'),
		(2, 5, 'SA', 35, 'SA0001', 4, '2019-07-30 10:58:14'),
		(2, 5, 'SA', 36, 'SA0001', 4, '2019-07-30 10:58:14'),
		(2, 5, 'SA', 37, 'SA0001', 4, '2019-07-30 10:58:14');
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setStockDistributeItemData error: %v", err.Error())
		log.Println()
	}
}

// func setStockDistributeItemData() {
// 	filename, err := filepath.Abs("test/data/test_in_storage_etl_stock_distribute_item_data.csv")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	headers, data := readDataFromCSV(filename)
// 	transactionItems := buildStockDistributeItem(headers, data)

// 	loadStockDistributeItemData(transactionItems)
// }

// func buildStockDistributeItem(headers map[int]string, data [][]string) []models.StockDistributeItem {
// 	transactionItems := make([]models.StockDistributeItem, 0)
// 	for _, row := range data {
// 		txn := new(models.StockDistributeItem)
// 		setObjectValue(headers, row, txn)
// 		local, _ := time.LoadLocation("Local")
// 		t, _ := time.ParseInLocation("2006-01-02 15:04:05", row[6], local)
// 		txn.CreatedAt = t
// 		transactionItems = append(transactionItems, *txn)
// 	}

// 	return transactionItems
// }

// func loadStockDistributeItemData(transactionItems []models.StockDistributeItem) {
// 	for _, txn := range transactionItems {
// 		if affected, err := factory.GetP2BrandEngine().Insert(&txn); err != nil {
// 			fmt.Printf("loadStockDistributeItemData error: %v", err.Error())
// 			fmt.Println()
// 			fmt.Printf("affected: %v", affected)
// 			fmt.Println()
// 		}
// 	}
// }
