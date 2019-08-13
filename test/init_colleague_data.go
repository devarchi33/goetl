package test

import (
	"clearance-adapter/factory"
	"log"
)

func initColleague() {
	createEmployeeDB()
	setEmployeeData()
}

func createEmployeeDB() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := factory.GetP2BrandEngine().Exec("DROP DATABASE IF EXISTS pangpang_common_colleague_employee;"); err != nil {
		log.Printf("createEmployeeDB error: %v", err.Error())
		log.Println()
	}

	if _, err := factory.GetP2BrandEngine().Exec("CREATE DATABASE pangpang_common_colleague_employee;"); err != nil {
		log.Printf("createEmployeeDB error: %v", err.Error())
		log.Println()
	}
}

func setEmployeeData() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_common_colleague_employee;"); err != nil {
		log.Printf("setEmployeeData error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE employees
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			name VARCHAR(255),
			employee_no VARCHAR(255),
			mobile VARCHAR(255),
			email VARCHAR(255),
			enable TINYINT(1),
			created_at DATETIME,
			updated_at DATETIME,
			version INT(11) DEFAULT '1'
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setEmployeeData error: %v", err.Error())
		log.Println()
	}
	sql = `
		INSERT INTO pangpang_common_colleague_employee.employees
		(name, employee_no, mobile, email, enable, created_at, updated_at, version)
		VALUES
		('史妍珣', '7000028260', '1333333333', '33@qq.com', 1, '2019-06-27 08:44:29', '2019-06-27 08:44:31', 1);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setEmployeeData error: %v", err.Error())
		log.Println()
	}
}
