package test

import (
	"clearance-adapter/factory"
	"log"
)

func initStore() {
	createPlaceDB()
	setStoreData()
}

func createPlaceDB() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := factory.GetP2BrandEngine().Exec("DROP DATABASE IF EXISTS pangpang_brand_place_management;"); err != nil {
		log.Printf("createPlaceDB error: %v", err.Error())
		log.Println()
	}

	if _, err := factory.GetP2BrandEngine().Exec("CREATE DATABASE pangpang_brand_place_management;"); err != nil {
		log.Printf("createPlaceDB error: %v", err.Error())
		log.Println()
	}
}

func setStoreData() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_place_management;"); err != nil {
		log.Printf("setStoreData error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE store
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			tenant_id BIGINT(20) NOT NULL,
			code VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			manager VARCHAR(255) NOT NULL,
			tel_no VARCHAR(255),
			area VARCHAR(255) NOT NULL,
			address VARCHAR(255) NOT NULL,
			status_code VARCHAR(255) NOT NULL,
			cashier TINYINT(1) NOT NULL,
			contract_no VARCHAR(255),
			open_date VARCHAR(255) NOT NULL,
			close_date VARCHAR(255) NOT NULL,
			enable TINYINT(1) NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			version INT(11) DEFAULT '1'
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setStoreData error: %v", err.Error())
		log.Println()
	}
	sql = `
		INSERT INTO pangpang_brand_place_management.store (tenant_id, code, name, manager, tel_no, area, address, status_code, cashier, contract_no, open_date, close_date, enable, created_at, updated_at, version)
		VALUES
		(1, 'PLANT1200', '物流仓库', 'test1', '12345678910', '北京市,市辖区,东城区', 'test12', 'open', 0, '', '2019-06-28', '9999-12-31', 1, '2019-06-28 10:37:33', '2019-07-08 02:27:58', 4),
		(1, 'CEGP', 'SA-CEGP', 'jxy', '17611242222', '北京市,市辖区,东城区', '恒通商务园', 'open', 1, '', '2009-09-30', '9999-12-31', 1, '2019-07-30 08:38:56', '2019-07-30 08:38:59', 1),
		(1, 'CFGY', 'SA-CFGY', 'jxy', '17611242222', '北京市,市辖区,东城区', '恒通商务园', 'open', 1, '', '2009-09-30', '9999-12-31', 1, '2019-07-30 08:38:56', '2019-07-30 08:38:59', 1);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setStoreData error: %v", err.Error())
		log.Println()
	}
}
