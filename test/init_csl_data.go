package test

import (
	"clearance-adapter/factory"
	"fmt"
	"log"
)

func initCSL() {
	createRecvSuppMstTable()
	createRecvSuppDtlTable()
	createStockMisDtlTable()
	createProductTable()
	createMonthlyBizFuctionClosingTable()
	createIFConfigTable()
	createUserInfoTable()
	createEmployeeTable()
	createComplexShopMappingTable()
	createSP()

	initReturnToWarehouseData()
	initTransferData()
}

// 创建的时候有问题，暂未使用，需要事先创建CSL数据库
func createCSLDB() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()
	if _, err := session.Exec("USE master;"); err != nil {
		log.Printf("createCSLDB error: %v", err.Error())
		log.Println()
	}
	if _, err := session.Exec("DROP DATABASE IF EXISTS CSL;"); err != nil {
		log.Printf("createCSLDB error: %v", err.Error())
		log.Println()
	}
	if _, err := session.Exec("CREATE DATABASE CSL;"); err != nil {
		log.Printf("createCSLDB error: %v", err.Error())
		log.Println()
	}
}

// 创建一些SP中用到的表
func createMonthlyBizFuctionClosingTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createMonthlyBizFuctionClosingTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.MonthlyBizFuctionClosing;"); err != nil {
		log.Printf("createMonthlyBizFuctionClosingTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE CSL.dbo.MonthlyBizFuctionClosing
		(
			BizFuctionID CHAR(2) NOT NULL,
			ClosingYearMonth CHAR(6) NOT NULL,
			ClosingSeq BIT NOT NULL,
			StartDateTime DATETIME,
			EndDateTime DATETIME,
			InformKr NVARCHAR(500),
			InformZn NVARCHAR(500),
			ClosingChk BIT DEFAULT '1' NOT NULL,
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME
		);
	`
	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createMonthlyBizFuctionClosingTable error: %v", err.Error())
		fmt.Println()
	}
}

func createIFConfigTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createIFConfigTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.IFConfig;"); err != nil {
		log.Printf("createIFConfigTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE CSL.dbo.IFConfig
		(
			IFSendFlagChk BIT,
    		TopCount INT
		);
		INSERT INTO CSL.dbo.IFConfig (IFSendFlagChk, TopCount) VALUES (1, 1000);
	`
	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createIFConfigTable error: %v", err.Error())
		fmt.Println()
	}
}

func createUserInfoTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createUserInfoTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.UserInfo;"); err != nil {
		log.Printf("createUserInfoTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE CSL.dbo.UserInfo
		(
			UserID VARCHAR(20) PRIMARY KEY NOT NULL,
			EmpID CHAR(10) NOT NULL,
			UserName NVARCHAR(100) NOT NULL,
			LoginID VARCHAR(25) NOT NULL,
			Domain VARCHAR(20),
			EndDate CHAR(8),
			UseChk BIT,
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME,
			Password VARCHAR(20),
			EuroIFChk BIT,
			ShopCode CHAR(4)
		);
		INSERT INTO CSL.dbo.UserInfo (UserID, EmpID, UserName, LoginID, Domain, EndDate, UseChk, InUserID, InDateTime, Password, EuroIFChk, ShopCode) 
		VALUES ('shi.yanxun', '7000028260', N'史妍珣', 'Shop-7000028260', 'CHINA', null, 1, 'LI_JING', '2015-09-07 17:07:00.000', '7000028260', null, null);
	`
	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createUserInfoTable error: %v", err.Error())
		fmt.Println()
	}
}

func createEmployeeTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createEmployeeTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.Employee;"); err != nil {
		log.Printf("createEmployeeTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE CSL.dbo.Employee
		(
			EmpID CHAR(10) PRIMARY KEY NOT NULL,
			EmpName NVARCHAR(100),
			HREmpNo CHAR(20),
			AreaCode VARCHAR(10),
			BranchCode VARCHAR(4),
			FactoryID VARCHAR(10),
			Birthday CHAR(8),
			SolarChk BIT,
			CelNo VARCHAR(25),
			ConfirmDate VARCHAR(8),
			ConfirmEmpID CHAR(10),
			ConfirmTypeCode CHAR,
			EmailAddress NVARCHAR(200),
			EnterDate VARCHAR(8),
			Address NVARCHAR(400),
			ZipNo VARCHAR(6),
			FamilyAddress NVARCHAR(400),
			FamilyZipNo VARCHAR(6),
			RetireDate VARCHAR(8),
			SexCode CHAR,
			ShopEmpChk BIT,
			SocialNo VARCHAR(20),
			TelNo VARCHAR(25),
			WeddingDay VARCHAR(8),
			EmpTypeCode CHAR,
			PeopleCode VARCHAR(2),
			CountryCode VARCHAR(2),
			DeptCode CHAR(2),
			ScholarshipCode VARCHAR(2),
			PositionCode CHAR(2),
			BrandCode VARCHAR(4),
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME,
			ModiUserID VARCHAR(20) NOT NULL,
			ModiDateTime DATETIME,
			SendState VARCHAR(2) DEFAULT '' NOT NULL,
			SendFlag CHAR DEFAULT 'R' NOT NULL,
			SendSeqNo BIGINT NOT NULL IDENTITY,
			SendDateTime DATETIME,
			AppInstallChk BIT,
			VMDAppReg_ID VARCHAR(255),
			KIS_ScholarshipName NVARCHAR(200),
			OldEnterDate VARCHAR(8),
			EX_EMPTP VARCHAR(1),
			NewPhoneNumber VARCHAR(20)
		);
		INSERT INTO CSL.dbo.Employee (EmpID, EmpName, HREmpNo, AreaCode, BranchCode, FactoryID, Birthday, SolarChk, CelNo, ConfirmDate, ConfirmEmpID, ConfirmTypeCode, EmailAddress, EnterDate, Address, ZipNo, FamilyAddress, FamilyZipNo, RetireDate, SexCode, ShopEmpChk, SocialNo, TelNo, WeddingDay, EmpTypeCode, PeopleCode, CountryCode, DeptCode, ScholarshipCode, PositionCode, BrandCode, InUserID, InDateTime, ModiUserID, ModiDateTime, SendState, SendFlag, SendDateTime, AppInstallChk, VMDAppReg_ID, KIS_ScholarshipName, OldEnterDate, EX_EMPTP, NewPhoneNumber) 
		VALUES ('7000028260', N'史妍珣', null, '0000001826', 'B413', '8573', '19860919', 1, '13677694277', '20150903', '3000000745', 'E', null, '20150903', '0', null, null, null, null, 'F', 1, '500236198609190863', null, null, 'C', '01', 'CN', 'B ', null, '20', 'EK', 'LI_JING', '2015-09-07 17:07:00.000', 'diao_xianyu', '2016-06-01 16:30:03.603', '', 'R', null, null, null, null, null, null, null);
	`
	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createEmployeeTable error: %v", err.Error())
		fmt.Println()
	}
}

func createComplexShopMappingTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createComplexShopMappingTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.ComplexShopMapping;"); err != nil {
		log.Printf("createComplexShopMappingTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE CSL.dbo.ComplexShopMapping
		(
			BrandCode VARCHAR(4) NOT NULL,
			ShopCode CHAR(4) NOT NULL,
			ChiefBrandCode VARCHAR(4) NOT NULL,
			ChiefShopCode CHAR(4) NOT NULL,
			DelChk BIT NOT NULL,
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME
		);

		INSERT INTO CSL.dbo.ComplexShopMapping (BrandCode, ShopCode, ChiefBrandCode, ChiefShopCode, DelChk, InUserID, InDateTime) 
		VALUES 
		('MC', 'CGND', 'SA', 'CG4R', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CGTJ', 'SA', 'CEGP', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CGXU', 'SA', 'CG4U', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CH69', 'SA', 'CFZV', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CH6D', 'SA', 'CFRW', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CHRR', 'SA', 'CEYL', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CHSB', 'SA', 'CFGY', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CJ29', 'SA', 'CDGT', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CJD7', 'SA', 'CCX4', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CJD8', 'SA', 'CEKG', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CJDE', 'SA', 'CFTN', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CJJN', 'SA', 'CJJL', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CJJU', 'SA', 'CJJS', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CJRG', 'SA', 'CFW5', 0, 'system', '2019-08-09 05:17:32.373'),
		('MC', 'CJXA', 'SA', 'CJJD', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CH0D', 'SA', 'CG4R', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CH71', 'SA', 'CFZV', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CH7A', 'SA', 'CFRW', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CHGA', 'SA', 'CG4U', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CHTM', 'SA', 'CFGY', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CHZM', 'SA', 'CEYL', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CJ3B', 'SA', 'CDGT', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CJC1', 'SA', 'CEGP', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CJDF', 'SA', 'CFTN', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CJDM', 'SA', 'CEKG', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CJFV', 'SA', 'CCX4', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CJJM', 'SA', 'CJJL', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CJJT', 'SA', 'CJJS', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CJRH', 'SA', 'CFW5', 0, 'system', '2019-08-09 05:17:32.373'),
		('Q3', 'CJXB', 'SA', 'CJJD', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CCX4', 'SA', 'CCX4', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CDGT', 'SA', 'CDGT', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CEGP', 'SA', 'CEGP', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CEKG', 'SA', 'CEKG', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CEYL', 'SA', 'CEYL', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CFGY', 'SA', 'CFGY', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CFRW', 'SA', 'CFRW', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CFTN', 'SA', 'CFTN', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CFW5', 'SA', 'CFW5', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CFZV', 'SA', 'CFZV', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CG4R', 'SA', 'CG4R', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CG4U', 'SA', 'CG4U', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CJJD', 'SA', 'CJJD', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CJJL', 'SA', 'CJJL', 0, 'system', '2019-08-09 05:17:32.373'),
		('SA', 'CJJS', 'SA', 'CJJS', 0, 'system', '2019-08-09 05:17:32.373');
	`
	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createComplexShopMappingTable error: %v", err.Error())
		fmt.Println()
	}
}

// RecvSuppMaster 部分
func createRecvSuppMstTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createRecvSuppMstTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.RecvSuppMst;"); err != nil {
		log.Printf("createRecvSuppMstTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE CSL.dbo.RecvSuppMst
		(
			RecvSuppNo CHAR(14) PRIMARY KEY NOT NULL,
			BrandCode VARCHAR(4),
			ShopCode CHAR(4),
			Dates CHAR(8) NOT NULL,
			SeqNo INT NOT NULL,
			SAPDeliveryNo CHAR(10),
			SAPDeliveryDate CHAR(8),
			RecvSuppType CHAR,
			NormalProductType CHAR,
			ShopSuppRecvDate CHAR(8),
			BrandSuppRecvDate VARCHAR(8),
			TransTypeCode CHAR,
			ShippingTypeCode CHAR(2),
			WayBillNo VARCHAR(13) NOT NULL,
			RecvSuppStatusCode CHAR NOT NULL,
			RequestNo CHAR(14),
			BoxNo CHAR(20),
			ShopDesc NVARCHAR(400),
			BrandDesc NVARCHAR(400),
			PlantCode CHAR(4),
			RoundRecvSuppNo CHAR(14),
			RoundSAPDeliveryNo CHAR(10),
			TargetShopCode CHAR(4),
			RecvChk BIT,
			OrderControlNo CHAR(12),
			RecvEmpID CHAR(10),
			RecvEmpName NVARCHAR(100),
			SuppEmpID CHAR(10),
			SuppEmpName NVARCHAR(200),
			SAPMenuType CHAR,
			DelChk BIT DEFAULT 0 NOT NULL,
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME,
			ModiUserID VARCHAR(20) NOT NULL,
			ModiDateTime DATETIME,
			SendState VARCHAR(2) DEFAULT '' NOT NULL,
			SendFlag CHAR DEFAULT 'R' NOT NULL,
			SendSeqNo BIGINT NOT NULL IDENTITY,
			SendDateTime DATETIME,
			InvtBaseDate CHAR(8),
			BoxAmount INT,
			StockOutUseAmt DECIMAL(9,2),
			ExpressNo VARCHAR(13),
			ShippingCompanyCode CHAR(2),
			BoxGram DECIMAL(18,3),
			DeliveryID VARCHAR(250),
			DeliveryOrderNo VARCHAR(250),
			VolumeType NVARCHAR(20),
			VolumesSize VARCHAR(20),
			VolumesUnit NVARCHAR(10),
			Channel VARCHAR(20),
			ProvinceCode CHAR(3),
			CityCode CHAR(5),
			DistrictCode CHAR(8),
			Area NVARCHAR(100),
			ShopManagerName NVARCHAR(10),
			MobilePhone VARCHAR(25),
			DeliverySendTime DATETIME,
			DeliveryReceiveTime DATETIME,
			BoxType CHAR(2)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createRecvSuppMstTable error: %v", err.Error())
		fmt.Println()
	}

	sql = `
		INSERT INTO CSL.dbo.RecvSuppMst 
		(
			RecvSuppNo, BrandCode, ShopCode, Dates, SeqNo, SAPDeliveryNo, SAPDeliveryDate, RecvSuppType, NormalProductType, ShopSuppRecvDate, BrandSuppRecvDate, TransTypeCode, ShippingTypeCode, WayBillNo, RecvSuppStatusCode, RequestNo, BoxNo, ShopDesc, BrandDesc, PlantCode, RoundRecvSuppNo, RoundSAPDeliveryNo, TargetShopCode, 
			RecvChk, OrderControlNo, RecvEmpID, RecvEmpName, SuppEmpID, SuppEmpName, SAPMenuType, DelChk, InUserID, InDateTime, ModiUserID, ModiDateTime, SendState, SendFlag, SendDateTime, InvtBaseDate, BoxAmount, StockOutUseAmt, ExpressNo, ShippingCompanyCode, BoxGram, DeliveryID, DeliveryOrderNo, VolumeType, VolumesSize, 
			VolumesUnit, Channel, ProvinceCode, CityCode, DistrictCode, Area, ShopManagerName, MobilePhone, DeliverySendTime, DeliveryReceiveTime, BoxType
		) 
		VALUES 
		('CEGP1907236002', 'SA', 'CEGP', '20190723', 6002, '8074783302', '20190723', 'R', 'A', '20190723', '20190723', '5', '01', '1010590009008', 'R', '              ', '1010590009008       ', ' ', ' ', '1200', '              ', '          ', '    ', 0, '            ', '          ', ' ', '          ', 'CNFBGLEA02', '1', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '20190723', 0, 0.00, '', 'SR', 0.000, '', '', '', '', '', '', '   ', '     ', '        ', '', '', '', null, null, '  '),
		('CFGY1907236000', 'SA', 'CFGY', '20190723', 6000, '8074783296', '20190723', 'R', 'A', '20190723', '20190723', '5', '01', '1010590009014', 'R', '              ', '1010590009014       ', ' ', ' ', '1200', '              ', '          ', '    ', 0, '            ', '          ', ' ', '          ', 'CNFBGLEA02', '1', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '20190723', 0, 0.00, '', 'SR', 0.000, '', '', '', '', '', '', '   ', '     ', '        ', '', '', '', null, null, '  '),
		('CFGY1907236010', 'SA', 'CFGY', '20190723', 6010, '8074783305', '20190723', 'R', 'A', '20190723', '20190723', '5', '01', '1010590009014', 'R', '              ', '1010590009014       ', ' ', ' ', '1200', '              ', '          ', '    ', 0, '            ', '          ', ' ', '          ', 'CNFBGLEA02', '1', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '20190723', 0, 0.00, '', 'SR', 0.000, '', '', '', '', '', '', '   ', '     ', '        ', '', '', '', null, null, '  '),
		('CJC11907236000', 'Q3', 'CJC1', '20190723', 6002, '8074783302', '20190723', 'R', 'A', '20190723', '20190723', '5', '01', '1010590009009', 'R', '              ', '1010590009009       ', ' ', ' ', '1200', '              ', '          ', '    ', 0, '            ', '          ', ' ', '          ', 'CNFBGLEA02', '1', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '20190723', 0, 0.00, '', 'SR', 0.000, '', '', '', '', '', '', '   ', '     ', '        ', '', '', '', null, null, '  '),
		('CEGP1907236003', 'SA', 'CEGP', '20190723', 6002, '8074783302', '20190723', 'R', 'A', '20190723', '20190723', '5', '01', '1010590009007', 'R', '              ', '1010590009007       ', ' ', ' ', '1200', '              ', '          ', '    ', 0, '            ', '          ', ' ', '          ', 'CNFBGLEA02', '1', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '20190723', 0, 0.00, '', 'SR', 0.000, '', '', '', '', '', '', '   ', '     ', '        ', '', '', '', null, null, '  '),
		('CFGY1908196000', 'SA', 'CFGY', '20190819', 6000, '8074783296', '20190819', 'R', 'A', '20190819', '20190819', '5', '01', '1010590009015', 'R', '              ', '1010590009015       ', ' ', ' ', '1200', '              ', '          ', '    ', 0, '            ', '          ', ' ', '          ', 'CNFBGLEA02', '1', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '20190723', 0, 0.00, '', 'SR', 0.000, '', '', '', '', '', '', '   ', '     ', '        ', '', '', '', null, null, '  '),
		('CFGY1908196010', 'SA', 'CFGY', '20190819', 6010, '8074783305', '20190819', 'R', 'A', '20190819', '20190819', '5', '01', '1010590009015', 'F', '              ', '1010590009015       ', ' ', ' ', '1200', '              ', '          ', '    ', 1, '            ', '          ', ' ', '          ', 'CNFBGLEA02', '1', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '20190819', 0, 0.00, '', 'SR', 0.000, '', '', '', '', '', '', '   ', '     ', '        ', '', '', '', null, null, '  '),
		('CFGY1908206000', 'SA', 'CFGY', '20190820', 6000, '8074783296', '20190820', 'R', 'A', '20190820', '20190820', '5', '01', '1010590009016', 'R', '              ', '1010590009016       ', ' ', ' ', '1200', '              ', '          ', '    ', 0, '            ', '          ', ' ', '          ', 'CNFBGLEA02', '1', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '20190820', 0, 0.00, '', 'SR', 0.000, '', '', '', '', '', '', '   ', '     ', '        ', '', '', '', null, null, '  '),
		('CFGY1908206010', 'SA', 'CFGY', '20190820', 6010, '8074783305', '20190820', 'R', 'A', '20190820', '20190820', '5', '01', '1010590009016', 'R', '              ', '1010590009016       ', ' ', ' ', '1200', '              ', '          ', '    ', 0, '            ', '          ', ' ', '          ', 'CNFBGLEA02', '1', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '20190820', 0, 0.00, '', 'SR', 0.000, '', '', '', '', '', '', '   ', '     ', '        ', '', '', '', null, null, '  '),
		('CJC11909026000', 'Q3', 'CJC1', '20190817', 6002, '8074783302', '20190817', 'R', 'A', '20190817', '20190817', '5', '01', '20190902001', 'R', '              ', '20190902001       ', ' ', ' ', '1200', '              ', '          ', '    ', 0, '            ', '          ', ' ', '          ', 'CNFBGLEA02', '1', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '20190817', 0, 0.00, '', 'SR', 0.000, '', '', '', '', '', '', '   ', '     ', '        ', '', '', '', null, null, '  '),
		('CJUT1909026000', 'SA', 'CJUT', '20190817', 6002, '8074783302', '20190817', 'R', 'A', '20190817', '20190817', '5', '01', '20190902002', 'R', '              ', '20190902002       ', ' ', ' ', '1200', '              ', '          ', '    ', 0, '            ', '          ', ' ', '          ', 'CNFBGLEA02', '1', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '20190817', 0, 0.00, '', 'SR', 0.000, '', '', '', '', '', '', '   ', '     ', '        ', '', '', '', null, null, '  ');
	`

	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createRecvSuppMstTable error: %v", err.Error())
		fmt.Println()
	}
}

// RecvSuppDetail 部分
func createRecvSuppDtlTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createRecvSuppDtlTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.RecvSuppDtl;"); err != nil {
		log.Printf("createRecvSuppDtlTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE CSL.dbo.RecvSuppDtl
		(
			RecvSuppNo CHAR(14) NOT NULL,
			RecvSuppSeqNo INT NOT NULL,
			SupGroupCode CHAR(2),
			BrandCode VARCHAR(4),
			ShopCode CHAR(4),
			Dates CHAR(8),
			SeqNo INT,
			SAPDeliveryNo CHAR(10),
			SAPDeliveryItemNo CHAR(10),
			RoundRecvSuppNo CHAR(14),
			RoundRecvSuppDtSeq INT,
			RoundSAPDeliveryNo CHAR(10),
			RoundSAPDeliveryItemNo CHAR(10),
			ProdCode VARCHAR(18),
			PriceTypeCode CHAR(2),
			SaipType CHAR(2),
			RecvSuppQty INT,
			RecvSuppFixedQty INT,
			SalePrice DECIMAL(19,2),
			AbnormalProdReasonCode CHAR(2),
			DelChk BIT DEFAULT 0 NOT NULL,
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME,
			ModiUserID VARCHAR(20) NOT NULL,
			ModiDateTime DATETIME,
			SendState VARCHAR(2) DEFAULT '' NOT NULL,
			SendFlag CHAR DEFAULT 'R' NOT NULL,
			SendSeqNo BIGINT NOT NULL IDENTITY,
			SendDateTime DATETIME,
			AbnormalChkCode CHAR(2),
			AbnormalSerialNo VARCHAR(7),
			ModiReason NVARCHAR(800),
			ApplyID VARCHAR(30)
		);
	`
	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createRecvSuppDtlTable error: %v", err.Error())
		fmt.Println()
	}

	sql = `
		INSERT INTO CSL.dbo.RecvSuppDtl 
		(
			RecvSuppNo, RecvSuppSeqNo, SupGroupCode, BrandCode, ShopCode, Dates, SeqNo, SAPDeliveryNo, SAPDeliveryItemNo, RoundRecvSuppNo, RoundRecvSuppDtSeq, RoundSAPDeliveryNo, RoundSAPDeliveryItemNo, ProdCode, PriceTypeCode, SaipType, 
			RecvSuppQty, RecvSuppFixedQty, SalePrice, AbnormalProdReasonCode, DelChk, InUserID, InDateTime, ModiUserID, ModiDateTime, SendState, SendFlag, SendDateTime, AbnormalChkCode, AbnormalSerialNo, ModiReason, ApplyID
		) 
		VALUES 
		('CEGP1907236002', 1, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000310    ', '              ', 0, '          ', '          ', 'SPWJ948S2255070', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236002', 2, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000320    ', '              ', 0, '          ', '          ', 'SPWJ948S2255075', '  ', '  ', 3, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236002', 3, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000330    ', '              ', 0, '          ', '          ', 'SPWJ948S2255080', '  ', '  ', 2, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236002', 4, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000340    ', '              ', 0, '          ', '          ', 'SPWJ948S2256070', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236002', 5, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000350    ', '              ', 0, '          ', '          ', 'SPWJ948S2256075', '  ', '  ', 3, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236002', 9, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000400    ', '              ', 0, '          ', '          ', 'SPWJ948S2356075', '  ', '  ', 3, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236002', 10, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000470    ', '              ', 0, '          ', '          ', 'SPYC949H2130095', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236002', 11, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000480    ', '              ', 0, '          ', '          ', 'SPYC949H2130100', '  ', '  ', 3, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236002', 12, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000490    ', '              ', 0, '          ', '          ', 'SPYC949H2130105', '  ', '  ', 3, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236002', 13, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000500    ', '              ', 0, '          ', '          ', 'SPYC949H2159095', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236002', 14, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000510    ', '              ', 0, '          ', '          ', 'SPYC949H2159100', '  ', '  ', 3, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236002', 6, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000360    ', '              ', 0, '          ', '          ', 'SPWJ948S2355070', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236002', 7, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000380    ', '              ', 0, '          ', '          ', 'SPWJ948S2355080', '  ', '  ', 2, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236002', 8, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000390    ', '              ', 0, '          ', '          ', 'SPWJ948S2356070', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1907236000', 1, '01', 'SA', 'CFGY', '20190723', 6000, '8074783296', '000260    ', '              ', 0, '          ', '          ', 'SPWH936D5430075', '  ', '  ', 8, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1907236000', 2, '01', 'SA', 'CFGY', '20190723', 6000, '8074783296', '000270    ', '              ', 0, '          ', '          ', 'SPWH936D5430080', '  ', '  ', 8, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1907236010', 1, '01', 'SA', 'CFGY', '20190723', 6010, '8074783305', '000610    ', '              ', 0, '          ', '          ', 'SPYC949S1139085', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1907236010', 2, '01', 'SA', 'CFGY', '20190723', 6010, '8074783305', '000620    ', '              ', 0, '          ', '          ', 'SPYC949S1139090', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1907236010', 3, '01', 'SA', 'CFGY', '20190723', 6010, '8074783305', '000630    ', '              ', 0, '          ', '          ', 'SPYC949S1139095', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1907236010', 4, '01', 'SA', 'CFGY', '20190723', 6010, '8074783305', '000640    ', '              ', 0, '          ', '          ', 'SPYC949S1159085', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1907236010', 5, '01', 'SA', 'CFGY', '20190723', 6010, '8074783305', '000650    ', '              ', 0, '          ', '          ', 'SPYC949S1159090', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1907236010', 6, '01', 'SA', 'CFGY', '20190723', 6010, '8074783305', '000660    ', '              ', 0, '          ', '          ', 'SPYC949S1159095', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1907236010', 7, '01', 'SA', 'CFGY', '20190723', 6010, '8074783305', '000850    ', '              ', 0, '          ', '          ', 'SPYS949H2250095', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1907236010', 8, '01', 'SA', 'CFGY', '20190723', 6010, '8074783305', '000860    ', '              ', 0, '          ', '          ', 'SPYS949H2250100', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1907236010', 9, '01', 'SA', 'CFGY', '20190723', 6010, '8074783305', '000870    ', '              ', 0, '          ', '          ', 'SPYS949H2250105', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CJC11907236000', 1, '01', 'Q3', 'CJC1', '20190723', 6002, '8074783302', '000310    ', '              ', 0, '          ', '          ', 'Q3AFAFDU6S2100230', '  ', '  ', 1, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CJC11907236000', 1, '01', 'Q3', 'CJC1', '20190723', 6002, '8074783302', '000310    ', '              ', 0, '          ', '          ', 'Q3AFAFDU6S2100240', '  ', '  ', 2, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CJC11907236000', 1, '01', 'Q3', 'CJC1', '20190723', 6002, '8074783302', '000310    ', '              ', 0, '          ', '          ', 'Q3AFAFDU6S2100250', '  ', '  ', 3, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CJC11907236000', 1, '01', 'Q3', 'CJC1', '20190723', 6002, '8074783302', '000310    ', '              ', 0, '          ', '          ', 'Q3AFAFDU6S2100260', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CJC11907236000', 1, '01', 'Q3', 'CJC1', '20190723', 6002, '8074783302', '000310    ', '              ', 0, '          ', '          ', 'Q3AFAFDU6S2100270', '  ', '  ', 5, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236003', 1, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000310    ', '              ', 0, '          ', '          ', 'SPWJ948S2255070', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CEGP1907236003', 2, '01', 'SA', 'CEGP', '20190723', 6002, '8074783302', '000320    ', '              ', 0, '          ', '          ', 'SPWJ948S2255075', '  ', '  ', 3, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1908196000', 1, '01', 'SA', 'CFGY', '20190819', 6000, '8074783296', '000260    ', '              ', 0, '          ', '          ', 'SPWH936D5430075', '  ', '  ', 8, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1908196000', 2, '01', 'SA', 'CFGY', '20190819', 6000, '8074783296', '000270    ', '              ', 0, '          ', '          ', 'SPWH936D5430080', '  ', '  ', 8, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1908196010', 1, '01', 'SA', 'CFGY', '20190819', 6010, '8074783305', '000610    ', '              ', 0, '          ', '          ', 'SPYC949S1139085', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1908196010', 2, '01', 'SA', 'CFGY', '20190819', 6010, '8074783305', '000620    ', '              ', 0, '          ', '          ', 'SPYC949S1139090', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1908196010', 3, '01', 'SA', 'CFGY', '20190819', 6010, '8074783305', '000630    ', '              ', 0, '          ', '          ', 'SPYC949S1139095', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1908196010', 4, '01', 'SA', 'CFGY', '20190819', 6010, '8074783305', '000640    ', '              ', 0, '          ', '          ', 'SPYC949S1159085', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1908196010', 5, '01', 'SA', 'CFGY', '20190819', 6010, '8074783305', '000650    ', '              ', 0, '          ', '          ', 'SPYC949S1159090', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1908196010', 6, '01', 'SA', 'CFGY', '20190819', 6010, '8074783305', '000660    ', '              ', 0, '          ', '          ', 'SPYC949S1159095', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1908196010', 7, '01', 'SA', 'CFGY', '20190819', 6010, '8074783305', '000850    ', '              ', 0, '          ', '          ', 'SPYS949H2250095', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1908196010', 8, '01', 'SA', 'CFGY', '20190819', 6010, '8074783305', '000860    ', '              ', 0, '          ', '          ', 'SPYS949H2250100', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1908196010', 9, '01', 'SA', 'CFGY', '20190819', 6010, '8074783305', '000870    ', '              ', 0, '          ', '          ', 'SPYS949H2250105', '  ', '  ', 4, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1908206000', 1, '01', 'SA', 'CFGY', '20190820', 6000, '8074783296', '000260    ', '              ', 0, '          ', '          ', 'SPWJ948S2255075', '  ', '  ', 1, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CFGY1908206010', 1, '01', 'SA', 'CFGY', '20190820', 6000, '8074783296', '000260    ', '              ', 0, '          ', '          ', 'SPWJ948S2255075', '  ', '  ', 1, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CJC11909026000', 1, '01', 'Q3', 'CJC1', '20190902', 6000, '8074783296', '000260    ', '              ', 0, '          ', '          ', 'SPWJ948S2255075', '  ', '  ', 1, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', ''),
		('CJUT1909026000', 1, '01', 'SA', 'CJUT', '20190902', 6000, '8074783296', '000260    ', '              ', 0, '          ', '          ', 'SPWJ948S2255075', '  ', '  ', 1, 0, 0.00, '  ', 0, 'SYS_BG_USER', null, 'SYS_BG_USER', null, '', 'S', null, '  ', '', '', '');
	`

	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createRecvSuppDtlTable error: %v", err.Error())
		fmt.Println()
	}
}

// 误差部分
func createStockMisDtlTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createStockMisDtlTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.StockMisDtl;"); err != nil {
		log.Printf("createStockMisDtlTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE CSL.dbo.StockMisDtl
		(
			BrandCode VARCHAR(4) NOT NULL,
			ShopCode CHAR(4) NOT NULL,
			Dates CHAR(8) NOT NULL,
			SeqNo INT NOT NULL,
			DtSeq INT NOT NULL,
			SAPDeliveryNo CHAR(10),
			SAPDeliveryItemNo CHAR(10),
			ProdCode VARCHAR(18) NOT NULL,
			RecvSuppQty INT,
			StockMisQty INT,
			StockMisFixQty INT,
			StockMisReasonCode CHAR(2),
			StockMisStatusCode CHAR(2),
			StockMisResultCode CHAR(2),
			RecvDate CHAR(8),
			RecvSuppNo CHAR(14),
			StyleCode VARCHAR(18) NOT NULL,
			StockMisNo CHAR(14) NOT NULL,
			RecvSuppDate CHAR(8),
			ShopDesc NVARCHAR(200),
			BrandDesc NVARCHAR(200),
			PlantCode CHAR(4) NOT NULL,
			InsertEmpID CHAR(10),
			InsertEmpName NVARCHAR(100) NOT NULL,
			RecvEmpID CHAR(10),
			RecvEmpName NVARCHAR(100),
			StockMisRecvDateTime SMALLDATETIME,
			ProcessEmpID VARCHAR(20),
			ProcessEmpName NVARCHAR(100),
			StockMisProcDateTime SMALLDATETIME,
			NewSAPDeliveryNo CHAR(10),
			BoxNo CHAR(20) NOT NULL,
			RecvSuppType CHAR,
			DelChk BIT DEFAULT 0 NOT NULL,
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME,
			ModiUserID VARCHAR(20) NOT NULL,
			ModiDateTime DATETIME,
			SendState VARCHAR(2) DEFAULT '' NOT NULL,
			SendFlag CHAR DEFAULT 'R' NOT NULL,
			SendSeqNo BIGINT NOT NULL IDENTITY,
			SendDateTime DATETIME,
			InvtBaseDate CHAR(8) NOT NULL,
			RecvSuppNoNew CHAR(14),
			SendEmpName NVARCHAR(100),
			WayBillNo01 CHAR(13)
		);
	`

	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createRecvSuppMstTable error: %v", err.Error())
		fmt.Println()
	}
}

// 商品部分，误差需要
func createProductTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createProductTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.Product;"); err != nil {
		log.Printf("createProductTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE CSL.dbo.Product
		(
			BrandCode VARCHAR(4) NOT NULL,
			ProdCode VARCHAR(18) NOT NULL,
			ProdName NVARCHAR(200),
			StyleCode VARCHAR(18) NOT NULL,
			ColorCode CHAR(3) NOT NULL,
			ColorName NVARCHAR(200),
			SizeCode CHAR(3) NOT NULL,
			SizeName NVARCHAR(200),
			BaseUnitCnt INT,
			BaseUnitCode CHAR(3),
			BaseUnitName NVARCHAR(60),
			UseChk BIT DEFAULT 1 NOT NULL,
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME,
			Price DECIMAL(13,2) DEFAULT 0 NOT NULL,
			PreProdCode VARCHAR(18),
			ModiUserID VARCHAR(20),
			ModiDateTime DATETIME
		);
	`

	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createProductTable error: %v", err.Error())
		fmt.Println()
	}

	sql = `
		INSERT INTO CSL.dbo.Product (BrandCode, ProdCode, ProdName, StyleCode, ColorCode, ColorName, SizeCode, SizeName, BaseUnitCnt, BaseUnitCode, BaseUnitName, UseChk, InUserID, InDateTime, Price, PreProdCode, ModiUserID, ModiDateTime) 
		VALUES 
		('SA', 'SPWH936D5430075', N'蜡笔小新彩色短裙, (30)Yellow, 165/70A(M)', 'SPWH936D54', '30 ', '(30)Yellow', '075', 'M', null, null, null, 1, 'system', '2019-04-22 00:00:00.000', 199.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPWH936D5430080', N'蜡笔小新彩色短裙, (30)Yellow, 170/74A(L)', 'SPWH936D54', '30 ', '(30)Yellow', '080', 'L', null, null, null, 1, 'system', '2019-04-22 00:00:00.000', 199.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPWJ948S2255070', N'女式牛仔裙, (55)Indigo, 160/66A(S)', 'SPWJ948S22', '55 ', '(55)Indigo', '070', 'S', null, null, null, 1, 'system', '2019-05-16 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPWJ948S2255075', N'女式牛仔裙, (55)Indigo, 165/70A(M)', 'SPWJ948S22', '55 ', '(55)Indigo', '075', 'M', null, null, null, 1, 'system', '2019-05-16 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPWJ948S2255080', N'女式牛仔裙, (55)Indigo, 170/74A(L)', 'SPWJ948S22', '55 ', '(55)Indigo', '080', 'L', null, null, null, 1, 'system', '2019-05-16 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPWJ948S2256070', N'女式牛仔裙, (56)Light Indigo, 160/66A(S)', 'SPWJ948S22', '56 ', '(56)Light Indigo', '070', 'S', null, null, null, 1, 'system', '2019-05-16 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPWJ948S2256075', N'女式牛仔裙, (56)Light Indigo, 165/70A(M)', 'SPWJ948S22', '56 ', '(56)Light Indigo', '075', 'M', null, null, null, 1, 'system', '2019-05-16 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPWJ948S2355070', N'女式牛仔裙, (55)Indigo, 160/66A(S)', 'SPWJ948S23', '55 ', '(55)Indigo', '070', 'S', null, null, null, 1, 'system', '2019-05-16 00:00:00.000', 199.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPWJ948S2355080', N'女式牛仔裙, (55)Indigo, 170/74A(L)', 'SPWJ948S23', '55 ', '(55)Indigo', '080', 'L', null, null, null, 1, 'system', '2019-05-16 00:00:00.000', 199.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPWJ948S2356070', N'女式牛仔裙, (56)Light Indigo, 160/66A(S)', 'SPWJ948S23', '56 ', '(56)Light Indigo', '070', 'S', null, null, null, 1, 'system', '2019-05-16 00:00:00.000', 199.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPWJ948S2356075', N'女式牛仔裙, (56)Light Indigo, 165/70A(M)', 'SPWJ948S23', '56 ', '(56)Light Indigo', '075', 'M', null, null, null, 1, 'system', '2019-05-16 00:00:00.000', 199.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYC949H2130095', N'大格纹衬衫, (30)Yellow, 170/92A(M)', 'SPYC949H21', '30 ', '(30)Yellow', '095', 'M', null, null, null, 1, 'system', '2019-05-17 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYC949H2130100', N'大格纹衬衫, (30)Yellow, 175/96A(L)', 'SPYC949H21', '30 ', '(30)Yellow', '100', 'L', null, null, null, 1, 'system', '2019-05-17 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYC949H2130105', N'大格纹衬衫, (30)Yellow, 180/100A(XL)', 'SPYC949H21', '30 ', '(30)Yellow', '105', 'XL', null, null, null, 1, 'system', '2019-05-17 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYC949H2159095', N'大格纹衬衫, (59)Navy, 170/92A(M)', 'SPYC949H21', '59 ', '(59)Navy', '095', 'M', null, null, null, 1, 'system', '2019-06-24 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYC949H2159100', N'大格纹衬衫, (59)Navy, 175/96A(L)', 'SPYC949H21', '59 ', '(59)Navy', '100', 'L', null, null, null, 1, 'system', '2019-06-24 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYC949S1139085', N'女式格子衬衫, (39)Ivory, 160/84A(S)', 'SPYC949S11', '39 ', '(39)Ivory', '085', 'S', null, null, null, 1, 'system', '2019-05-15 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYC949S1139090', N'女式格子衬衫, (39)Ivory, 165/88A(M)', 'SPYC949S11', '39 ', '(39)Ivory', '090', 'M', null, null, null, 1, 'system', '2019-05-15 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYC949S1139095', N'女式格子衬衫, (39)Ivory, 170/92A(L)', 'SPYC949S11', '39 ', '(39)Ivory', '095', 'L', null, null, null, 1, 'system', '2019-05-15 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYC949S1159085', N'女式格子衬衫, (59)Navy, 160/84A(S)', 'SPYC949S11', '59 ', '(59)Navy', '085', 'S', null, null, null, 1, 'system', '2019-05-15 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYC949S1159090', N'女式格子衬衫, (59)Navy, 165/88A(M)', 'SPYC949S11', '59 ', '(59)Navy', '090', 'M', null, null, null, 1, 'system', '2019-05-15 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYC949S1159095', N'女式格子衬衫, (59)Navy, 170/92A(L)', 'SPYC949S11', '59 ', '(59)Navy', '095', 'L', null, null, null, 1, 'system', '2019-05-15 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYS949H2250095', N'条纹衬衫, (50)Blue, 170/92A(M)', 'SPYS949H22', '50 ', '(50)Blue', '095', 'M', null, null, null, 1, 'system', '2019-05-17 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYS949H2250100', N'条纹衬衫, (50)Blue, 175/96A(L)', 'SPYS949H22', '50 ', '(50)Blue', '100', 'L', null, null, null, 1, 'system', '2019-05-17 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('SA', 'SPYS949H2250105', N'条纹衬衫, (50)Blue, 180/100A(XL)', 'SPYS949H22', '50 ', '(50)Blue', '105', 'XL', null, null, null, 1, 'system', '2019-05-17 00:00:00.000', 259.00, ' ', null, '2019-08-06 05:29:20.887'),
		('Q3', 'Q3AFAFDU6S2100230', N'男女同款拖鞋, (00)生产代表颜色, 230MM', 'Q3AFAFDU6S21', '00 ', N'(00)生产代表颜色', '230', '230MM', null, null, null, 1, 'system', '2018-02-12 05:07:58.857', 79.00, ' ', null, '2019-08-09 05:27:44.677'),
		('Q3', 'Q3AFAFDU6S2100240', N'男女同款拖鞋, (00)生产代表颜色, 240MM', 'Q3AFAFDU6S21', '00 ', N'(00)生产代表颜色', '240', '240MM', null, null, null, 1, 'system', '2018-02-12 05:07:58.857', 79.00, ' ', null, '2019-08-09 05:27:44.677'),
		('Q3', 'Q3AFAFDU6S2100250', N'男女同款拖鞋, (00)生产代表颜色, 250MM', 'Q3AFAFDU6S21', '00 ', N'(00)生产代表颜色', '250', '250MM', null, null, null, 1, 'system', '2018-02-12 05:07:58.857', 79.00, ' ', null, '2019-08-09 05:27:44.677'),
		('Q3', 'Q3AFAFDU6S2100260', N'男女同款拖鞋, (00)生产代表颜色, 260MM', 'Q3AFAFDU6S21', '00 ', N'(00)生产代表颜色', '260', '260MM', null, null, null, 1, 'system', '2018-02-12 05:07:58.857', 79.00, ' ', null, '2019-08-09 05:27:44.677'),
		('Q3', 'Q3AFAFDU6S2100270', N'男女同款拖鞋, (00)生产代表颜色, 270MM', 'Q3AFAFDU6S21', '00 ', N'(00)生产代表颜色', '270', '270MM', null, null, null, 1, 'system', '2018-02-12 05:07:58.857', 79.00, ' ', null, '2019-08-09 05:27:44.677'),
		('Q3', 'Q3AFAFDU6S2100280', N'男女同款拖鞋, (00)生产代表颜色, 280MM', 'Q3AFAFDU6S21', '00 ', N'(00)生产代表颜色', '280', '280MM', null, null, null, 1, 'system', '2018-02-12 05:07:58.857', 79.00, ' ', null, '2019-08-09 05:27:44.677');
	`

	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createProductTable error: %v", err.Error())
		fmt.Println()
	}
}

// SP 部分
func createSP() {
	createMonthlyClosingChk()
	createComonRaiseError()
	createChkRecvSupp()
	createUpdateStockInEnterConfirmSaveRecvSuppMstR1ClearanceByWaybillNo()
	createInsertStockInMissSaveStockMisDtlC1()
}

func createMonthlyClosingChk() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createMonthlyClosingChk error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP FUNCTION IF EXISTS dbo.udf_CSLK_MonthlyClosingChk"); err != nil {
		log.Printf("createMonthlyClosingChk error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE FUNCTION [dbo].[udf_CSLK_MonthlyClosingChk] (@Id varchar(50), @language char(2))  
		RETURNS nvarchar(max)
		AS  
		BEGIN  
			DECLARE @return nvarchar(max)  
			set @return ='0'  
				
			if (@language is null or @language='')  
			set @language='Zn'  
				
			if exists (select *  
			from MonthlyBizFuctionClosing  
			where BizFuctionID =@Id  
			and  GetDate() between StartDateTime and EndDatetime  
			and ClosingChk=1)  
			
			select @return='1'+case when @language ='Kr' then InformKr else InformZn end  
			from MonthlyBizFuctionClosing  
			where BizFuctionID=@Id  
			and  GetDate() between StartDateTime and EndDatetime  
			and ClosingChk=1  
			RETURN @return  
		END;
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createMonthlyClosingChk error: %v", err.Error())
		log.Println()
	}
}

func createComonRaiseError() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createComonRaiseError error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP PROCEDURE IF EXISTS dbo.up_CSLK_ComonRaiseError"); err != nil {
		log.Printf("createComonRaiseError error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE PROCEDURE [dbo].[up_CSLK_ComonRaiseError]  
			@ErrorCode NVARCHAR(50),  
			@ErrorParam1 NVARCHAR(4000),  
			@ErrorParam2 NVARCHAR(4000)  
			
		AS  
		BEGIN  
			DECLARE @ErrorMessage NVARCHAR(4000);  
			DECLARE @ErrorSeverity INT;  
			DECLARE @ErrorState INT;  
			
			SET @ErrorMessage = 'C=%s|P=%s|P=%s|S='+ '[ERROR Proc]'+ERROR_PROCEDURE()  +',[ERROR_LINE]'+convert(varchar,ERROR_LINE()) +',[ERROR_MESSAGE]'+ERROR_MESSAGE()    
			SET @ErrorSeverity = ERROR_SEVERITY()    
			SET @ErrorState = ERROR_STATE()    
			IF @ErrorSeverity is null 
			BEGIN
				SET @ErrorSeverity = 16 ; 
				SET @ErrorState = 1;
			END 
			
			
			IF(@ErrorCode ='' AND ERROR_NUMBER()<>50000)  
			SELECT @ErrorCode= 'EDB_' + convert(varchar,ERROR_NUMBER())    
			IF(@ErrorCode ='' AND ERROR_NUMBER()=50000 AND CHARINDEX ( CHAR(12),ERROR_MESSAGE())>0)  
			SELECT @ErrorCode= 'EDB_' +   LEFT( ERROR_MESSAGE(), CHARINDEX ( CHAR(12),ERROR_MESSAGE()) -1)  
		
			RAISERROR (@ErrorMessage,    
			@ErrorSeverity,     
			@ErrorState,     
			@ErrorCode,    
			@ErrorParam1,    
			@ErrorParam2);    
		END;
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createComonRaiseError error: %v", err.Error())
		log.Println()
	}
}

func createChkRecvSupp() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createChkRecvSupp error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP PROCEDURE IF EXISTS dbo.up_CSLK_IF_CHK_RecvSupp"); err != nil {
		log.Printf("createChkRecvSupp error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE PROCEDURE [dbo].[up_CSLK_IF_CHK_RecvSupp]
		(@p_RECVSUPPNO varchar(14)='', @o_SENDF char(1) OUTPUT)
		WITH 
		EXECUTE AS CALLER
		AS
		BEGIN  
			SET  XACT_ABORT ON;  
			SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED  
		
			DECLARE @Sql   nvarchar (500)  
		
			DECLARE @ParmDefinition   nvarchar (500);	
			DECLARE @ChkFlag BIT
		
			SET @ChkFlag  = (SELECT IFSendFlagChk FROM IFConfig)
		
			SET @o_SENDF = 'S'
		
			IF (@ChkFlag <> 1)
			BEGIN
				SET @o_SENDF = 'S'
				RETURN;
			END
		END;
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createChkRecvSupp error: %v", err.Error())
		log.Println()
	}
}

func createUpdateStockInEnterConfirmSaveRecvSuppMstR1ClearanceByWaybillNo() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createUpdateStockInEnterConfirmSaveRecvSuppMstR1ClearanceByWaybillNo error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP PROCEDURE IF EXISTS dbo.up_CSLK_IOM_UpdateStockInEnterConfirmSave_RecvSuppMst_R1_Clearance_By_WaybillNo"); err != nil {
		log.Printf("createUpdateStockInEnterConfirmSaveRecvSuppMstR1ClearanceByWaybillNo error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE PROCEDURE [dbo].[up_CSLK_IOM_UpdateStockInEnterConfirmSave_RecvSuppMst_R1_Clearance_By_WaybillNo]
			@BrandCode VARCHAR(4),   -- 브랜드 코드
			@ShopCode CHAR(4),       -- 主卖场
			@WaybillNo VARCHAR(13),  -- 운송장번호
			@InDate CHAR(8),			-- 入库时间
			@EmpID CHAR(10)      -- 入库人EmpNo
			AS
			--SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED
			SET XACT_ABORT ON;
			SET NOCOUNT ON;
			BEGIN

				DECLARE @RecvSuppNo CHAR(14);
				DECLARE @SENDF CHAR(1);
				DECLARE @SendFlag CHAR(1);
				DECLARE @TempRecvSuppStatusCode CHAR(1);
				DECLARE @ErrorCode NVARCHAR(1000) = '';
				DECLARE @ErrorParam1 NVARCHAR(4000) = '';
				DECLARE @ErrorParam2 NVARCHAR(4000) = '';
				DECLARE @RecvEmpName NVARCHAR(200);
				DECLARE @UserID VARCHAR(20);

				BEGIN TRY
					-- 마감체크
					IF (
						LEFT(dbo.udf_CSLK_MonthlyClosingChk('01', 'Zn'), 1) = 1
					)
					BEGIN
						SELECT @ErrorCode
							= SUBSTRING(
								dbo.udf_CSLK_MonthlyClosingChk('01', 'Zn'), 2, 510
									);
						EXEC [up_CSLK_ComonRaiseError] @ErrorCode, @ErrorParam1, @ErrorParam2;
					END;


					IF(@InDate = '' OR @InDate IS NULL)
					BEGIN
						SET @InDate = CONVERT(CHAR(8),GETDATE(),112)
					END

					--Query Start
					PRINT ('@RecvSuppNo1');

					SELECT TOP 1
						@RecvSuppNo = RecvSuppNo
					FROM RecvSuppMst
					WHERE BrandCode = @BrandCode
						AND ShopCode = @ShopCode
						AND WayBillNo = @WaybillNo
						AND ShippingTypeCode IN ( '01', '66', '16' )
						AND DelChk = '0';

					PRINT ('@RecvSuppNo2');
					PRINT (@RecvSuppNo);


			SELECT @RecvEmpName = EmpName
			FROM Employee
			WHERE EmpID = @EmpID


			SELECT @UserID = UserID
			FROM UserInfo
			WHERE EmpID = @EmpID

					-- 인터페이스 실시간 처리 확인
					EXEC [up_CSLK_IF_CHK_RecvSupp] @p_RECVSUPPNO = @RecvSuppNo,
						@o_SENDF = @SENDF OUTPUT;

					/* 수신 시스템 체크
					R : 등록/수정된 상태, 전송 전 ( 수정 불가능 )
					I : 전송 중 ( 수정불가능 )
					S : 전송 후 ( 수정가능 )
					*/
					IF (
						@SENDF = 'S'
					)
					BEGIN

						-- 이미 입고확인 되었을경우 입고불가

						SELECT @TempRecvSuppStatusCode = RecvSuppStatusCode
						FROM RecvSuppMst
						WHERE RecvSuppNo = @RecvSuppNo;

						IF @TempRecvSuppStatusCode != 'F'
						BEGIN

							-- SAP진행상황 체크
							SELECT @SendFlag = SendFlag
							FROM RecvSuppMst
							WHERE RecvSuppNo = @RecvSuppNo;

							-- R:등록  I:전송중  S:완료
							IF (
								@SendFlag = 'I'
							)
							BEGIN
								SET @ErrorCode = 'COM100';
								SET @ErrorParam1 = '';
								SET @ErrorParam2 = '';
								EXEC [up_CSLK_ComonRaiseError] @ErrorCode, @ErrorParam1, @ErrorParam2;
							END;

							SELECT TOP 1
								@SendFlag = SendFlag
							FROM RecvSuppDtl
							WHERE RecvSuppNo = @RecvSuppNo;

							-- R:등록  I:전송중  S:완료
							IF (
								@SendFlag = 'I'
							)
							BEGIN
								SET @ErrorCode = 'COM100';
								SET @ErrorParam1 = '';
								SET @ErrorParam2 = '';
								EXEC [up_CSLK_ComonRaiseError] @ErrorCode, @ErrorParam1, @ErrorParam2;
							END;
							ELSE
							BEGIN
								--위는 다 SAP에서 하는 체크.
								UPDATE RecvSuppMst
								SET RecvChk = 1,                  -- 매장입고여부  조회할때 RecvChk 만 보고  RecvSuppStatusCode = 'F'	는 안본다. RecvChk = 0 ,RecvSuppStatusCode = 'F' 면 데이터 꼬인것.
									ShopSuppRecvDate = @InDate, -- 매장입고일자
									RecvSuppStatusCode = 'F',     -- 입고확정
									RecvEmpID = @EmpID,           -- 입고자사번
									RecvEmpName = @RecvEmpName,   -- 입고자명
									ModiUserID = @UserID,
									ModiDateTime = GETDATE(),
									SendFlag = 'R',
									InvtBaseDate = @InDate,     -- 재고기준일자.
									Channel = 'Clearance'
								WHERE BrandCode = @BrandCode
									AND ShopCode = @ShopCode
									AND WayBillNo = @WaybillNo
									AND BoxNo = @WaybillNo
									AND ShippingTypeCode IN ( '01', '66', '16' ); --wangpengda 20100203

								UPDATE RecvSuppDtl
								SET RecvSuppFixedQty = RecvSuppQty, -- 출고수량과 입고수량을 동일하게 맞춤
									ModiUserID = @UserID,
									ModiDateTime = GETDATE(),
									SendFlag = 'R'
								FROM
								(
									SELECT Mst.RecvSuppNo,
										Dtl.RecvSuppSeqNo
									FROM RecvSuppMst AS Mst
										INNER JOIN RecvSuppDtl AS Dtl
											ON (Mst.RecvSuppNo = Dtl.RecvSuppNo)
									WHERE Mst.BrandCode = @BrandCode
										AND Mst.ShopCode = @ShopCode
										AND WayBillNo = @WaybillNo
										AND BoxNo = @WaybillNo
										AND ShippingTypeCode IN ( '01', '66', '16' )
										AND Mst.RecvChk = 1
								) AS A
								WHERE RecvSuppDtl.RecvSuppNo = A.RecvSuppNo
									AND RecvSuppDtl.RecvSuppSeqNo = A.RecvSuppSeqNo;
							END;
						END;
						SELECT * FROM dbo.RecvSuppMst
							WHERE BrandCode = @BrandCode
							AND ShopCode = @ShopCode
							AND WayBillNo = @WaybillNo
							AND BoxNo = @WaybillNo
							AND ShippingTypeCode IN ( '01', '66', '16' );
					END;
					ELSE
					BEGIN

						-- 본사 수정중입니다.
						SET @ErrorCode = 'IOM132';
						SET @ErrorParam1 = '';
						SET @ErrorParam2 = '';
						EXEC [up_CSLK_ComonRaiseError] @ErrorCode, @ErrorParam1, @ErrorParam2;
					END;
				--Query End
				END TRY
				BEGIN CATCH
					EXEC [up_CSLK_ComonRaiseError] @ErrorCode, @ErrorParam1, @ErrorParam2;
				END CATCH;
			END;
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createUpdateStockInEnterConfirmSaveRecvSuppMstR1ClearanceByWaybillNo error: %v", err.Error())
		log.Println()
	}
}

func createInsertStockInMissSaveStockMisDtlC1() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createInsertStockInMissSaveStockMisDtlC1 error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP PROCEDURE IF EXISTS dbo.up_CSLK_IOM_InsertStockInMissSave_StockMisDtl_C1_Clearance"); err != nil {
		log.Printf("createInsertStockInMissSaveStockMisDtlC1 error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE PROCEDURE [dbo].[up_CSLK_IOM_InsertStockInMissSave_StockMisDtl_C1_Clearance]
		@BrandCode  VARCHAR(4)
		,@ShopCode  CHAR(4)
		,@WaybillNo  VARCHAR(13)
		,@ProdCode  VARCHAR(18)
		,@ShopRecvSuppQty INT
		,@ShopInFixQty  INT
		,@ErrorRegEmpID VARCHAR(20) -- 登记这条误差的人的EmpID
		,@RecvDate	CHAR(8)='' -- 入库日期
		AS
		SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED
		SET XACT_ABORT ON
		SET NOCOUNT ON

		BEGIN

		DECLARE @ErrorCode  NVARCHAR(50) ='';
		DECLARE @ErrorParam1 NVARCHAR(4000)='';
		DECLARE @ErrorParam2 NVARCHAR(4000)='';


		DECLARE @CurrDate    CHAR(8)   -- 현재일자
		DECLARE @SeqNo     INT    -- 착오순번
		DECLARE @DtSeq     INT    -- 착오상세순번
		DECLARE @ErrorRegEmpName  NVARCHAR(200) -- 등록자명
		DECLARE @ErrorRegUserID  VARCHAR(20)
		DECLARE @SuppEmpName NVARCHAR(100)

		DECLARE @RecvSuppNo    CHAR(14)  -- 대표 입출고번호
		DECLARE @RecvSuppNo_Temp  CHAR(14)  -- 상품이 포함된 입출고번호

		BEGIN TRY

		-- 현재년월일 입력
		SELECT @CurrDate = CONVERT(CHAR(8),GETDATE(),112)

		SELECT @ErrorRegEmpName = E.EmpName
			, @ErrorRegUserID = U.UserID
		FROM CSL.dbo.Employee E
		JOIN CSL.dbo.UserInfo U
		ON E.EmpID = U.EmpID
		WHERE E.EmpID = @ErrorRegEmpID
		AND U.UseChk = 1

		/***착오등록막음(복합매장)  2011.12.22  ***/
		DECLARE @StyleCode VARCHAR(18)

		DECLARE  @Exists table ( e bit)

		select @Stylecode = StyleCode
		from Product where ProdCode = @ProdCode

		-- 이미 해당 상품에 대한 오차가 등록되었을 경우
		IF EXISTS ( SELECT * FROM StockMisDtl WHERE BrandCode = @BrandCode AND ShopCode = @ShopCode AND BoxNo = @WaybillNo AND ProdCode = @ProdCode  AND DelChk = 0 )
		BEGIN
		-- 에러 처리
		SET @ErrorCode = 'IOM138'
		EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2

		--PRINT N'이미 해상 상품에 대한 오차가 등록되었습니다.'

		END


		IF ( (@ShopInFixQty - @ShopRecvSuppQty ) >= 0 )
		BEGIN
		/******** 실물이 많을 경우 *********/
		--PRINT N'실물이 많을 경우'

		DECLARE @RecvSuppQty_ProdCode INT    -- 상품별 출고수량 ( 총합 아님 )
		DECLARE @SAPDeliveryNo   CHAR(10)  -- 물류팀에서 필요로 하는 값(RecvSuppDtl에서 가져옴, 신규착오는 NULL)
		DECLARE @SAPDeliveryNo_Temp  CHAR(10)
		DECLARE @SAPDeliveryItemNo  CHAR(10)  -- 물류팀에서 필요로 하는 값(RecvSuppDtl에서 가져옴, 신규착오는 NULL)


		-- 해당 박스의 첫번째 입출고 번호 구함
		SELECT TOP 1
		@RecvSuppNo = MST.RecvSuppNo,
		@SeqNo = MST.SeqNo,
		@SAPDeliveryNo = MST.SAPDeliveryNo,
		@SuppEmpName = MST.SuppEmpName
		FROM RecvSuppMst AS MST
		WHERE MST.BrandCode = @BrandCode AND MST.ShopCode = @ShopCode AND MST.WayBillNo = @WaybillNo
		AND MST.BoxNo = @WaybillNo
		AND ShippingTypeCode IN ('01','66','16') --wangpengda 20100203
		AND MST.DelChk = 0
		AND MST.ShopSuppRecvDate=@RecvDate			--ADD BY SHEN.XUE 20140804


		-- 해당상품이 속한 첫번째 입출고 번호 구함
		SELECT TOP 1 @RecvSuppNo_Temp = MST.RecvSuppNo, @SeqNo = MST.SeqNo, @RecvSuppQty_ProdCode = DTL.RecvSuppQty,
			@SAPDeliveryNo_Temp = DTL.SAPDeliveryNo, @SAPDeliveryItemNo = DTL.SAPDeliveryItemNo
		FROM RecvSuppMst AS MST
		INNER JOIN RecvSuppDtl AS DTL
		ON (MST.RecvSuppNo = DTL.RecvSuppNo AND DTL.DelChk = 0)
		WHERE MST.BrandCode = @BrandCode AND MST.ShopCode = @ShopCode AND MST.WayBillNo = @WaybillNo AND MST.BoxNo = @WaybillNo
		AND DTL.ProdCode = @ProdCode
		AND ShippingTypeCode IN ('01','66','16') --wangpengda 20100203
		AND MST.DelChk = 0
		AND MST.ShopSuppRecvDate=@RecvDate			--ADD BY SHEN.XUE 20140804
		IF @RecvSuppNo_Temp IS NOT NULL
		BEGIN
		SET @RecvSuppNo = @RecvSuppNo_Temp
		SET @SAPDeliveryNo = @SAPDeliveryNo_Temp
		END

		-- DtSeq 일렬번호 구함
		SELECT @DtSeq = ISNULL(MAX(DTL.DtSeq),0) + 1 FROM StockMisDtl AS DTL
		WHERE DTL.BrandCode = @BrandCode
		AND DTL.ShopCode = @ShopCode
		AND DTL.SeqNo = @SeqNo --modify by li.guolin 20151013
		AND DTL.Dates=@CurrDate  --modify by li.guolin 20151013
		print @DtSeq

		PRINT @RecvSuppNo
		-- 오차 데이터 생성
		INSERT StockMisDtl
			( BrandCode
			, ShopCode
			, Dates
			, SeqNo
			, DtSeq
			, SAPDeliveryNo
			, SAPDeliveryItemNo
			, ProdCode
			, RecvSuppQty
			, StockMisQty
			, StockMisFixQty
			, StockMisReasonCode
			, StockMisStatusCode
			, StockMisResultCode
			, RecvDate
			, RecvSuppNo
			, StyleCode
			, StockMisNo
			, RecvSuppDate
			, ShopDesc
			, BrandDesc
			, PlantCode
			, InsertEmpID
			, RecvEmpID
			, StockMisRecvDateTime
			, ProcessEmpID
			, StockMisProcDateTime
			, NewSAPDeliveryNo
			, RecvSuppType
			, BoxNo
			, DelChk
			, InUserID
			, InDateTime
			, ModiUserID
			, ModiDateTime
			, RecvEmpName
			, ProcessEmpName
			, InsertEmpName
			, SendState
			, SendFlag
			, SendDateTime
			, InvtBaseDate
			, SendEmpName
			, WayBillNo01                              -- 운송장번호 김재훈 20110503
			)
		SELECT TOP 1 B.BrandCode
			, B.ShopCode        -- 매장코드
			, @CurrDate        -- 현재일자 Dates
			, A.SeqNo         -- SeqNo
			, @DtSeq         -- DtSeq -- 상세순번
			, @SAPDeliveryNo       -- SAPDeliveryNo
			, @SAPDeliveryItemNo      -- SAPDeliveryItemNo
			, @ProdCode        --, A.ProdCode        -- ProdCode
			, ISNULL(@RecvSuppQty_ProdCode,0)       -- RecvSuppQty 출고수량
			, ABS(@ShopRecvSuppQty - @ShopInFixQty) -- StockMisQty 착오수량
			, 0          -- StockMisFixQty 착오확정수량
			, Null          -- StockMisReasonCode 착오
			, '10'          -- StockMisStatusCode 착오 처리 상태(10: 매장입고)
			, Null          -- StockMisResultCode
			, B.ShopSuppRecvDate      -- RecvDate
			, @RecvSuppNo        -- RecvSuppNo
			, (SELECT StyleCode FROM Product WHERE ProdCode = @ProdCode) AS StyleCode -- StyleCode (스타일코드)
			, B.BrandCode + B.ShopCode + @CurrDate  -- StockMisNo (상품착오번호 = 브랜드코드 + 매장코드 + 날짜)
			, B.BrandSuppRecvDate      -- RecvSuppDate (출고일)
			, null          -- ShopDesc
			, null          -- BrandDesc
			, B.PlantCode        -- PlantCode
			, @ErrorRegEmpID            -- InsertEmpID 등록자    WangPengDa 20100127
			, @ErrorRegEmpID            -- RecvEmpID 접수자    WangPengDa 20100127
			, null          -- StockMisRecvDateTime (smalldatetime)
			, null          -- ProcessEmpID 처리자
			, null          -- StockMisProcDateTime
			, null          -- NewSAPDeliveryNo
			, CASE WHEN (@ShopInFixQty - @ShopRecvSuppQty ) < 0 THEN 'S' ELSE 'R' END -- RecvSuppType 물류출고수량 > 매장입고수량 : 'S' else 'R'
			, B.BoxNo         -- BoxNo
			, 0          -- DelChk
			, @ErrorRegUserID        -- InUserID
			, GETDATE()        -- InDateTime
			, @ErrorRegUserID        -- ModiUserID
			, GETDATE()        -- ModiDateTime
			, @ErrorRegEmpName      -- RecvEmpName  접수자
			, Null          -- ProcessEmpName 처리자
			, @ErrorRegEmpName      -- InsertEmpName 등록자
			, ''          -- SendState
			, 'R'          -- SendFlag
			, null
			,@CurrDate         -- SendDateTime
			,@SuppEmpName        -- SendEmpName 王鹏达 20100629
			,B.WayBillNo                               -- WayBillNo01 운송장번호 김재훈 20110503
			FROM RecvSuppDtl AS A
		INNER JOIN RecvSuppMst AS B ON (B.RecvSuppNo = A.RecvSuppNo)
		INNER JOIN Product    AS C ON (C.BrandCode = A.BrandCode AND C.ProdCode = A.ProdCode)
		WHERE B.RecvSuppNo = @RecvSuppNo

		END
		ELSE
		BEGIN
		/************************ 데이터가 많을 경우 ************************
		이 경우는 무조건 입고확정 정보가 상품별로 모두 존재하는 경우임.
		*********************************************************************/
		--PRINT N'데이터가 많을 경우'

		DECLARE @RecvSuppNo_Inner   CHAR(14)
		DECLARE @SeqNo_Inner    INT
		DECLARE @ProdCode_Inner    VARCHAR(18)
		DECLARE @RecvSuppQty_Inner   INT

		DECLARE @StockMisQty_Inner   INT  -- 임시 잔여 오차수량 저장변수
		DECLARE @StockMisQtyProduct_Inner INT  -- 임시 오차수량 저장변수
		DECLARE @IsBreak_Inner    BIT  -- 중지 가능여부 저장변수
		DECLARE @SAPDeliveryNo_Inner   CHAR(10)  -- 물류팀에서 필요로 하는 값(RecvSuppDtl에서 가져옴, 신규착오는 NULL)
		DECLARE @SAPDeliveryItemNo_Inner  CHAR(10)  -- 물류팀에서 필요로 하는 값(RecvSuppDtl에서 가져옴, 신규착오는 NULL)

		SET @StockMisQty_Inner = @ShopInFixQty - @ShopRecvSuppQty
		--PRINT N'00 @StockMisQty_Inner' + CONVERT(VARCHAR(10), @StockMisQty_Inner)

		DECLARE @StockMisCursor CURSOR

		-- 대상 상품에 대한 입출고정보 결과 커서 ( 예측되는 결과 : 1..*, 상품이 여러개의 레코드에 포함되는 경우가 발생함 )
		SET @StockMisCursor = CURSOR LOCAL SCROLL FOR
		SELECT 
		MST.RecvSuppNo, 
		MST.SeqNo, 
		DTL.ProdCode, 
		DTL.RecvSuppQty, 
		DTL.SAPDeliveryNo, 
		DTL.SAPDeliveryItemNo, 
		MST.SuppEmpName
		FROM RecvSuppMst AS MST
		INNER JOIN RecvSuppDtl AS DTL
		ON (MST.RecvSuppNo = DTL.RecvSuppNo AND DTL.DelChk = 0)
		WHERE MST.BrandCode = @BrandCode AND MST.ShopCode = @ShopCode
		AND ShippingTypeCode IN ('01','66','16')  --wangpengda 20100203
		AND MST.DelChk = 0
		AND MST.WayBillNo = @WaybillNo
		AND MST.BoxNo = @WaybillNo
		AND DTL.ProdCode = @ProdCode					--modify by li.guolin 20151013
		AND MST.ShopSuppRecvDate=@RecvDate			--ADD BY SHEN.XUE 20140804
		OPEN @StockMisCursor

		FETCH NEXT FROM @StockMisCursor
		INTO @RecvSuppNo_Inner,
		@SeqNo_Inner, 
		@ProdCode_Inner, 
		@RecvSuppQty_Inner, 
		@SAPDeliveryNo_Inner, 
		@SAPDeliveryItemNo_Inner, 
		@SuppEmpName

		-- 최초 한번 동작해야함.
		SET @IsBreak_Inner = '0'

		WHILE @@FETCH_STATUS = 0
		BEGIN
		SELECT @DtSeq = ISNULL(MAX(DTL.DtSeq),0) + 1 FROM StockMisDtl AS DTL
		WHERE DTL.BrandCode = @BrandCode   --modify by li.guolin 20151013
			AND DTL.ShopCode = @ShopCode   --modify by li.guolin 20151013
			AND DTL.SeqNo = @SeqNo_Inner  --modify by li.guolin 20151013
			AND DTL.Dates=@CurrDate  --modify by li.guolin 20151013
		IF ( @IsBreak_Inner = '1' )
		BEGIN
		--PRINT N'BREAK'
		BREAK
		END


		--PRINT N'01 @StockMisQty_Inner' + CONVERT(VARCHAR(10), @StockMisQty_Inner)

		-- "해당상품의 출고수량 총합 >= 오차수량" 의 조건에 위배되면 커서실행을 멈춘다.
		--  IF (@ProdCode_Inner=@ProdCode)
		--BEGIN

			IF ( @StockMisQty_Inner < 0)
			BEGIN
			IF ( @RecvSuppQty_Inner < ABS(@StockMisQty_Inner) )
			BEGIN
			SET @StockMisQtyProduct_Inner = @RecvSuppQty_Inner
			SET @StockMisQty_Inner = ( ABS(@StockMisQty_Inner) - @RecvSuppQty_Inner ) * -1
			--PRINT N'IF @StockMisQtyProduct_Inner' + CONVERT(VARCHAR(10), @StockMisQtyProduct_Inner)
			END
			ELSE
			BEGIN
			SET @StockMisQtyProduct_Inner = ISNULL(ABS(@StockMisQty_Inner),0)
			SET @IsBreak_Inner = '1'
			--PRINT N'ELSE @StockMisQtyProduct_Inner' + CONVERT(VARCHAR(10), @StockMisQtyProduct_Inner)
			END
			END
			ELSE
			BEGIN
			SET @IsBreak_Inner = '1'
			END

		--PRINT N'02 @StockMisQty_Inner' + CONVERT(VARCHAR(10), @StockMisQty_Inner)

		-- 오차 데이터 생성
			INSERT StockMisDtl
				( BrandCode
				, ShopCode
				, Dates
				, SeqNo
				, DtSeq
				, SAPDeliveryNo
				, SAPDeliveryItemNo
				, ProdCode
				, RecvSuppQty
				, StockMisQty
				, StockMisFixQty
				, StockMisReasonCode
				, StockMisStatusCode
				, StockMisResultCode
				, RecvDate
				, RecvSuppNo
				, StyleCode
				, StockMisNo
				, RecvSuppDate
				, ShopDesc
				, BrandDesc
				, PlantCode
				, InsertEmpID
				, RecvEmpID
				, StockMisRecvDateTime
				, ProcessEmpID
				, StockMisProcDateTime
				, NewSAPDeliveryNo
				, RecvSuppType
				, BoxNo
				, DelChk
				, InUserID
				, InDateTime
				, ModiUserID
				, ModiDateTime
				, RecvEmpName
				, ProcessEmpName
				, InsertEmpName
				, SendState
				, SendFlag
				, SendDateTime
				, InvtBaseDate
				, SendEmpName
				, WayBillNo01                              -- 운송장번호 김재훈 20110503
				)
			SELECT B.BrandCode
				, B.ShopCode        -- 매장코드
				, @CurrDate        -- 현재일자 Dates
				, A.SeqNo         -- SeqNo
				, @DtSeq         -- DtSeq -- 상세순번
				, @SAPDeliveryNo_Inner       -- SAPDeliveryNo
				, @SAPDeliveryItemNo_Inner      -- SAPDeliveryItemNo
				, @ProdCode_Inner       --, A.ProdCode        -- ProdCode
				, @RecvSuppQty_Inner      -- RecvSuppQty 출고수량
				, ABS(@StockMisQtyProduct_Inner)   -- StockMisQty 착오수량
				, 0          -- StockMisFixQty 착오확정수량
				, Null          -- StockMisReasonCode 착오
				, '10'          -- StockMisStatusCode 착오 처리 상태(10: 매장입고)
				, Null          -- StockMisResultCode
				, B.ShopSuppRecvDate      -- RecvDate
				, @RecvSuppNo_Inner      -- RecvSuppNo
				, (SELECT StyleCode FROM Product WHERE ProdCode = @ProdCode) AS StyleCode -- StyleCode (스타일코드)
				, B.BrandCode + B.ShopCode + @CurrDate  -- StockMisNo (상품착오번호 = 브랜드코드 + 매장코드 + 날짜)
				, B.BrandSuppRecvDate      -- RecvSuppDate (출고일)
				, null          -- ShopDesc
				, null          -- BrandDesc
				, B.PlantCode        -- PlantCode
				, @ErrorRegEmpID            -- InsertEmpID      WangPengDa 20100127
				, @ErrorRegEmpID            -- RecvEmpID        WangPengDa 20100127
				, null          -- StockMisRecvDateTime (smalldatetime)
				, null          -- ProcessEmpID
				, null          -- StockMisProcDateTime
				, null          -- NewSAPDeliveryNo
				, CASE WHEN (@ShopInFixQty - @ShopRecvSuppQty ) < 0 THEN 'S' ELSE 'R' END -- RecvSuppType 물류출고수량 > 매장입고수량 : 'S' else 'R'
				, B.BoxNo         -- BoxNo
				, 0          -- DelChk
				, @ErrorRegUserID        -- InUserID
				, GETDATE()        -- InDateTime
				, @ErrorRegUserID        -- ModiUserID
				, GETDATE()        -- ModiDateTime
				, @ErrorRegEmpName      -- RecvEmpName
				, Null          -- ProcessEmpName
				, @ErrorRegEmpName      -- InsertEmpName
				, ''          -- SendState
				, 'R'          -- SendFlag
				, null
				,B.InvtBaseDate        -- SendDateTime
				,@SuppEmpName        --SendEmpName
					,B.WayBillNo                               -- WayBillNo01 운송장번호 김재훈 20110503
				FROM RecvSuppDtl AS A
				INNER JOIN RecvSuppMst AS B ON (B.RecvSuppNo = A.RecvSuppNo)
				INNER JOIN Product    AS C ON (C.BrandCode = A.BrandCode AND C.ProdCode = A.ProdCode)
				WHERE B.RecvSuppNo = @RecvSuppNo_Inner
				AND A.ProdCode = @ProdCode_Inner
		--END
		FETCH NEXT FROM @StockMisCursor
		INTO @RecvSuppNo_Inner,
			@SeqNo_Inner, 
			@ProdCode_Inner, 
			@RecvSuppQty_Inner, 
			@SAPDeliveryNo_Inner, 
			@SAPDeliveryItemNo_Inner, 
			@SuppEmpName
		END

		CLOSE @StockMisCursor
		DEALLOCATE @StockMisCursor
		END


		END TRY
		BEGIN CATCH

		EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2
		END CATCH
		END;
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("createInsertStockInMissSaveStockMisDtlC1 error: %v", err.Error())
		log.Println()
	}
}
