package test

import (
	"clearance-adapter/factory"
	"log"
)

func initLocation() {
	createLocationDB()
	// 入库分配的表
	createStockDistributeTable()
	createStockDistributeItemTable()
	// 退仓的表
	createReturnToWarehouseTable()
	createReturnToWarehouseItemTable()
	// 调货的表
	createStockRoundTable()
	createStockRoundItemTable()
	// 物流中心的表
	createCMMPlantTable()
	// 库存相关表
	createStockTable()
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

func createStockDistributeTable() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_sku_location;"); err != nil {
		log.Printf("createStockDistributeTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE stock_distribute
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			tenant_code VARCHAR(255),
			brand_code VARCHAR(255),
			box_no VARCHAR(255),
			waybill_no VARCHAR(255),
			shipment_location_code VARCHAR(255),
			receipt_location_id BIGINT(20),
			created_at DATETIME,
			colleague_id BIGINT(20),
			version VARCHAR(255),
			synced TINYINT(1),
			last_updated_at DATETIME
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockDistributeTable error: %v", err.Error())
		log.Println()
	}

	sql = `
		INSERT INTO pangpang_brand_sku_location.stock_distribute 
		(tenant_code, brand_code, box_no, waybill_no, shipment_location_code, receipt_location_id, created_at, colleague_id, version, synced, last_updated_at) 
		VALUES 
		('pangpang', 'SA', '1010590009008', '1010590009008', 'ES', 2, '2019-07-30 10:56:43', 1, NULL, false, '2019-07-30 10:56:43'),
		('pangpang', 'SA', '1010590009014', '1010590009014', 'ES', 3, '2019-07-30 10:56:43', 1, NULL, false, '2019-07-30 10:56:43'),
		('pangpang', 'Q3', '1010590009009', '1010590009009', 'ES', 2, '2019-07-30 10:56:43', 1, NULL, false, '2019-07-30 10:56:43'),
		('pangpang', 'SA', '1010590009007', '1010590009007', 'ES', 2, '2019-07-30 10:56:43', 1, NULL, true, '2019-07-30 10:56:43'),
		('pangpang', 'SA', '1010590009015', '1010590009015', 'ES', 3, '2019-08-19 22:56:43', 1, NULL, false, '2019-08-19 10:56:43'),
		('pangpang', 'SA', '1010590009016', '1010590009016', 'ES', 3, '2019-08-19 10:56:43', 1, NULL, false, '2019-08-19 10:56:43');
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockDistributeTable error: %v", err.Error())
		log.Println()
	}
}

func createStockDistributeItemTable() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_sku_location;"); err != nil {
		log.Printf("createStockDistributeItemTable error: %v", err.Error())
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
			created_at DATETIME,
			out_quantity BIGINT(20)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockDistributeItemTable error: %v", err.Error())
		log.Println()
	}

	sql = `
		INSERT INTO pangpang_brand_sku_location.stock_distribute_item 
		(stock_distribute_id, product_id, brand_code, sku_id, barcode, quantity, created_at, out_quantity)
		VALUES
		(1, 2, 'SA', 8, 'SPWJ948S2255070', 4, '2019-07-30 10:58:14', 4),
		(1, 2, 'SA', 9, 'SPWJ948S2255075', 3, '2019-07-30 10:58:14', 4),
		(1, 2, 'SA', 10, 'SPWJ948S2255080', 2, '2019-07-30 10:58:14', 2),
		(1, 2, 'SA', 12, 'SPWJ948S2256075', 3, '2019-07-30 10:58:14', 3),
		(1, 3, 'SA', 15, 'SPWJ948S2355070', 4, '2019-07-30 10:58:14', 3),
		(1, 3, 'SA', 17, 'SPWJ948S2355080', 2, '2019-07-30 10:58:14', 3),
		(1, 3, 'SA', 18, 'SPWJ948S2356070', 4, '2019-07-30 10:58:14', 3),
		(1, 3, 'SA', 19, 'SPWJ948S2356075', 3, '2019-07-30 10:58:14', 3),
		(1, 4, 'SA', 22, 'SPYC949H2130095', 4, '2019-07-30 10:58:14', 3),
		(1, 4, 'SA', 23, 'SPYC949H2130100', 3, '2019-07-30 10:58:14', 3),
		(1, 4, 'SA', 24, 'SPYC949H2130105', 3, '2019-07-30 10:58:14', 3),
		(1, 4, 'SA', 29, 'SPYC949H2159095', 5, '2019-07-30 10:58:14', 3),
		(1, 4, 'SA', 30, 'SPYC949H2159100', 2, '2019-07-30 10:58:14', 3),
		(1, 2, 'SA', 1, 'SPYC949S1139085', 3, '2019-07-30 10:58:14', 3),
		(2, 1, 'SA', 1, 'SPYC949S1139085', 4, '2019-07-30 10:58:14', 3),
		(2, 1, 'SA', 2, 'SPYC949S1139090', 4, '2019-07-30 10:58:14', 3),
		(2, 1, 'SA', 3, 'SPYC949S1139095', 4, '2019-07-30 10:58:14', 3),
		(2, 1, 'SA', 4, 'SPYC949S1159085', 4, '2019-07-30 10:58:14', 3),
		(2, 1, 'SA', 5, 'SPYC949S1159090', 4, '2019-07-30 10:58:14', 3),
		(2, 1, 'SA', 6, 'SPYC949S1159095', 4, '2019-07-30 10:58:14', 3),
		(2, 5, 'SA', 35, 'SPYS949H2250095', 4, '2019-07-30 10:58:14', 3),
		(2, 5, 'SA', 36, 'SPYS949H2250100', 4, '2019-07-30 10:58:14', 3),
		(2, 5, 'SA', 37, 'SPYS949H2250105', 4, '2019-07-30 10:58:14', 3),
		(3, 6, 'Q3', 45, 'Q3AFAFDU6S2100230', 3, '2019-07-30 10:58:14', 3),
		(3, 6, 'Q3', 46, 'Q3AFAFDU6S2100240', 2, '2019-07-30 10:58:14', 3),
		(3, 6, 'Q3', 47, 'Q3AFAFDU6S2100250', 1, '2019-07-30 10:58:14', 3),
		(3, 6, 'Q3', 48, 'Q3AFAFDU6S2100260', 4, '2019-07-30 10:58:14', 3),
		(3, 6, 'Q3', 49, 'Q3AFAFDU6S2100270', 5, '2019-07-30 10:58:14', 3),
		(4, 2, 'SA', 8, 'SPWJ948S2255070', 4, '2019-08-19 10:58:14', 3),
		(4, 2, 'SA', 9, 'SPWJ948S2255075', 3, '2019-08-19 10:58:14', 3),
		(5, 1, 'SA', 1, 'SPYC949S1139085', 4, '2019-08-19 10:58:14', 3),
		(5, 1, 'SA', 2, 'SPYC949S1139090', 4, '2019-08-19 10:58:14', 3),
		(5, 1, 'SA', 3, 'SPYC949S1139095', 4, '2019-08-19 10:58:14', 3),
		(5, 1, 'SA', 4, 'SPYC949S1159085', 4, '2019-08-19 10:58:14', 3),
		(5, 1, 'SA', 5, 'SPYC949S1159090', 4, '2019-08-19 10:58:14', 3),
		(5, 1, 'SA', 6, 'SPYC949S1159095', 4, '2019-08-19 10:58:14', 3),
		(5, 5, 'SA', 35, 'SPYS949H2250095', 4, '2019-08-19 10:58:14', 3),
		(5, 5, 'SA', 36, 'SPYS949H2250100', 4, '2019-08-19 10:58:14', 3),
		(5, 5, 'SA', 37, 'SPYS949H2250105', 4, '2019-08-19 10:58:14', 3),
		(6, 2, 'SA', 9, 'SPWJ948S2255075', 3, '2019-07-30 10:58:14', 3),
		(6, 2, 'SA', 9, 'SPWJ948S2255075', 3, '2019-07-30 10:58:14', 3);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockDistributeItemTable error: %v", err.Error())
		log.Println()
	}
}

func createReturnToWarehouseTable() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_sku_location;"); err != nil {
		log.Printf("createReturnToWarehouseTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE return_to_warehouse
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			tenant_code VARCHAR(255),
			brand_code VARCHAR(255),
			status VARCHAR(255),
			waybill_no VARCHAR(255),
			shipment_location_id BIGINT(20),
			created_at DATETIME,
			last_updated_at DATETIME,
			colleague_id BIGINT(20),
			receipt_location_code VARCHAR(255),
			synced TINYINT(1) DEFAULT '0' NOT NULL,
			last_colleague_id BIGINT(20),
			out_created_at DATETIME,
			out_colleague_id BIGINT(20)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createReturnToWarehouseTable error: %v", err.Error())
		log.Println()
	}

	sql = `
		INSERT INTO pangpang_brand_sku_location.return_to_warehouse 
		(tenant_code, brand_code, status, waybill_no, shipment_location_id, created_at, last_updated_at, 
			colleague_id, receipt_location_code, synced, last_colleague_id, out_created_at, out_colleague_id) 
		VALUES
		('pangpang', 'SA', 'R', '20190813001', 2, '2019-08-13 09:03:12', '2019-08-13 09:03:13', 1, 'ES', false, 1, '2019-08-13 09:03:13', 1),
		('pangpang', 'Q3', 'R', '20190814001', 2, '2019-08-13 09:03:12', '2019-08-13 09:03:13', 1, 'ES', false, 1, '2019-08-13 19:03:13', 1),
		('pangpang', 'SA', 'R', '20190819001', 2, '2019-08-14 09:03:12', '2019-08-14 09:03:13', 1, 'ES', true, 1, '2019-08-18 19:03:13', 1);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createReturnToWarehouseTable error: %v", err.Error())
		log.Println()
	}
}

func createReturnToWarehouseItemTable() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_sku_location;"); err != nil {
		log.Printf("createReturnToWarehouseItemTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE return_to_warehouse_item
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			return_to_warehouse_id BIGINT(20),
			product_id BIGINT(20),
			brand_code VARCHAR(255),
			sku_id BIGINT(20),
			barcode VARCHAR(255),
			quantity BIGINT(20),
			created_at DATETIME,
			created_colleague_id BIGINT(20),
			updated_at DATETIME,
			updated_colleague_id BIGINT(20)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createReturnToWarehouseItemTable error: %v", err.Error())
		log.Println()
	}

	sql = `
		INSERT INTO pangpang_brand_sku_location.return_to_warehouse_item 
		(return_to_warehouse_id, product_id, brand_code, sku_id, barcode, quantity, created_at, created_colleague_id, updated_at, updated_colleague_id) 
		VALUES 
		(1, 2, 'SA', 8, 'SPWJ948S2255070', 4, '2019-08-13 09:09:13', 1, '2019-08-13 09:09:18', 1),
		(1, 1, 'SA', 3, 'SPYC949S1139095', 1, '2019-08-13 09:10:08', 1, '2019-08-13 09:10:11', 1),
		(2, 6, 'Q3', 45, 'Q3AFAFDU6S2100230', 2, '2019-08-14 09:09:13', 1, '2019-08-14 09:09:18', 1),
		(2, 6, 'Q3', 46, 'Q3AFAFDU6S2100240', 3, '2019-08-14 09:10:08', 1, '2019-08-14 09:10:11', 1),
		(3, 2, 'SA', 8, 'SPWJ948S2255070', 4, '2019-08-13 09:09:13', 1, '2019-08-13 09:09:18', 1),
		(3, 1, 'SA', 3, 'SPYC949S1139095', 1, '2019-08-13 09:10:08', 1, '2019-08-13 09:10:11', 1);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createReturnToWarehouseItemTable error: %v", err.Error())
		log.Println()
	}
}

func createStockRoundTable() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_sku_location;"); err != nil {
		log.Printf("createStockRoundTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE stock_round
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			tenant_code VARCHAR(255),
			brand_code VARCHAR(255),
			box_no VARCHAR(255),
			waybill_no VARCHAR(255),
			shipment_location_id BIGINT(20),
			receipt_location_id BIGINT(20),
			synced TINYINT(1),
			created_at DATETIME,
			last_updated_at DATETIME,
			colleague_id BIGINT(20),
			status VARCHAR(255),
			in_created_at DATETIME,
			out_created_at DATETIME,
			in_colleague_id BIGINT(20),
			out_colleague_id BIGINT(20),
			shipping_company_code VARCHAR(255)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockRoundTable error: %v", err.Error())
		log.Println()
	}

	sql = `
		INSERT INTO pangpang_brand_sku_location.stock_round 
		(tenant_code, brand_code, box_no, waybill_no, shipment_location_id, receipt_location_id, synced, created_at, last_updated_at, colleague_id, status, in_created_at, out_created_at, in_colleague_id, out_colleague_id, shipping_company_code) 
		VALUES 
		('pangpang', 'SA', '20190821001', '20190821001', 2, 3, 0, '2019-08-20 13:42:13', '2019-08-20 15:02:09', 1, 'F', '2019-08-21 22:02:09', '2019-08-20 22:02:09', 1, 1, 'SR');
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockRoundTable error: %v", err.Error())
		log.Println()
	}
}

func createStockRoundItemTable() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_sku_location;"); err != nil {
		log.Printf("createStockRoundItemTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE stock_round_item
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			stock_round_id BIGINT(20),
			product_id BIGINT(20),
			brand_code VARCHAR(255),
			sku_id BIGINT(20),
			barcode VARCHAR(255),
			quantity BIGINT(20),
			created_at DATETIME,
			created_colleague_id BIGINT(20),
			updated_at DATETIME,
			updated_colleague_id BIGINT(20)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockRoundItemTable error: %v", err.Error())
		log.Println()
	}

	sql = `
		INSERT INTO pangpang_brand_sku_location.stock_round_item 
		(stock_round_id, product_id, brand_code, sku_id, barcode, quantity, created_at, created_colleague_id, updated_at, updated_colleague_id)
		VALUES 
		(1, 2, 'SA', 8, 'SPWJ948S2255070', 4, '2019-08-21 13:42:13', 1, '2019-08-21 13:42:13', 1),
		(1, 1, 'SA', 3, 'SPYC949S1139095', 1, '2019-08-21 13:42:13', 1, '2019-08-21 13:42:13', 1);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockRoundItemTable error: %v", err.Error())
		log.Println()
	}
}

func createCMMPlantTable() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_sku_location;"); err != nil {
		log.Printf("createCMMPlantTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE cmm_plant
		(
			id BIGINT(20) PRIMARY KEY NOT NULL,
			plant_code VARCHAR(255),
			plant_name VARCHAR(255),
			user_name VARCHAR(255),
			mobile VARCHAR(50),
			area VARCHAR(255),
			address VARCHAR(255),
			created_at DATETIME,
			uaer_name VARCHAR(255)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createCMMPlantTable error: %v", err.Error())
		log.Println()
	}

	sql = `
		INSERT INTO pangpang_brand_sku_location.cmm_plant 
		(id, plant_code, plant_name, user_name, mobile, area, address, created_at, uaer_name) 
		VALUES (1, 'ES', 'ES', '张殿利', '13917778786', '上海市,上海市,闵行区', '莲花南路3130号 衣恋物流中心', '2019-08-09 09:32:19', null);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createCMMPlantTable error: %v", err.Error())
		log.Println()
	}
}

func createStockTable() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_sku_location;"); err != nil {
		log.Printf("createStockTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE stock_history
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			stock_id BIGINT(20),
			quantity BIGINT(20),
			operation VARCHAR(255),
			current_quantity BIGINT(20),
			colleague_id BIGINT(20),
			created_at DATETIME,
			tenant_code VARCHAR(255),
			order_id BIGINT(20),
			event_type VARCHAR(255),
			event_id BIGINT(20)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockTable error: %v", err.Error())
		log.Println()
	}

	sql = `
		CREATE TABLE stock
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			location_id BIGINT(20),
			product_id BIGINT(20),
			sku_id BIGINT(20),
			quantity BIGINT(20),
			barcode VARCHAR(255),
			tenant_code VARCHAR(255),
			created_at DATETIME,
			updated_at DATETIME
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockTable error: %v", err.Error())
		log.Println()
	}
}
