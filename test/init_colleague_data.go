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

	if _, err := factory.GetP2BrandEngine().Exec("DROP DATABASE IF EXISTS pangpang_common_colleague_auth;"); err != nil {
		log.Printf("createEmployeeDB error: %v", err.Error())
		log.Println()
	}

	if _, err := factory.GetP2BrandEngine().Exec("CREATE DATABASE pangpang_common_colleague_auth;"); err != nil {
		log.Printf("createEmployeeDB error: %v", err.Error())
		log.Println()
	}
}

func setEmployeeData() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_common_colleague_auth;"); err != nil {
		log.Printf("setEmployeeData error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE employees
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			colleague_id BIGINT(20),
			emp_id VARCHAR(255),
			name VARCHAR(255),
			mobile VARCHAR(255),
			email VARCHAR(255),
			enable TINYINT(1)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setEmployeeData error: %v", err.Error())
		log.Println()
	}
	sql = `
		INSERT INTO pangpang_common_colleague_auth.employees
		(colleague_id, emp_id, name, mobile, email, enable)
		VALUES
		(1, '7000028260', '史妍珣', '1333333333', '33@qq.com', 1);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setEmployeeData error: %v", err.Error())
		log.Println()
	}
}
