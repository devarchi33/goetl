package test

import (
	"clearance-adapter/factory"
	"log"
)

func initClearance() {
	createClearanceDB()
	createStockDistributionErrorTable()
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
			error_message VARCHAR(2000) NOT NULL DEFAULT '',
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
