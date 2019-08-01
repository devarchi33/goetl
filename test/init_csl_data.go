package test

import (
	"clearance-adapter/factory"
	"clearance-adapter/models"
	"fmt"
	"log"
	"path/filepath"
)

func initCSL() {
	// createCSLDB()
	createRecvSuppMstTable()
	setRecvSuppMstData()
	createRecvSuppDtlTable()
	setRecvSuppDtlData()
	createMonthlyBizFuctionClosingTable()
	createIFConfigTable()
	createUserInfoTable()
	createEmployeeTable()
	createSP()
}

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

// Master 部分
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
}

func setRecvSuppMstData() {
	filename, err := filepath.Abs("test/data/test_in_storage_etl_RecvSuppMst_data.csv")
	if err != nil {
		panic(err.Error())
	}
	headers, data := readDataFromCSV(filename)
	masters := buildRecvSuppMsts(headers, data)

	loadRecvSuppMstData(masters)
}

func buildRecvSuppMsts(headers map[int]string, data [][]string) []models.RecvSuppMst {
	masters := make([]models.RecvSuppMst, 0)
	for _, row := range data {
		master := new(models.RecvSuppMst)
		setObjectValue(headers, row, master)
		masters = append(masters, *master)
	}

	return masters
}

func loadRecvSuppMstData(masters []models.RecvSuppMst) {
	for _, master := range masters {
		if affected, err := factory.GetCSLEngine().Insert(&master); err != nil {
			fmt.Printf("loadRecvSuppMstData error: %v", err.Error())
			fmt.Println()
			fmt.Printf("affected: %v", affected)
			fmt.Println()
		}
	}
}

// Detail 部分
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
}

func setRecvSuppDtlData() {
	filename, err := filepath.Abs("test/data/test_in_storage_etl_RecvSuppDtl_data.csv")
	if err != nil {
		panic(err.Error())
	}
	headers, data := readDataFromCSV(filename)
	details := buildRecvSuppDtls(headers, data)

	loadRecvSuppDtlData(details)
}

func buildRecvSuppDtls(headers map[int]string, data [][]string) []models.RecvSuppDtl {
	details := make([]models.RecvSuppDtl, 0)
	for _, row := range data {
		detail := new(models.RecvSuppDtl)
		setObjectValue(headers, row, detail)
		details = append(details, *detail)
	}

	return details
}

func loadRecvSuppDtlData(details []models.RecvSuppDtl) {
	for _, detail := range details {
		if affected, err := factory.GetCSLEngine().Insert(&detail); err != nil {
			fmt.Printf("loadRecvSuppDtlData error: %v", err.Error())
			fmt.Println()
			fmt.Printf("affected: %v", affected)
			fmt.Println()
		}
	}
}

// SP 部分
func createSP() {
	createMonthlyClosingChk()
	createComonRaiseError()
	createChkRecvSupp()
	createUpdateStockInEnterConfirmSaveRecvSuppMstR1ClearanceByWaybillNo()
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
		@ShopCode CHAR(4),       -- 매장코드
		@WaybillNo VARCHAR(13),  -- 운송장번호
		@UserID VARCHAR(20)      -- 入库人ID
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
			DECLARE @CurrDate CHAR(8);
			DECLARE @RecvEmpName NVARCHAR(200);
			DECLARE @EmpId VARCHAR(20);
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

				--Query Start
				PRINT ('@RecvSuppNo1');
				SELECT TOP 1
					@RecvSuppNo = RecvSuppNo
				FROM RecvSuppMst
				WHERE BrandCode = @BrandCode
					AND ShopCode = @ShopCode
					AND WayBillNo = @WaybillNo
					AND ShippingTypeCode IN ( '01', '66', '16' ) --wangpengda 20100203
					AND DelChk = '0'; --20091210 WangPengDa \=-
				PRINT ('@RecvSuppNo2');
				PRINT (@RecvSuppNo);


				SELECT @EmpId = A.EmpID,
					@RecvEmpName = B.EmpName
				FROM UserInfo AS A
					INNER JOIN Employee AS B
						ON (B.EmpID = A.EmpID)
				WHERE UserID = @UserID;

				SELECT @CurrDate = CONVERT(CHAR(8), GETDATE(), 112);
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
						--SET @ErrorCode = 'IOM119';
						--SET @ErrorParam1 = '';
						--SET @ErrorParam2 = '';
						--EXEC [up_CSLK_ComonRaiseError] @ErrorCode, @ErrorParam1, @ErrorParam2;
						--注释by wang.wanyue 已入库不再是错误，直接查询数据

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
						--	   AND RecvSuppSeqNo = @RecvSuppSeqNo

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
								ShopSuppRecvDate = @CurrDate, -- 매장입고일자
								RecvSuppStatusCode = 'F',     -- 입고확정
								RecvEmpID = @EmpId,           -- 입고자사번
								RecvEmpName = @RecvEmpName,   -- 입고자명
								ModiUserID = @UserID,
								ModiDateTime = GETDATE(),
								SendFlag = 'R',
								InvtBaseDate = @CurrDate,     -- 재고기준일자.
								Channel = 'Clearance'
							WHERE BrandCode = @BrandCode
								AND ShopCode = @ShopCode
								--AND RecvSuppNo=@RecvSuppNo
								--MODIFY BY SHEN.XUE 20140806 同箱号,同运单号,允许入库 --
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

									--AND MST.RecvSuppNo=@RecvSuppNo --modify by shen.xue 20140806 同箱号，同款号，允许入库
									AND WayBillNo = @WaybillNo
									AND BoxNo = @WaybillNo
									AND ShippingTypeCode IN ( '01', '66', '16' ) --wangpengda 20100203
									AND Mst.RecvChk = 1 /* add by li.guolin 20150930 20151012 */
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
