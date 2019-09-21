package test

import (
	"clearance-adapter/factory"
	"fmt"
	"log"
)

func initReturnToWarehouseData() {
	createOrderControlTable()
	createWaybillNoTable()
	createStyleTable()
	createReturnToWarehouseMasterSP()
	createReturnToWarehouseDetailSP()
	createBrandTable()
}

func createOrderControlTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createOrderControlTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.OrderControl;"); err != nil {
		log.Printf("createOrderControlTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE OrderControl
		(
			OrderControlNo CHAR(12) PRIMARY KEY NOT NULL,
			ShippingTypeCode CHAR(2) NOT NULL,
			BrandCode VARCHAR(4),
			OrderControlName NVARCHAR(100),
			StyleRangeChk BIT,
			ShopRangeChk BIT,
			StartDate CHAR(8),
			EndDate CHAR(8),
			ExpectedLeadTime SMALLINT,
			StorageLoc CHAR(4),
			PlantCode CHAR(4),
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME,
			ModiUserID VARCHAR(20) NOT NULL,
			ModiDateTime DATETIME,
			SubStyleTypeCode CHAR NOT NULL,
			DelChk BIT NOT NULL
		);

		INSERT INTO CSL.dbo.OrderControl 
		(OrderControlNo, ShippingTypeCode, BrandCode, OrderControlName, StyleRangeChk, ShopRangeChk, StartDate, EndDate, ExpectedLeadTime, StorageLoc, PlantCode, InUserID, InDateTime, ModiUserID, ModiDateTime, SubStyleTypeCode, DelChk) 
		VALUES 
		('SA1811220001', '41', 'SA', N'SA随时退仓', 1, 1, '20190101', '99991231', 0, '1300', '1201', 'CNFBGLEA02', '2019-08-12 10:18:00.000', ' ', '2019-08-12 10:18:00.000', 'P', 0),
		('Q31811220001', '41', 'Q3', N'Q3随时退仓', 1, 1, '20190101', '99991231', 0, '1300', '1201', 'CNFBGLEA02', '2019-08-14 10:18:00.000', ' ', '2019-08-14 10:18:00.000', 'P', 0),
		('SA1811220002', '42', 'SA', N'SA季节退仓', 1, 1, '20190101', '99991231', 0, '1300', '1201', 'CNFBGLEA02', '2019-08-12 10:18:00.000', ' ', '2019-08-12 10:18:00.000', 'P', 0),
		('Q31811220002', '42', 'Q3', N'Q3季节退仓', 1, 1, '20190101', '99991231', 0, '1300', '1201', 'CNFBGLEA02', '2019-08-14 10:18:00.000', ' ', '2019-08-14 10:18:00.000', 'P', 0),
		('SA1811220003', '47', 'SA', N'SA次品退仓', 1, 1, '20190101', '99991231', 0, '1300', '1201', 'CNFBGLEA02', '2019-08-12 10:18:00.000', ' ', '2019-08-12 10:18:00.000', 'P', 0),
		('Q31811220003', '47', 'Q3', N'Q3次品退仓', 1, 1, '20190101', '99991231', 0, '1300', '1201', 'CNFBGLEA02', '2019-08-14 10:18:00.000', ' ', '2019-08-14 10:18:00.000', 'P', 0);
	`
	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createOrderControlTable error: %v", err.Error())
		fmt.Println()
	}
}

func createWaybillNoTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createWaybillNoTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.WayBillNo;"); err != nil {
		log.Printf("createWaybillNoTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE WayBillNo
		(
			ShippingCompanyCode CHAR(2) NOT NULL,
			WayBillNo VARCHAR(30) NOT NULL,
			AllowDulpChk BIT,
			InUserID VARCHAR(20),
			InDateTime DATETIME,
			ModiUserID VARCHAR(20),
			ModiDateTime DATETIME
		);
	`
	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createWaybillNoTable error: %v", err.Error())
		fmt.Println()
	}
}

func createStyleTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createStyleTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.Style;"); err != nil {
		log.Printf("createStyleTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE Style
		(
			BrandCode VARCHAR(4) NOT NULL,
			StyleCode VARCHAR(18) NOT NULL,
			ItemSCode CHAR(2),
			StyleName NVARCHAR(100),
			Year CHAR(4),
			SeasonCode CHAR(4),
			MateTypeCode CHAR(4) NOT NULL,
			SubStyleTypeCode CHAR NOT NULL,
			NAONAChk BIT NOT NULL,
			SupGroupCode CHAR(2) NOT NULL,
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME,
			UseChk BIT,
			SaleMonthCode CHAR(2),
			PreStyleCode VARCHAR(18),
			ModiDateTime DATETIME
		);

		INSERT INTO CSL.dbo.Style 
		(BrandCode, StyleCode, ItemSCode, StyleName, Year, SeasonCode, MateTypeCode, SubStyleTypeCode, NAONAChk, SupGroupCode, InUserID, InDateTime, UseChk, SaleMonthCode, PreStyleCode, ModiDateTime) 
		VALUES 
		('Q3', 'Q3AFAFDU6S21', 'PT', '男女同款拖鞋', '2016', 'F   ', 'HAWA', 'P', 0, '01', 'system', '2018-12-26 05:08:02.593', 1, 'A ', ' ', '2019-08-12 05:31:17.437'),
		('SA', 'SPWH936D54', 'WH', '蜡笔小新彩色短裙', '2019', '3   ', 'FERT', 'P', 0, '01', 'system', '2019-04-22 00:00:00.000', 1, '6 ', ' ', '2019-08-12 05:31:17.437'),
		('SA', 'SPWJ948S22', 'WJ', '女式牛仔裙', '2019', '4   ', 'FERT', 'P', 0, '01', 'system', '2019-05-16 00:00:00.000', 1, '8 ', ' ', '2019-08-12 05:31:17.437'),
		('SA', 'SPWJ948S23', 'WJ', '女式牛仔裙', '2019', '4   ', 'FERT', 'P', 0, '01', 'system', '2019-05-16 00:00:00.000', 1, '8 ', ' ', '2019-08-12 05:31:17.437'),
		('SA', 'SPYC949H21', 'YC', '大格纹衬衫', '2019', '4   ', 'FERT', 'P', 0, '01', 'system', '2019-05-17 00:00:00.000', 1, '9 ', ' ', '2019-08-12 05:31:17.437'),
		('SA', 'SPYC949S11', 'YC', '女式格子衬衫', '2019', '4   ', 'FERT', 'P', 0, '01', 'system', '2019-05-15 00:00:00.000', 1, '9 ', ' ', '2019-08-12 05:31:17.437'),
		('SA', 'SPYS949H22', 'YS', '条纹衬衫', '2019', '4   ', 'FERT', 'P', 0, '01', 'system', '2019-05-17 00:00:00.000', 1, '9 ', ' ', '2019-08-12 05:31:17.437');
	`
	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createStyleTable error: %v", err.Error())
		fmt.Println()
	}
}

func createReturnToWarehouseMasterSP() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createReturnToWarehouseMasterSP error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP PROCEDURE IF EXISTS dbo.up_CSLK_IOM_InsertReturnGoodReservation_RecvSuppMst_C1_Clearance"); err != nil {
		log.Printf("createReturnToWarehouseMasterSP error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE PROCEDURE [dbo].[up_CSLK_IOM_InsertReturnGoodReservation_RecvSuppMst_C1_Clearance]
			@BrandCode			VARCHAR(4)  = NULL   
			,@ShopCode				CHAR(4)		= NULL
			,@OutDate						CHAR(8) = null -- 退仓出库日期
			,@WayBillNo			VARCHAR(13) = NULL
			,@ShippingTypeCode		CHAR(2)		= NULL   -- 출하유형
			,@ShippingCompanyCode  CHAR(2)     = NULL
			,@EmpID	CHAR(10) = NULL	 -- 등록자
			,@DeliveryID VARCHAR(250)=NULL -- 春风的billno moidfy by li.guolin 20170811
			,@DeliveryOrderNo VARCHAR(250)=NULL -- 春风的orderno moidfy by li.guolin 20170811
		AS     
		--SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED        
		SET XACT_ABORT ON     
		SET NOCOUNT ON  
		BEGIN  
			DECLARE @NewRecvSuppNo		CHAR(14)
			DECLARE @NewSeq			INT
			DECLARE @OrderControlNo    CHAR(12)
			DECLARE @SendState			VARCHAR(2)
			DECLARE @SendFlag			CHAR(1)
			DECLARE @PlantCode			CHAR(4)    -- 물류창고
			DECLARE @RecvSuppType		CHAR(1)    -- 입출고구분
			DECLARE @SAPMenuType		CHAR(1)
			DECLARE @UserID VARCHAR(20);
			DECLARE @EmpName NVARCHAR(200);
			DECLARE @ErrorCode			NVARCHAR(1000) = '';
			DECLARE @ErrorParam1		NVARCHAR(4000) = '';
			DECLARE @ErrorParam2		NVARCHAR(4000) = '';
			DECLARE @TransTypeCode		CHAR(1)
			DECLARE @RecvSuppStatusCode	CHAR(1)
			DECLARE @NormalProductType	CHAR(1)
			DECLARE @BrandSuppRecvDate	VARCHAR(8)	= NULL
			DECLARE @VolumeType NVARCHAR(20)=N'中箱子' --moidfy by li.guolin 20170811
			DECLARE @VolumesSize VARCHAR(20)='0.09486' --moidfy by li.guolin 20170811
			DECLARE @VolumesUnit NVARCHAR(10)=N'm³' --moidfy by li.guolin 20170811
			DECLARE @BoxAmount  int=1 --moidfy by li.guolin 20170811
			DECLARE @BoxType CHAR(2)='MB'  --moidfy by li.guolin 20171009
		
		BEGIN TRY
		
		IF(@OutDate = '' OR @OutDate IS NULL)
		BEGIN 
			SET @OutDate = CONVERT(CHAR(8),GETDATE(),112)
		END
		
			IF(@ShippingTypeCode = '47')
				BEGIN
					SET @NormalProductType='B'
				END
			ELSE
				BEGIN
					SET @NormalProductType='A'
				END
		
			IF (@ShippingCompanyCode = 'SR')
				BEGIN
					SET @VolumeType = N'中箱子'
					SET @VolumesSize = '0.09486'
					SET @VolumesUnit = N'm³' --
					SET @BoxAmount  = 1 --moidfy by li.g
					SET @BoxType = 'MB'  --moidfy by
			END
			ELSE
				BEGIN
					SET @VolumeType = NULL
					SET @VolumesSize = NULL
					SET @VolumesUnit = NULL --
					SET @BoxAmount  = 1 --moidfy by li.g
					SET @BoxType = NULL  --moidfy by
				END
		
		
			-- 마감 체크  
			IF (left(dbo.udf_CSLK_MonthlyClosingChk('02','Zn'),1) = 1)  
			BEGIN  
				SELECT @ErrorCode = SUBSTRING(dbo.udf_CSLK_MonthlyClosingChk('02','Zn'),2,510)  
					EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2  
				
			END  
			
									
			IF @ShippingCompanyCode=''
				SET @ShippingCompanyCode=NULL;
		
		
			--zhang.fengcheng
			--同一运单号下的箱号不能超过999
			IF (SELECT COUNT(DISTINCT BoxNo) FROM RecvSuppMst WITH(NOLOCK) WHERE  WayBillNo = @WayBillNo AND DelChk = 0)>=999 	
			BEGIN  
				SET @ErrorCode = 'IOM184'  
				SET @ErrorParam1 = 'ReturnGoodReservation'  
				EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2  
			END 
		
			-- 순번 생성
			SELECT @NewSeq = IsNull(MAX(SeqNo),0) + 1  
				FROM RecvSuppMst 
			WHERE BrandCode = @BrandCode  
				AND ShopCode = @ShopCode  
				AND Dates = @OutDate
				AND SeqNo <6000  			--2015年6月29日14:05:35翟  
				
			--ShopCode(4)+yymmdd+9999(14자리)     
			SET @NewRecvSuppNo = @ShopCode + Right(@OutDate,6)+ Right(replicate('0',4) + convert(varchar,@NewSeq),4)
		
		/*
			임시 저장 모드로 내용 변경
		*/
			SET @RecvSuppType = 'S'   -- 입출고구분 S:출고/ R:입고  
			SET @TransTypeCode = '5'  -- 운송타입  
			SET @RecvSuppStatusCode = 'R' -- 卖场退仓确定
			
			-- 통제번호/반품통제 창고 가져오기  
			SELECT @OrderControlNo = OrderControlNo  
				, @PlantCode = PlantCode  
				FROM OrderControl  WITH(NOLOCK)
			WHERE ShippingTypeCode = @ShippingTypeCode  
				AND BrandCode = @BrandCode  
				AND StartDate <=  @OutDate
				AND EndDate >= @OutDate
				AND DelChk = 0
			
			-- 창고코드 : 통제에 창고코드가 없는경우   
			IF @PlantCode IS NULL OR @PlantCode = ''  
			BEGIN  
				SELECT @PlantCode = PlantCode  
					FROM Brand  WITH(NOLOCK)
				WHERE BrandCode = @BrandCode  
			END   
			
			-- 통제대상 여부 확인  
			IF @OrderControlNo IS NULL OR @OrderControlNo = ''  
			BEGIN  
				print '1'
				SET @ErrorCode = 'COM051'  
				SET @ErrorParam1 = 'ReturnGoodReservation'  
				EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2  
			END
		
		/*
			임시 저장 모드로 내용 변경
		*/  
			-- SAP관련세팅  
			SET @SendFlag = 'R'		-- SAP 인터 페이스 방지
			SET @SendState = ''   
			SET @SAPMenuType = '3'		-- 반품:3  
		
		
			SELECT @EmpName = EmpName
				FROM Employee
				WHERE EmpID = @EmpID
		
		
				SELECT @UserID = UserID
				FROM UserInfo
				WHERE EmpID = @EmpID
		
		
				IF EXISTS ( SELECT 1 FROM RecvSuppMst  WITH(NOLOCK)
										WHERE BrandCode = @BrandCode 
											AND ShopCode = @ShopCode
											AND WayBillNo = @WayBillNo
											AND BoxNo = RTRIM(@WayBillNo)
											AND ShippingCompanyCode=@ShippingCompanyCode
											AND DelChk = 0
											--AND RecvSuppStatusCode<> 'T'
											)
				BEGIN
					SELECT RecvSuppNo AS RecvSuppNo FROM RecvSuppMst  WITH(NOLOCK)
										WHERE BrandCode = @BrandCode
											AND ShopCode = @ShopCode
											AND WayBillNo = @WayBillNo
											AND BoxNo = RTRIM(@WayBillNo)
											AND ShippingCompanyCode=@ShippingCompanyCode
											AND DelChk = 0
					RETURN
				END
		
			-- 입출고 공통 (반품Master)  
			INSERT INTO RecvSuppMst  
						(
							RecvSuppNo
						,BrandCode
						,ShopCode
						,Dates
						,SeqNo
						,RecvSuppType
						,NormalProductType
						,ShopSuppRecvDate
						,TransTypeCode
						,ShippingTypeCode
						,WayBillNo
						,RecvSuppStatusCode
						,BoxNo
						,PlantCode
						,OrderControlNo
						,RecvEmpID
						,RecvEmpName
						,SuppEmpID
						,SuppEmpName
						,InUserID
						,InDateTime
						,ModiUserID
						,ModiDateTime
						,BrandSuppRecvDate
						,SAPMenuType
						,SendFlag
						,InvtBaseDate
						,ShippingCompanyCode --add by li.guolin 20170214
						,DeliveryID
						,DeliveryOrderNo
						,VolumeType
						,VolumesSize
						,VolumesUnit
						,BoxAmount
						,ProvinceCode
						,CityCode
						,DistrictCode
						,Area
						,ShopManagerName
						,MobilePhone
						,BoxType
						,Channel --add by zhang.xinshuai 2018-1-26 13:41:57
						)  
				VALUES(
							@NewRecvSuppNo
							,@BrandCode
							,@ShopCode
							,@OutDate
							,@NewSeq
							,@RecvSuppType
							,@NormalProductType
							,@OutDate
							,@TransTypeCode
							,@ShippingTypeCode
							,@WayBillNo
							,@RecvSuppStatusCode
							,@WayBillNo
							,@PlantCode
							,@OrderControlNo
							,''
							,''
							,@EmpID
							,@EmpName
							,@UserID     -- 등록자
							,GETDATE()
							,@UserID     -- 수정자
							,GETDATE()
							,@BrandSuppRecvDate
							,@SAPMenuType
							,@SendFlag
							,@OutDate
							,@ShippingCompanyCode --add by li.guolin 20170214
							,@DeliveryID
							,@DeliveryOrderNo
							,@VolumeType
							,@VolumesSize
							,@VolumesUnit
							,@BoxAmount
							,'SHH'
							,'SHH01'
							,'SHH01015'
							,N'莲花南路3130号 衣恋物流中心'
							,N'张殿利'
							,'13917778786'
							,@BoxType
							,'Clearance'
						)    -- 재고기준일자  
		
			/*
			添加物流信息
			*/
				IF NOT EXISTS(
					SELECT 1
					FROM WayBillNo
					WHERE ShippingCompanyCode=@ShippingCompanyCode
					AND WayBillNo=@WayBillNo
				)
				BEGIN
					INSERT INTO WayBillNo
					(
						ShippingCompanyCode, WayBillNo, AllowDulpChk, InUserID, InDateTime, ModiUserID, ModiDateTime
					)
					VALUES
					(
						@ShippingCompanyCode,@WayBillNo,0,@UserID,GETDATE(),@UserID,GETDATE()
					)
				END
		
			SELECT @NewRecvSuppNo AS RecvSuppNo   
			
		END TRY  
		BEGIN CATCH  
			
			EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2  
			
		END CATCH   
		END
	`

	if _, err := session.Exec(sql); err != nil {
		log.Printf("createReturnToWarehouseMasterSP error: %v", err.Error())
		log.Println()
	}
}

func createReturnToWarehouseDetailSP() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createReturnToWarehouseMasterSP error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP PROCEDURE IF EXISTS dbo.up_CSLK_IOM_InsertReturnGoodReservation_RecvSuppDtl_C1_Clearance"); err != nil {
		log.Printf("createReturnToWarehouseMasterSP error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE PROCEDURE [dbo].[up_CSLK_IOM_InsertReturnGoodReservation_RecvSuppDtl_C1_Clearance]
			@RecvSuppNo							CHAR(14) 			= NULL
			,@RecvSuppSeqNo						INT		 				= NULL
			,@BrandCode								VARCHAR(4) 		= NULL
			,@ShopCode								CHAR(4)				= NULL
			,@Dates						  			CHAR(8)				= NULL
			,@ProdCode								VARCHAR(18) 	= NULL
			,@RecvSuppQty							INT 					= NULL
			,@EmpID										CHAR(10) 			= NULL
			,@AbnormalProdReasonCode	CHAR(2) 			= NULL
			,@AbnormalChkCode    			CHAR(2) 			= NULL
			,@AbnormalSerialNo   			VARCHAR(7) 		= NULL


			AS   
			SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED      
			SET XACT_ABORT ON   
			SET NOCOUNT ON
			BEGIN
			
			DECLARE @ErrorCode   NVARCHAR(50) ='';
			DECLARE @ErrorParam1 NVARCHAR(4000)='';
			DECLARE @ErrorParam2 NVARCHAR(4000)='';

			DECLARE @SendFlag		 CHAR(1)
			DECLARE @SeqNo			 INT = NULL
			DECLARE @UserID		   VARCHAR(20) = NULL
			DECLARE @EmpName     NVARCHAR(200)
			DECLARE @SalePrice	DECIMAL(19,2) = NULL

			BEGIN TRY	

				IF (@RecvSuppSeqNo IS NULL OR @RecvSuppSeqNo = 0)
				BEGIN
					IF( EXISTS (select 1 from RecvSuppDtl where RecvSuppNo=@RecvSuppNo))
						BEGIN
						SET @RecvSuppSeqNo=(select MAX(RecvSuppSeqNo)+1 from RecvSuppDtl where RecvSuppNo=@RecvSuppNo)
						END
					ELSE
						BEGIN
								SET @RecvSuppSeqNo=1
								SET @SeqNo=1
						END 

				END

				SET @SendFlag = 'R'
				
				IF(@Dates = '' or @Dates IS NULL)   --WANGPENGDA 20091126
					SET @Dates = CONVERT(CHAR(8),GETDATE(),112)
				SET @SeqNo=@RecvSuppSeqNo


				------------------------------------------------------------
				-- 공급구분/단가구분 가져오기
				------------------------------------------------------------
				DECLARE @SupGroupCode	char(2)		-- 공급구분
				DECLARE @PriceTypeCode	char(2)		-- 단가구분
				DECLARE @SaipType		char(2)		-- 사입

				SET @SaipType = '00'
				SET @PriceTypeCode = '01'

				-- 공급구분/단가구분 가져오기
				SELECT @SupGroupCode = B.SupGroupCode
				From Product AS A 
					INNER JOIN Style AS B On (B.BrandCode = A.BrandCode AND B.StyleCode = A.StyleCode)
				WHERE A.BrandCode = @BrandCode
					AND A.ProdCode = @ProdCode

				SELECT @EmpName = EmpName
					FROM Employee
					WHERE EmpID = @EmpID

				SELECT @UserID = UserID
					FROM UserInfo
					WHERE EmpID = @EmpID


				SELECT @SalePrice = Price
				FROM Product
				WHERE BrandCode = @BrandCode
				AND ProdCode = @ProdCode

				-- 입출고 상세 (반품상세)
				INSERT INTO RecvSuppDtl
				(
						RecvSuppNo
					,RecvSuppSeqNo
					,BrandCode
					,ShopCode
					,Dates
					,SeqNo
					,ProdCode
					,RecvSuppQty
					,SalePrice
					,AbnormalProdReasonCode
					,SupGroupCode
					,PriceTypeCode
					,SaipType
					,InUserID
					,InDateTime
					,ModiUserID
					,ModiDateTime
					,SendFlag
					,AbnormalChkCode
					,AbnormalSerialNo

				)
				VALUES 
				(
						@RecvSuppNo
					,@RecvSuppSeqNo
					,@BrandCode
					,@ShopCode
					,@Dates
					,@SeqNo
					,@ProdCode
					,@RecvSuppQty
					,@SalePrice
					,CASE WHEN @AbnormalProdReasonCode = '' THEN NULL ELSE @AbnormalProdReasonCode END
					,@SupGroupCode
					,@PriceTypeCode
					,@SaipType			
					,@UserID
					,GETDATE()
					,@UserID
					,GETDATE()
					,@SendFlag
					,@AbnormalChkCode
					,@AbnormalSerialNo
				)

			END TRY
			BEGIN CATCH

				EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2
			END CATCH 
			END
	`

	if _, err := session.Exec(sql); err != nil {
		log.Printf("createReturnToWarehouseMasterSP error: %v", err.Error())
		log.Println()
	}
}

func createBrandTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createBrandTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.Brand;"); err != nil {
		log.Printf("createBrandTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE Brand
		(
			BrandCode VARCHAR(4) PRIMARY KEY NOT NULL,
			BrandName NVARCHAR(200),
			Initial VARCHAR(2),
			CompanyCode CHAR(4),
			RivalBrandDisplayOrder INT NOT NULL,
			PlantCode CHAR(4) NOT NULL,
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME,
			UseChks BIT
		);

		INSERT INTO 
		CSL.dbo.Brand (BrandCode, BrandName, Initial, CompanyCode, RivalBrandDisplayOrder, PlantCode, InUserID, InDateTime, UseChks) 
		VALUES 
		('Q3', 'SHOOPEN', 'Q3', 'F201', 0, '1222', 'system', '2015-06-03 05:00:01.140', 1),
		('SA', 'SPAO(CN)', 'SA', 'F201', 0, '1200', 'system', '2013-05-04 05:00:03.410', 1);
	`
	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createBrandTable error: %v", err.Error())
		fmt.Println()
	}
}
