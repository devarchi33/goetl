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
	INSERT INTO pangpang_brand_sku_location.stock_distribute 
	(tenant_code, brand_code, type, box_no, waybill_no, shipment_location_id, receipt_location_id, created_at, colleague_id) 
	VALUES 
	('demo', 'SA', 'IN', '1010590009008', '1010590009008', 1, 2, '2019-07-30 10:56:43', 1),
	('demo', 'SA', 'IN', '1010590009014', '1010590009014', 1, 3, '2019-07-30 10:56:43', 1),
	('demo', 'Q3', 'IN', '1010590009009', '1010590009009', 1, 2, '2019-07-30 10:56:43', 1);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setStockDistributeData error: %v", err.Error())
		log.Println()
	}
}

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
		INSERT INTO pangpang_brand_sku_location.stock_distribute_item 
		(stock_distribute_id, product_id, brand_code, sku_id, barcode, quantity, created_at)
		VALUES
		(1, 2, 'SA', 8, 'SPWJ948S2255070', 4, '2019-07-30 10:58:14'),
		(1, 2, 'SA', 9, 'SPWJ948S2255075', 3, '2019-07-30 10:58:14'),
		(1, 2, 'SA', 10, 'SPWJ948S2255080', 2, '2019-07-30 10:58:14'),
		(1, 2, 'SA', 12, 'SPWJ948S2256075', 3, '2019-07-30 10:58:14'),
		(1, 3, 'SA', 15, 'SPWJ948S2355070', 4, '2019-07-30 10:58:14'),
		(1, 3, 'SA', 17, 'SPWJ948S2355080', 2, '2019-07-30 10:58:14'),
		(1, 3, 'SA', 18, 'SPWJ948S2356070', 4, '2019-07-30 10:58:14'),
		(1, 3, 'SA', 19, 'SPWJ948S2356075', 3, '2019-07-30 10:58:14'),
		(1, 4, 'SA', 22, 'SPYC949H2130095', 4, '2019-07-30 10:58:14'),
		(1, 4, 'SA', 23, 'SPYC949H2130100', 3, '2019-07-30 10:58:14'),
		(1, 4, 'SA', 24, 'SPYC949H2130105', 3, '2019-07-30 10:58:14'),
		(1, 4, 'SA', 29, 'SPYC949H2159095', 5, '2019-07-30 10:58:14'),
		(1, 4, 'SA', 30, 'SPYC949H2159100', 2, '2019-07-30 10:58:14'),
		(1, 2, 'SA', 1, 'SPYC949S1139085', 3, '2019-07-30 10:58:14'),
		(2, 1, 'SA', 1, 'SPYC949S1139085', 4, '2019-07-30 10:58:14'),
		(2, 1, 'SA', 2, 'SPYC949S1139090', 4, '2019-07-30 10:58:14'),
		(2, 1, 'SA', 3, 'SPYC949S1139095', 4, '2019-07-30 10:58:14'),
		(2, 1, 'SA', 4, 'SPYC949S1159085', 4, '2019-07-30 10:58:14'),
		(2, 1, 'SA', 5, 'SPYC949S1159090', 4, '2019-07-30 10:58:14'),
		(2, 1, 'SA', 6, 'SPYC949S1159095', 4, '2019-07-30 10:58:14'),
		(2, 5, 'SA', 35, 'SPYS949H2250095', 4, '2019-07-30 10:58:14'),
		(2, 5, 'SA', 36, 'SPYS949H2250100', 4, '2019-07-30 10:58:14'),
		(2, 5, 'SA', 37, 'SPYS949H2250105', 4, '2019-07-30 10:58:14'),
		(3, 6, 'Q3', 45, 'Q3AFAFDU6S2100230', 3, '2019-07-30 10:58:14'),
		(3, 6, 'Q3', 46, 'Q3AFAFDU6S2100240', 2, '2019-07-30 10:58:14'),
		(3, 6, 'Q3', 47, 'Q3AFAFDU6S2100250', 1, '2019-07-30 10:58:14'),
		(3, 6, 'Q3', 48, 'Q3AFAFDU6S2100260', 4, '2019-07-30 10:58:14'),
		(3, 6, 'Q3', 49, 'Q3AFAFDU6S2100270', 5, '2019-07-30 10:58:14');
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setStockDistributeItemData error: %v", err.Error())
		log.Println()
	}
}
