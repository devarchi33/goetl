package test

import (
	"clearance-adapter/factory"
	"log"
)

func initClearance() {
	createClearanceDB()
	createStockDistributionErrorTable()
	createReturnToWarehouseErrorTable()
	createStockRoundErrorTable()
}

func createClearanceDB() {
	if _, err := factory.GetClrEngine().Exec("DROP DATABASE IF EXISTS clearance;"); err != nil {
		log.Printf("createClearanceDB error: %v", err.Error())
		log.Println()
	}
	if _, err := factory.GetClrEngine().Exec("CREATE DATABASE clearance;"); err != nil {
		log.Printf("createClearanceDB error: %v", err.Error())
		log.Println()
	}
}

func createStockDistributionErrorTable() {
	session := factory.GetClrEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE clearance;"); err != nil {
		log.Printf("createStockDistributionErrorTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE stock_distribution_error
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			brand_code VARCHAR(20) NOT NULL DEFAULT '',
			receipt_location_code VARCHAR(20) NOT NULL DEFAULT '',
			waybill_no VARCHAR(30) NOT NULL DEFAULT '',
			type VARCHAR(50) NOT NULL DEFAULT '',
			error_message VARCHAR(4000) NOT NULL DEFAULT '',
			is_processed TINYINT(1) NOT NULL DEFAULT 0 ,
			created_at DATETIME NOT NULL DEFAULT NOW(),
			created_by VARCHAR(50) NOT NULL DEFAULT '',
			UNIQUE INDEX uidx_waybill_no (brand_code, receipt_location_code, waybill_no)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockDistributionErrorTable error: %v", err.Error())
		log.Println()
	}
}

func createReturnToWarehouseErrorTable() {
	session := factory.GetClrEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE clearance;"); err != nil {
		log.Printf("createReturnToWarehouseErrorTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE return_to_warehouse_error
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			brand_code VARCHAR(20) NOT NULL DEFAULT '',
			shipment_location_code VARCHAR(20) NOT NULL DEFAULT '',
			waybill_no VARCHAR(30) NOT NULL DEFAULT '',
			type VARCHAR(50) NOT NULL DEFAULT '',
			error_message VARCHAR(4000) NOT NULL DEFAULT '',
			is_processed TINYINT(1) NOT NULL DEFAULT 0 ,
			created_at DATETIME NOT NULL DEFAULT NOW(),
			created_by VARCHAR(50) NOT NULL DEFAULT '',
			UNIQUE INDEX uidx_waybill_no (brand_code, shipment_location_code, waybill_no)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createReturnToWarehouseErrorTable error: %v", err.Error())
		log.Println()
	}
}

func createStockRoundErrorTable() {
	session := factory.GetClrEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE clearance;"); err != nil {
		log.Printf("createStockRoundErrorTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE stock_round_error
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			brand_code VARCHAR(20) NOT NULL DEFAULT '',
			shipment_location_code VARCHAR(20) NOT NULL DEFAULT '',
			receipt_location_code VARCHAR(20) NOT NULL DEFAULT '',
			waybill_no VARCHAR(30) NOT NULL DEFAULT '',
			type VARCHAR(50) NOT NULL DEFAULT '',
			error_message VARCHAR(4000) NOT NULL DEFAULT '',
			is_processed TINYINT(1) NOT NULL DEFAULT 0 ,
			created_at DATETIME NOT NULL DEFAULT NOW(),
			created_by VARCHAR(50) NOT NULL DEFAULT '',
			UNIQUE INDEX uidx_waybill_no (brand_code, shipment_location_code, receipt_location_code, waybill_no)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createStockRoundErrorTable error: %v", err.Error())
		log.Println()
	}
}
