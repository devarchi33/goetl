package test

import (
	"clearance-adapter/factory"
	"fmt"
	"log"
)

// 调货
func initTransferData() {
	createShopTable()
	createBranchTable()
	createTransferMasterSP()
	createTransferDetailSP()
	createTransferConfirmMasterSP()
	createTransferConfirmDetailSP()
}

func createBranchTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createBranchTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.Branch;"); err != nil {
		log.Printf("createBranchTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE Branch
		(
			BranchCode VARCHAR(4) PRIMARY KEY NOT NULL,
			BranchName NVARCHAR(100),
			DisplayOrder SMALLINT,
			SaleChk BIT,
			UseChk BIT DEFAULT 1 NOT NULL,
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME,
			WideBranchCode CHAR(4)
		);

		INSERT INTO CSL.dbo.Branch 
		(BranchCode, BranchName, DisplayOrder, SaleChk, UseChk, InUserID, InDateTime, WideBranchCode) 
		VALUES 
		('B411', N'Fashion华东-上海支社', null, 1, 1, 'system', '2013-12-28 05:00:02.000', 'D211'),
		('B412', N'Fashion华北-北京支社', null, 1, 1, 'system', '2013-12-28 05:00:02.000', 'D212'),
		('B413', N'Fashion华西-成都支社', null, 1, 1, 'system', '2013-12-28 05:00:02.000', 'D218'),
		('B414', N'Fashion华西-西北支社', null, 1, 1, 'system', '2013-12-28 05:00:02.000', 'D218'),
		('B415', N'Fashion华南-深圳支社', null, 1, 1, 'system', '2013-12-28 05:00:02.000', 'D215'),
		('B416', N'Fashion华东-南京支社', null, 1, 1, 'system', '2015-06-19 05:00:00.000', 'D211'),
		('B417', N'中国 E-Commerce本部', null, 1, 1, 'system', '2013-12-28 05:00:02.000', 'D217'),
		('B418', N'Fashion华北-东北支社', null, 1, 1, 'system', '2015-06-19 05:00:00.000', 'D212'),
		('B419', N'Fashion华北-天津支社', null, 1, 1, 'system', '2015-06-19 05:00:00.000', 'D212'),
		('B420', N'Fashion华南-武汉支社', null, 1, 1, 'system', '2015-06-19 05:00:00.000', 'D215'),
		('B428', N'Fashion华东-杭州支社', null, 1, 1, 'system', '2015-09-18 05:00:01.000', 'D211'),
		('B430', N'Fashion华西-重庆支社', null, 1, 1, 'system', '2016-12-17 05:00:00.000', 'D218'),
		('B437', N'SPA-电商支社(不使用)', null, 1, 1, 'system', '2013-12-28 05:00:02.000', 'D217'),
		('D211', N'华东广域支社', null, 1, 1, 'system', '2010-12-22 05:00:02.103', null),
		('D212', N'华北广域支社', null, 1, 1, 'system', '2010-12-22 05:00:02.103', null),
		('D215', N'华南广域支社', null, 1, 1, 'system', '2010-12-22 05:00:02.103', null),
		('D217', N'电商广域支社', null, 1, 1, 'system', '2013-12-28 05:00:02.273', null),
		('D218', N'华西广域支社', null, 1, 1, 'system', '2016-12-17 05:00:00.780', null);
	`

	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createShopTable error: %v", err.Error())
		fmt.Println()
	}
}

func createShopTable() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createShopTable error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP TABLE IF EXISTS CSL.dbo.Shop;"); err != nil {
		log.Printf("createShopTable error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE Shop
		(
			BrandCode VARCHAR(4) NOT NULL,
			ShopCode CHAR(4) NOT NULL,
			ShopName NVARCHAR(100),
			VendorGroupCode VARCHAR(3),
			VendorGroupShopCode VARCHAR(10),
			Floor SMALLINT,
			Address NVARCHAR(400),
			PosUseChk BIT NOT NULL,
			ShopStatusCode CHAR(2) NOT NULL,
			CompanyCode CHAR(4),
			BrandCompanyCode CHAR(4),
			ShopEntryTypeCode CHAR(2) NOT NULL,
			OpenDate CHAR(8) NOT NULL,
			CloseDate CHAR(8),
			SaleReturnBaseType CHAR,
			ChannelCode CHAR(2),
			PastDateSalePermitChk BIT,
			PastDateSalePermitPeriod INT,
			BranchCode VARCHAR(4),
			AreaGroupTypeCode CHAR(2),
			AreaGroupCode CHAR(2),
			AreaCode VARCHAR(10),
			EmpID CHAR(10),
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME,
			CostCenter NVARCHAR(10),
			ProfitCenter NCHAR(10),
			ShopOldNewCode CHAR,
			ComplexShopChk BIT DEFAULT 0 NOT NULL,
			ChiefShopChk BIT,
			BusinessCode NVARCHAR(4),
			ShopTypeCode CHAR(4),
			TelNo VARCHAR(25),
			PlantChk BIT
		);

		INSERT INTO CSL.dbo.Shop 
		(BrandCode, ShopCode, ShopName, VendorGroupCode, VendorGroupShopCode, Floor, Address, PosUseChk, ShopStatusCode, CompanyCode, BrandCompanyCode, ShopEntryTypeCode, OpenDate, CloseDate, SaleReturnBaseType, ChannelCode, PastDateSalePermitChk, PastDateSalePermitPeriod, BranchCode, AreaGroupTypeCode, AreaGroupCode, AreaCode, EmpID, InUserID, InDateTime, CostCenter, ProfitCenter, ShopOldNewCode, ComplexShopChk, ChiefShopChk, BusinessCode, ShopTypeCode, TelNo, PlantChk) 
		VALUES 
		('Q3', 'CFGH', N'Q3 天山店', 'R01', 'RJV0001', 1, N'长宁区天山路789号一层至六层、889号地下一层至', 1, '02', 'F201', 'F201', '21', '20151216', '99991231', '1', '12', 0, 0, 'B411', '11', 'AA', '0000000540', '3000000629', 'system', '2019-08-14 04:30:06.227', 'Q301SCFGH', 'Q301      ', 'A', 0, null, 'F201', '01  ', '13816549820', null),
		('Q3', 'CG08', N'Q3 杭州银泰湖滨店（SPM）', '124', '1845', -1, N'杭州市 杭州市上城区延安路258号', 1, '02', 'F253', 'F201', '70', '20160608', '99991231', '1', '27', 0, 0, 'B428', '11', 'AB', '0000000849', '3000002758', 'system', '2019-08-14 04:30:06.227', 'Q301SCG08', 'Q301SCG08 ', 'A', 0, null, '032C', '02  ', '0571-85870494', null),
		('Q3', 'CH0D', N'Q3 上海七宝万科广场', '786', '0003', 1, N'上海市闵行区漕宝路3366号万科广场L118-21', 1, '02', 'F253', 'F201', '70', '20170628', '99991231', '1', '27', 0, 0, 'B411', '11', 'AA', '0000000540', '3000009485', 'system', '2019-08-14 05:18:10.477', 'Q301SCH0D', 'Q301SCH0D ', 'A', 1, 0, 'F253', '02  ', '021-63895393', null),
		('Q3', 'CH6W', N'Q3 成都衣恋纽可商贸 成都衣恋纽可商贸SHOOPEN', 'R01', 'R000001', 1, N'成都市 成华区建设路53号B1至6层、55号1层至', 1, '02', 'F201', 'F201', '21', '20170901', '99991231', '1', '12', 0, 0, 'B413', '11', 'AP', '0000001712', '3000038371', 'system', '2019-08-14 04:30:06.227', 'Q301SCH6W', 'Q301      ', 'A', 0, null, 'F201', '01  ', '028-61398666', null),
		('Q3', 'CH6X', N'Q3 辽宁萃兮纽可尔商贸', 'R01', 'RJV0009', 1, N'辽宁省 沈阳市和平区市府大路188号', 1, '02', 'F201', 'F201', '21', '20170902', '99991231', '1', '12', 0, 0, 'B418', '11', 'AN', '0000000346', '3000038371', 'system', '2019-08-14 04:30:06.227', 'Q301SCH6X', 'Q301      ', 'A', 0, null, 'F201', '01  ', '024-82517297', null),
		('Q3', 'CH71', N'Q3 上海松江万达SPM', '466', '1881', 1, N'上海市 松江区茸梅路518号1幢964室', 1, '02', 'F253', 'F201', '70', '20170909', '99991231', '1', '27', 0, 0, 'B411', '11', 'AA', '0000000540', '3000009485', 'system', '2019-08-14 05:18:10.477', 'Q301SCH71', 'Q301SCH71 ', 'A', 1, 0, 'F253', '02  ', '67885917', null),
		('Q3', 'CH7A', N'Q3 上海宝乐汇', '960', '0001', 2, N'上海市宝山区宝杨路699号 上海市', 1, '02', 'F261', 'F201', '70', '20180909', '99991231', '1', '27', 0, 0, 'B411', '11', 'AA', '0000000540', '3000009485', 'system', '2019-08-14 05:18:10.477', 'Q301SCH7A', 'Q301SCH7A ', 'A', 1, 0, 'F261', '02  ', '021-56935650', null),
		('Q3', 'CH7Q', N'Q3 杭州湖滨银泰', '124', '1846', -1, N'杭州市 杭州市上城区延安路258号', 1, '02', 'F253', 'F201', '70', '20170901', '99991231', '1', '27', 0, 0, 'B428', '11', 'AB', '0000000849', '3000038371', 'system', '2019-08-14 04:30:06.227', 'Q301SCH7Q', 'Q301SCH7Q ', 'A', 0, null, '032C', '02  ', '0571-85870494', null),
		('Q3', 'CHGA', N'Q3 扬州京华城 SHOOPEN (復合)', '353', '1409', 1, N'扬州市 文昌西路南侧(京华城中城)', 1, '02', 'F261', 'F201', '70', '20171125', '99991231', '1', '27', 0, 0, 'B416', '11', 'AC', '0000000612', '3000038371', 'system', '2019-08-14 05:18:10.477', 'Q301SCHGA', 'Q301SCHGA ', 'A', 1, 0, 'D261', '02  ', '18616346162', null),
		('Q3', 'CHTM', N'Q3 哈尔滨红博会展花园店SPM', '348', '1633', 1, N'黑龙江省哈尔滨市 南岗区红旗大街339号', 1, '02', 'F251', 'F201', '70', '20180419', '99991231', '1', '27', 0, 0, 'B418', '11', 'AM', '0000000455', '3000023950', 'system', '2019-08-14 05:18:10.477', 'Q301SCHTM', 'Q301SCHTM ', 'A', 1, 0, '044A', '02  ', '82273300', null),
		('Q3', 'CHZM', N'Q3 广州富力海珠城', '913', '0001', 2, N'广州市 海珠区前进路安和巷26号之一首层', 1, '02', 'F256', 'F201', '70', '20180601', '99991231', '1', '27', 0, 0, 'B415', '11', 'AY', '0000001397', '3000038090', 'system', '2019-08-14 05:18:10.477', 'Q301SCHZM', 'Q301SCHZM ', 'A', 1, 0, '020F', '02  ', '020-89446045', null),
		('Q3', 'CJ3B', N'Q3 无锡T12时尚购物中心  (復合)', '841', '0002', 2, N'无锡市 中山路328号', 1, '02', 'F261', 'F201', '70', '20180711', '99991231', '1', '27', 0, 0, 'B416', '11', 'AC', '0000000572', '3000007235', 'system', '2019-08-14 05:18:10.477', 'Q301SCJ3B', 'Q301SCJ3B ', 'A', 1, 0, 'A261', '02  ', '0510-82739042', null),
		('Q3', 'CJ5K', N'Q3 杭州湖滨银泰AY', '124', '1859', -1, N'杭州市 杭州市上城区延安路258号', 1, '02', 'F253', 'F201', '70', '20180810', '99991231', '1', '27', 0, 0, 'B428', '11', 'AB', '0000000849', '3000038371', 'system', '2019-08-14 04:30:06.227', 'Q301SCJ5K', 'Q301SCJ5K ', 'A', 0, null, '045C', '02  ', '0571-85870494', null),
		('Q3', 'CJC1', N'Q3 成都龙湖金楠天街', '739', '0007', 2, N'成都市武侯区晋吉北路龙湖金楠天街购物中心1/2楼', 1, '02', 'F257', 'F201', '70', '20181019', '99991231', '1', '27', 0, 0, 'B413', '11', 'AP', '0000001712', '3000010229', 'system', '2019-08-14 05:18:10.477', 'Q301SCJC1', 'Q301SCJC1 ', 'A', 1, 0, 'F257', '02  ', '028-87775364', null),
		('Q3', 'CJDF', N'Q3 南京友谊商场', '957', '0002', 2, N'南京市 秦淮区汉中路27号', 1, '02', 'F261', 'F201', '70', '20190327', '99991231', '1', '27', 0, 0, 'B416', '11', 'AC', '0000000553', '3000008815', 'system', '2019-08-14 05:18:10.477', 'Q301SCJDF', 'Q301SCJDF ', 'A', 1, 0, 'B261', '02  ', '025-51806801', null),
		('Q3', 'CJDM', N'Q3 长沙德思勤（SPM)', '869', '0003', 2, N'长沙市 雨花区湘府中路18号德思勤城市广场卫视中心', 1, '02', 'F256', 'F201', '70', '20181116', '99991231', '1', '27', 0, 0, 'B420', '11', 'BB', '0000001082', '3000004793', 'system', '2019-08-14 05:18:10.477', 'Q301SCJDM', 'Q301SCJDM ', 'A', 1, 0, '018F', '02  ', '0731-89739166', null),
		('Q3', 'CJFV', N'Q3 武汉汉街万达', '466', '1793', 3, N'武汉 武昌区水果湖横路3号', 1, '02', 'F256', 'F201', '70', '20181203', '99991231', '1', '27', 0, 0, 'B420', '11', 'BE', '0000001186', '3000001112', 'system', '2019-08-14 05:18:10.477', 'Q301SCJFV', 'Q301SCJFV ', 'A', 1, 0, 'I256', '03  ', '027-87835318', null),
		('Q3', 'CJJM', N'Q3 上海三林印象城', 'A37', '0002', -1, N'上海市浦东新区环林东路799弄 65号1层', 1, '02', 'F261', 'F201', '70', '20190125', '99991231', '1', '27', 0, 0, 'B411', '11', 'AA', '0000000540', '3000009485', 'system', '2019-08-14 05:18:10.477', 'Q301SCJJM', 'Q301SCJJM ', 'A', 1, 0, 'F261', '02  ', '021-000000', null),
		('Q3', 'CJJT', N'Q3 上海青浦万达广场', '466', '1915', 1, N'上海市青浦区沪青平公路5251号 二楼C区188室', 1, '02', 'F261', 'F201', '70', '20190620', '99991231', '1', '27', 0, 0, 'B411', '11', 'AA', '0000000540', '3000009485', 'system', '2019-08-14 05:18:10.477', 'Q301SCJJT', 'Q301SCJJT ', 'A', 1, 0, 'F261', '02  ', '021-000000', null),
		('Q3', 'CJRH', N'Q3 温州万象城', '498', '0006', 1, N'温州市 温州市瓯海区横港头安昌锦园6幢105室', 1, '02', 'F261', 'F201', '70', '20190426', '99991231', '1', '27', 0, 0, 'B428', '11', 'AB', '0000000914', '3000007235', 'system', '2019-08-14 05:18:10.477', 'Q301SCJRH', 'Q301SCJRH ', 'A', 1, 0, 'C261', '02  ', '0577-85700809', null),
		('Q3', 'CJXB', N'Q3 上海百联东方商厦杨浦店', '111', '1845', -1, N'上海市杨浦区 四平路2500号', 1, '02', 'F201', 'F201', '21', '20190710', '99991231', '1', '12', 0, 0, 'B411', '11', 'AA', '0000000540', '3000009485', 'system', '2019-08-14 05:18:10.477', 'Q301SCJXB', 'Q301      ', 'B', 1, 0, 'F201', '01  ', '11111111', null),
		('SA', 'CCX4', N'SA 武汉汉街万达', '466', '1793', 3, N'武汉 武昌区水果湖横路3号', 1, '02', 'F256', 'F201', '70', '20140912', '99991231', '1', '27', 0, 0, 'B420', '11', 'BE', '0000001186', '3000001112', 'system', '2019-08-14 05:18:10.550', 'SA01ZCCX4', 'SA01SCCX4 ', 'A', 1, 1, 'I256', '02  ', '027-87835318', null),
		('SA', 'CDGT', N'SA 无锡T12时尚购物中心', '841', '0002', 2, N'无锡市 中山路328号', 1, '02', 'F261', 'F201', '70', '20141230', '99991231', '1', '27', 0, 0, 'B416', '11', 'AC', '0000000572', '3000007235', 'system', '2019-08-14 05:18:10.550', 'SA01ZCDGT', 'SA01SCDGT ', 'A', 1, 1, 'A261', '04  ', '0510-82739042', null),
		('SA', 'CEGP', N'SA 成都龙湖金楠天街', '739', '0007', 2, N'成都市武侯区晋吉北路龙湖金楠天街购物中心1/2楼', 1, '02', 'F257', 'F201', '70', '20151021', '99991231', '1', '27', 0, 0, 'B413', '11', 'AP', '0000001712', '3000010229', 'system', '2019-08-14 05:18:10.550', 'SA01ZCEGP', 'SA01SCEGP ', 'A', 1, 1, 'F257', '02  ', '028-87775364', null),
		('SA', 'CEKG', N'SA 长沙德思勤（SPM)', '869', '0003', 2, N'长沙市 雨花区湘府中路18号德思勤城市广场卫视中心', 1, '02', 'F256', 'F201', '70', '20150826', '99991231', '1', '27', 0, 0, 'B420', '11', 'BB', '0000001082', '3000004793', 'system', '2019-08-14 05:18:10.550', 'SA01ZCEKG', 'SA01SCEKG ', 'A', 1, 1, '018F', '02  ', '0731-89739166', null),
		('SA', 'CEYL', N'SA 广州富力海珠城', '913', '0001', 2, N'广州市 海珠区前进路安和巷26号之一首层', 1, '02', 'F256', 'F201', '70', '20151019', '99991231', '1', '27', 0, 0, 'B415', '11', 'AY', '0000001397', '3000038090', 'system', '2019-08-14 05:18:10.550', 'SA01ZCEYL', 'SA01SCEYL ', 'A', 1, 1, '020F', '02  ', '020-89446045', null),
		('SA', 'CFGY', N'SA 哈尔滨红博会展花园店SPM', '348', '1633', 1, N'黑龙江省哈尔滨市 南岗区红旗大街339号', 1, '02', 'F251', 'F201', '70', '20151222', '99991231', '1', '27', 0, 0, 'B418', '11', 'AM', '0000000455', '3000023950', 'system', '2019-08-14 05:18:10.550', 'SA01ZCFGY', 'SA01SCFGY ', 'A', 1, 1, '044A', '02  ', '82273300', null),
		('SA', 'CFRW', N'SA 上海宝乐汇', '960', '0001', 2, N'上海市宝山区宝杨路699号 上海市', 1, '02', 'F261', 'F201', '70', '20160427', '99991231', '1', '27', 0, 0, 'B411', '11', 'AA', '0000000540', '3000009485', 'system', '2019-08-14 05:18:10.550', 'SA01ZCFRW', 'SA01SCFRW ', 'A', 1, 1, 'F261', '02  ', '021-56935650', null),
		('SA', 'CFTN', N'SA 南京友谊广场', '957', '0002', 2, N'南京市 秦淮区汉中路27号', 1, '02', 'F261', 'F201', '70', '20160430', '99991231', '1', '27', 0, 0, 'B416', '11', 'AC', '0000000553', '3000008815', 'system', '2019-08-14 05:18:10.550', 'SA01ZCFTN', 'SA01SCFTN ', 'A', 1, 1, 'B261', '02  ', '025-51806801', null),
		('SA', 'CFW5', N'SA 温州万象城', '498', '0006', 1, N'温州市 温州市瓯海区横港头安昌锦园6幢105室', 1, '02', 'F261', 'F201', '70', '20160429', '99991231', '1', '27', 0, 0, 'B428', '11', 'AB', '0000000914', '3000007235', 'system', '2019-08-14 05:18:10.550', 'SA01ZCFW5', 'SA01SCFW5 ', 'A', 1, 1, 'C261', '02  ', '0577-85700809', null),
		('SA', 'CFZV', N'SA 上海松江万达SPM', '466', '1881', 1, N'上海市 松江区茸梅路518号1幢964室', 1, '02', 'F253', 'F201', '70', '20160520', '99991231', '1', '27', 0, 0, 'B411', '11', 'AA', '0000000540', '3000009485', 'system', '2019-08-14 05:18:10.550', 'SA01ZCFZV', 'SA01SCFZV ', 'A', 1, 1, 'F253', '02  ', '67885917', null),
		('SA', 'CG4R', N'SA 上海七宝万科广场', '786', '0003', 1, N'上海市闵行区漕宝路3366号万科广场L118-21', 1, '02', 'F253', 'F201', '70', '20161030', '99991231', '1', '27', 0, 0, 'B411', '11', 'AA', '0000000540', '3000009485', 'system', '2019-08-14 05:18:10.550', 'SA01ZCG4R', 'SA01SCG4R ', 'A', 1, 1, 'F253', '02  ', '021-63895393', null),
		('SA', 'CG4U', N'SA 扬州京华城', '353', '1409', 1, N'扬州市 文昌西路南侧(京华城中城)', 1, '02', 'F261', 'F201', '70', '20160923', '99991231', '1', '27', 0, 0, 'B416', '11', 'AC', '0000000612', '3000009328', 'system', '2019-08-14 05:18:10.550', 'SA01ZCG4U', 'SA01SCG4U ', 'A', 1, 1, 'D261', '02  ', '0514-82087918', null),
		('SA', 'CJ2F', N'SA 成都优客城市奥莱 (特卖)', 'R01', 'RJV0003', 4, 'N成都市 成华区建设路53号B1至6层、55号1层至', 1, '02', 'F201', 'F201', '21', '20180629', '99991231', '1', '12', 0, 0, 'B413', '11', 'AP', '0000001712', '3000037894', 'system', '2019-08-14 04:30:06.227', 'SA01SCJ2F', 'SA01      ', 'B', 0, null, 'F201', '04  ', '028-61398666', null),
		('SA', 'CJJD', N'SA 上海百联东方商厦杨浦店', '111', '1845', 1, N'上海市杨浦区 四平路2500号', 1, '02', 'F201', 'F201', '21', '20190104', '99991231', '1', '12', 0, 0, 'B411', '11', 'AA', '0000000540', '3000009485', 'system', '2019-08-14 05:18:10.550', 'SA01ZCJJD', 'SA01      ', 'A', 1, 1, 'F201', '01  ', '11111111', null),
		('SA', 'CJJL', N'SA 上海三林印象城', 'A37', '0002', -1, N'上海市浦东新区环林东路799弄 65号1层', 1, '02', 'F261', 'F201', '70', '20190124', '99991231', '1', '27', 0, 0, 'B411', '11', 'AA', '0000000540', '3000009485', 'system', '2019-08-14 05:18:10.550', 'SA01ZCJJL', 'SA01SCJJL ', 'A', 1, 1, 'F261', '02  ', '021-000000', null),
		('SA', 'CJJS', N'SA 上海青浦万达广场', '466', '1915', 1, N'上海市青浦区沪青平公路5251号 二楼C区188室', 1, '02', 'F261', 'F201', '70', '20190620', '99991231', '1', '27', 0, 0, 'B411', '11', 'AA', '0000000540', '3000009485', 'system', '2019-08-14 05:18:10.550', 'SA01ZCJJS', 'SA01SCJJS ', 'A', 1, 1, 'F261', '02  ', '021-000000', null),
		('SA', 'CJUT', N'SA 长沙通程华晨奥特莱斯', '252', '1485', 1, N'长沙市 劳动西路589号', 1, '02', 'F201', 'F201', '21', '20190626', '99991231', '1', '12', 0, 0, 'B420', '11', 'BB', '0000001082', '3000004793', 'system', '2019-08-14 04:30:06.227', 'SA01SCJUT', 'SA01      ', 'B', 0, null, 'F201', '04  ', '0731-85551143', null);
	`
	if _, err := session.Exec(sql); err != nil {
		fmt.Printf("createShopTable error: %v", err.Error())
		fmt.Println()
	}
}

func createTransferMasterSP() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createTransferMasterSP error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP PROCEDURE IF EXISTS dbo.up_CSLK_IOM_InsertRotationOuterReg_RecvSuppMst_C1_Clearance"); err != nil {
		log.Printf("createTransferMasterSP error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE PROCEDURE [dbo].[up_CSLK_IOM_InsertRotationOuterReg_RecvSuppMst_C1_Clearance]
			@BrandCode VARCHAR(4) = NULL,
			@ShopCode CHAR(4) = NULL,
			@TargetShopCode CHAR(4) = NULL,       -- 상대매장
			@OutDate CHAR(8) = NULL,    -- 出库日期
			@WaybillNo VARCHAR(13) = NULL,
			@BoxNo CHAR(20) = NULL,
			@EmpID	CHAR(10) = NULL,	 -- 등록자
			@ShippingCompanyCode CHAR(2) = NULL,
			@DeliveryOrderNo VARCHAR(250) = NULL, -- 快递凭证编号

			@IsBigSize INT,
			@ExpressNo VARCHAR(13) = NULL /*add by yuan.yujiao 20130617*/,
			@BoxAmount VARCHAR(8) = NULL, -- 箱数
			@StockOutUseAmt VARCHAR(16) = NULL,   -- 快递费
			@ProvinceCode VARCHAR(8) = NULL,      --moidfy by li.guolin 20170823
			@CityCode VARCHAR(8) = NULL,          --moidfy by li.guolin 20170823
			@DistrictCode CHAR(8) = NULL,         --moidfy by li.guolin 20170823
			@Area NVARCHAR(100) = NULL,           --moidfy by li.guolin 20170823
			@ShopManagerName NVARCHAR(10) = NULL, --moidfy by li.guolin 20170823
			@MobilePhone VARCHAR(25) = NULL       --moidfy by li.guolin 20170823
		AS
			--SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED
			SET XACT_ABORT ON;
			SET NOCOUNT ON;
			BEGIN

				DECLARE @NewRecvSuppNo CHAR(14);
				DECLARE @NewSeq INT;
				DECLARE @SendState VARCHAR(2);
				DECLARE @SendFlag CHAR(2);
				DECLARE @ShippingTypeCode CHAR(2); -- 운송구분코드
				DECLARE @RecvSuppType CHAR(1); -- 입출고구분
				DECLARE @PlantCode CHAR(4); -- 물류센터코드
				DECLARE @NormalProductType CHAR(1); -- 정품구분 A:정품 B:비품
				DECLARE @SAPMenuType CHAR(1);
				DECLARE @EmpName NVARCHAR(200);
				DECLARE @UserID VARCHAR(20);
				DECLARE @TransTypeCode CHAR(1) = NULL;   -- 운송구분
				DECLARE @RecvSuppStatusCode CHAR(1) = NULL;
				DECLARE @WideBranchRotation BIT = 0;


				DECLARE @ErrorCode NVARCHAR(1000) = '';
				DECLARE @ErrorParam1 NVARCHAR(4000) = '';
				DECLARE @ErrorParam2 NVARCHAR(4000) = '';

				BEGIN TRY

					--同一广域支社不需批准 add by zhai.weihao 20151027
					SELECT @WideBranchRotation = CASE
													WHEN
					(
						SELECT B.WideBranchCode
						FROM Shop AS A WITH (NOLOCK)
							JOIN Branch AS B WITH (NOLOCK)
								ON A.BranchCode = B.BranchCode
						WHERE A.ShopCode = @ShopCode
							AND A.BrandCode = @BrandCode
					)   =
					(
						SELECT B.WideBranchCode
						FROM Shop AS A WITH (NOLOCK)
							JOIN Branch AS B WITH (NOLOCK)
								ON A.BranchCode = B.BranchCode
						WHERE A.ShopCode = @TargetShopCode
							AND A.BrandCode = @BrandCode
					)   THEN
														0
													ELSE
														1
												END;
					--同一广域支社不需批准 add by zhai.weihao 20151027
					IF (
						LEFT(dbo.udf_CSLK_MonthlyClosingChk('03', 'Zn'), 1) = 1
					)
					BEGIN
						SELECT @ErrorCode
							= SUBSTRING(
								dbo.udf_CSLK_MonthlyClosingChk('03', 'Zn'), 2, 510
									);
						EXEC [up_CSLK_ComonRaiseError] @ErrorCode, @ErrorParam1, @ErrorParam2;
					END;



					-- 기본값 세팅
					SET @ShippingTypeCode = '20'; -- 회전(매장 직->직)
					SET @RecvSuppType = 'S'; -- 입출고구분 S:출고/ R:입고
					SET @TransTypeCode = '5'; -- 운송타입(5:운송사배송)
					IF @WideBranchRotation = 0 --add by song.hejia 20121031
						SET @RecvSuppStatusCode = 'R'; -- 출고
					ELSE
						SET @RecvSuppStatusCode = 'W'; --add by song.hejia 20121031 等待批准
					SET @NormalProductType = 'A';

					IF @BoxAmount = ''
						SET @BoxAmount = NULL;

					IF @StockOutUseAmt = ''
						SET @StockOutUseAmt = NULL;


					--날짜
					IF @OutDate IS NULL OR @OutDate = ''
						SET @OutDate = CONVERT(CHAR(8), GETDATE(), 112);

					IF @ShippingCompanyCode = ''
						SET @ShippingCompanyCode = NULL;

					SELECT @NewSeq = ISNULL(MAX(SeqNo), 0) + 1
					FROM RecvSuppMst
					WHERE BrandCode = @BrandCode
						AND ShopCode = @ShopCode
						AND Dates = @OutDate
						AND SeqNo < 6000;

					--ShopCode(4)+yymmdd+9999(14자리)
					SET @NewRecvSuppNo
						= @ShopCode + RIGHT(@OutDate, 6)
						+ RIGHT(REPLICATE('0', 4) + CONVERT(VARCHAR, @NewSeq), 4);


					--물류창고
					SELECT @PlantCode = PlantCode
					FROM Brand WITH (NOLOCK)
					WHERE BrandCode = @BrandCode;

					-- SAP관련세팅
					IF @WideBranchRotation = 0 --add by song.hejia 20121031
						SET @SendFlag = 'R'; -- 등록
					ELSE
						SET @SendFlag = ''; --批准完成 在SP里用R处理 add by song.hejia 20121031

					SET @SendState = '';
					SET @SAPMenuType = '4'; -- 매장회전

					IF @IsBigSize = 0
					BEGIN
						IF EXISTS
						(
							SELECT 1
							FROM RecvSuppMst
							WHERE BrandCode = @BrandCode
								AND ShopCode = @ShopCode
								AND ShippingCompanyCode=@ShippingCompanyCode
								AND WayBillNo = @WaybillNo
								AND BoxNo = @BoxNo
								--AND Dates = @Dates 注释by wang.wanyue 不再按日期检查
								AND DelChk = 0
						)
						BEGIN
							SET @ErrorCode = 'IOM164';
							EXEC dbo.[up_CSLK_ComonRaiseError] @ErrorCode, @ErrorParam1, @ErrorParam2;
						END;
					END;

				SELECT @EmpName = EmpName
					FROM Employee
					WHERE EmpID = @EmpID


					SELECT @UserID = UserID
					FROM UserInfo
					WHERE EmpID = @EmpID


					-- 입출고 공통 (회전출고등록MASTER)
					INSERT INTO RecvSuppMst
					(   RecvSuppNo,
						BrandCode,
						ShopCode,
						Dates,
						SeqNo,
						RecvSuppType,
						ShopSuppRecvDate,
						TransTypeCode,
						ShippingTypeCode,
						WayBillNo,
						RecvSuppStatusCode,
						NormalProductType,
						BoxNo,
						PlantCode,
						RoundRecvSuppNo,
						TargetShopCode,
						InUserID,
						InDateTime,
						ModiUserID,
						ModiDateTime,
						RecvEmpID,
						RecvEmpName,
						SuppEmpID,    -- 출고직원ID
						SuppEmpName,  -- 출고직원명
						BrandSuppRecvDate,
						SAPMenuType,
						SendState,
						SendFlag,
						InvtBaseDate, -- 재고기준일자
						BoxAmount /*add by song.hejia 20130415*/,
						StockOutUseAmt /*add by song.hejia 20130415*/,
						ExpressNo /*add by yuan.yujiao 20130617*/,
						ShippingCompanyCode,
						DeliveryID,
						DeliveryOrderNo,
						VolumeType,
						VolumesSize,
						VolumesUnit,
						ProvinceCode,
						CityCode,
						DistrictCode,
						Area,
						ShopManagerName,
						MobilePhone,
						BoxType,
						Channel
					)
					VALUES
					(   @NewRecvSuppNo,
						@BrandCode,
						@ShopCode,
						@OutDate,
						@NewSeq,
						@RecvSuppType,
						@OutDate,
						@TransTypeCode,
						@ShippingTypeCode,
						@WaybillNo,
						@RecvSuppStatusCode,
						@NormalProductType,
						@BoxNo,
						@PlantCode,
						'',
						@TargetShopCode,
						@UserID,         -- 등록자
						GETDATE(),
						@UserID,         -- 수정자
						GETDATE(),
						'',                -- RecvEmpID
						'',                -- RecvEmpName
						@EmpID,
						@EmpName,
						'',
						@SAPMenuType,
						@SendState,
						@SendFlag,
						@OutDate, -- 재고기준일자.
						@BoxAmount /*add by song.hejia 20130415*/,
						@StockOutUseAmt /*add by song.hejia 20130415*/,
						@ExpressNo /*add by yuan.yujiao 20130617*/,
						@ShippingCompanyCode,
						@WaybillNo,
						@DeliveryOrderNo,
						N'中箱子',
						'0.09486',
						N'm³',
						@ProvinceCode,
						@CityCode,
						@DistrictCode,
						@Area,
						@ShopManagerName,
						@MobilePhone,
						'MB',
						'Clearance'
					);

				/*
					添加物流信息
				*/
					IF NOT EXISTS(
					SELECT 1
					FROM WayBillNo
					WHERE ShippingCompanyCode=@ShippingCompanyCode
					AND WayBillNo=@WaybillNo
					)
					BEGIN
					INSERT INTO WayBillNo
					(
						ShippingCompanyCode, WayBillNo, AllowDulpChk, InUserID, InDateTime, ModiUserID, ModiDateTime
					)
					VALUES
					(
						@ShippingCompanyCode,@WaybillNo,0,@UserID,GETDATE(),@UserID,GETDATE()
					)
					END

					-- 입출고번호 리턴
					SELECT @NewRecvSuppNo AS RecvSuppNo;

				END TRY
				BEGIN CATCH
					EXEC [up_CSLK_ComonRaiseError] @ErrorCode, @ErrorParam1, @ErrorParam2;
				END CATCH;
			END;
	`

	if _, err := session.Exec(sql); err != nil {
		log.Printf("createTransferMasterSP error: %v", err.Error())
		log.Println()
	}
}

func createTransferDetailSP() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createTransferDetailSP error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP PROCEDURE IF EXISTS dbo.up_CSLK_IOM_InsertRotationOuterReg_RecvSuppDtl_C1_Clearance"); err != nil {
		log.Printf("createTransferDetailSP error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE PROCEDURE [dbo].[up_CSLK_IOM_InsertRotationOuterReg_RecvSuppDtl_C1_Clearance]
			@RecvSuppNo    CHAR(14) = NULL
			,@RecvSuppSeqNo    INT   = NULL
			,@BrandCode     VARCHAR(4) = NULL
			,@ShopCode     CHAR(4) = NULL
			,@OutDate      CHAR(8) = NULL -- 出库日期
			,@ProdCode     VARCHAR(18) = NULL
			,@RecvSuppQty    INT = NULL
			,@EmpID	CHAR(10) = NULL	 -- 등록자
			
			AS
			SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED
			SET XACT_ABORT ON
			SET NOCOUNT ON
			BEGIN
			DECLARE @SendState  varchar(2)
			DECLARE @SendFlag  char(1)
			
			BEGIN TRY
			
			DECLARE @ErrorCode  NVARCHAR(50) ='';
			DECLARE @ErrorParam1 NVARCHAR(4000)='';
			DECLARE @ErrorParam2 NVARCHAR(4000)='';
			DECLARE @StyleCode VARCHAR(18);
			DECLARE @WideBranchRotation BIT = 0;
			DECLARE @UserID     VARCHAR(20) = NULL;
			DECLARE @SeqNo			 INT = NULL;
			DECLARE @SalePrice	DECIMAL(19,2) = NULL;
			
			
			IF (@RecvSuppSeqNo IS NULL OR @RecvSuppSeqNo = 0)
			BEGIN
				IF( EXISTS (select 1 from RecvSuppDtl where RecvSuppNo=@RecvSuppNo))
				BEGIN
					SET @RecvSuppSeqNo=(select MAX(RecvSuppSeqNo)+1 from RecvSuppDtl where RecvSuppNo=@RecvSuppNo)
					SET @SeqNo=@RecvSuppSeqNo
				END
				ELSE
				BEGIN
					SET @RecvSuppSeqNo=1
					SET @SeqNo=1
				END
			
			END
			
			
			IF @RecvSuppSeqNo>400
			BEGIN
				SET @ErrorCode = 'IOM189'
				EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2
			END
			
			SELECT @Stylecode = StyleCode
			FROM Product WHERE ProdCode = @ProdCode
			
			SET @SendState = ''
		
			IF @OutDate = NULL OR @OutDate = ''
			SET @OutDate = CONVERT(CHAR(8),GETDATE(),112)
		
			--同一广域支社不需批准 add by zhai.weihao 20151027
			SELECT @WideBranchRotation =
			CASE WHEN(
			SELECT B.WideBranchCode
			FROM Shop AS A  WITH(NOLOCK)
				JOIN Branch AS B WITH(NOLOCK)
					ON A.BranchCode =B.BranchCode
			WHERE A.ShopCode =@ShopCode AND A.BrandCode=@BrandCode)=(SELECT B.WideBranchCode
			FROM Shop AS A  WITH(NOLOCK)
				JOIN Branch AS B WITH(NOLOCK)
					ON A.BranchCode =B.BranchCode
			WHERE A.ShopCode =(SELECT TargetShopCode  FROM RecvSuppMst WHERE RecvSuppNo =@RecvSuppNo ) AND A.BrandCode=@BrandCode)THEN 0 ELSE 1 END
			--同一广域支社不需批准 add by zhai.weihao 20151027
		
			IF @WideBranchRotation = 0  --add by song.hejia 20121031
				SET @SendFlag = 'R'
			ELSE
				SET @SendFlag='' --add by song.hejia 20121031
				SET @RecvSuppSeqNo=@SeqNo
			
			
				SELECT @UserID = UserID
				FROM UserInfo
				WHERE EmpID = @EmpID
			
				SELECT @SalePrice = Price
					FROM Product
					WHERE BrandCode = @BrandCode
					AND ProdCode = @ProdCode
			
				-- 입출고 상세 (회전출고상세)
				INSERT INTO RecvSuppDtl
				(
				RecvSuppNo
				,RecvSuppSeqNo
				,BrandCode
				,ShopCode
				,Dates
				,SeqNo
				,RoundRecvSuppNo
				,RoundRecvSuppDtSeq
				,ProdCode
				,RecvSuppQty
				,RecvSuppFixedQty
				,SalePrice
				,InUserID
				,InDateTime
				,ModiUserID
				,ModiDateTime
				,SendState
				,SendFlag
				)
				VALUES
				(
				@RecvSuppNo
				,@RecvSuppSeqNo
				,@BrandCode
				,@ShopCode
				,@OutDate
				,@SeqNo
				,null --@RoundRecvSuppNo
				,null --@RoundRecvSuppDtSeq
				,(SELECT Upper(@ProdCode)) --modify by zhangxia 20141126
				,@RecvSuppQty
				,@RecvSuppQty  -- 출고확정수량
				,@SalePrice
				,@UserID
				,GETDATE()
				,@UserID
				,GETDATE()
				,@SendState
				,@SendFlag
				)
			
			
			END TRY
			BEGIN CATCH
			
				EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2
			END CATCH
		END
	`

	if _, err := session.Exec(sql); err != nil {
		log.Printf("createTransferDetailSP error: %v", err.Error())
		log.Println()
	}
}

func createTransferConfirmMasterSP() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createTransferConfirmMasterSP error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP PROCEDURE IF EXISTS dbo.up_CSLK_IOM_InsertRotationEnterConfirm_RecvSuppMst_C1_Clearance"); err != nil {
		log.Printf("createTransferConfirmMasterSP error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE PROCEDURE [dbo].[up_CSLK_IOM_InsertRotationEnterConfirm_RecvSuppMst_C1_Clearance]
			@BrandCode    varchar(4) = null
			,@ShopCode    char(4) = null -- 入库卖场
			,@TargetShopCode  char(4) = null   -- 出库卖场
			,@InDate    char(8)   -- 入库日期
			,@WayBillNo    varchar(13) = null
			,@BoxNo     char(20) = null
			,@RoundRecvSuppNo  char(14) = null   -- 出库单RecvSuppNo
			,@EmpID    char(10) = null
		
		AS
		--SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED
		SET XACT_ABORT ON
		SET NOCOUNT ON
		BEGIN
		
			DECLARE @NewRecvSuppNo  char(14)
			DECLARE @NewSeq    int
			DECLARE @SendState   varchar(2)
			DECLARE @SendFlag   char(2)
			DECLARE @ShippingTypeCode char(2)    -- 운송구분코드
			DECLARE @RecvSuppType  char(1)    -- 입출고구분
			DECLARE @PlantCode   char(4)    -- 물류센터코드
			DECLARE @RequestNo   Varchar(20)
			DECLARE @SAPMenuType  char(1)
			DECLARE @EmpName  nvarchar(200)
			DECLARE @NormalProductType char(1)
			DECLARE @SuppEmpID   CHAR(10)
			DECLARE @SuppEmpName  nvarchar(200)
			DECLARE @TransTypeCode   CHAR(1)   -- 운송구분
			DECLARE @RecvSuppStatusCode CHAR(1)
			DECLARE @UserID    varchar(20)
		
			DECLARE @ErrorCode  NVARCHAR(1000) ='';
			DECLARE @ErrorParam1 NVARCHAR(4000)='';
			DECLARE @ErrorParam2 NVARCHAR(4000)='';
		
			BEGIN TRY
			-- 마감체크
			IF (left(dbo.udf_CSLK_MonthlyClosingChk('04','Zn'),1) = 1)
			BEGIN
			SELECT @ErrorCode = SUBSTRING(dbo.udf_CSLK_MonthlyClosingChk('04','Zn'),2,510)
			EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2
			END
		
		
			DECLARE @TempRecvSuppStatusCode CHAR(1)
			SELECT @TempRecvSuppStatusCode = RecvSuppStatusCode FROM RecvSuppMst WITH(NOLOCK) WHERE RecvSuppNo = @RoundRecvSuppNo
		
			IF @TempRecvSuppStatusCode = 'F'
			BEGIN
			SET @ErrorCode = 'COM011'
			SET @ErrorParam1 = ''
			EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2
			END
			---------------------------------------------------------------------
		
			--날짜
			IF @InDate IS NULL OR @InDate = ''
			SET @InDate =  CONVERT(CHAR(8),GETDATE(),112)
		
			SELECT @NewSeq = IsNull(MAX(SeqNo),0) + 1
			FROM RecvSuppMst
			WHERE BrandCode = @BrandCode
				AND ShopCode = @ShopCode
				AND Dates = @InDate
				AND SeqNo < 6000
		
			--ShopCode(4)+yymmdd+9999(14자리)
			SET @NewRecvSuppNo = @ShopCode + Right(@InDate,6)+ Right(replicate('0',4) + convert(varchar,@NewSeq),4)
		
			DECLARE  @ShippingCompanyCode  CHAR(2) --ADD BY GAOMEIFANG 2017.5.12
			,@DeliveryID VARCHAR(250)=NULL --moidfy by li.guolin 20170811
			,@DeliveryOrderNo VARCHAR(250)=NULL --moidfy by li.guolin 20170811
			,@VolumeType NVARCHAR(20)=NULL --moidfy by li.guolin 20170811
			,@VolumesSize VARCHAR(20)=NULL --moidfy by li.guolin 20170811
			,@VolumesUnit NVARCHAR(10)=NULL --moidfy by li.guolin 20170811
			,@ProvinceCode VARCHAR(8)=NULL --moidfy by li.guolin 20170823
			,@CityCode VARCHAR(8)=NULL --moidfy by li.guolin 20170823
			,@DistrictCode CHAR(8)=NULL --moidfy by li.guolin 20170823
			,@Area NVARCHAR(100)=NULL --moidfy by li.guolin 20170823
			,@ShopManagerName NVARCHAR(10)=NULL --moidfy by li.guolin 20170823
			,@MobilePhone VARCHAR(25)=NULL --moidfy by li.guolin 20170823
			,@BoxType CHAR(2)=NULL --moidfy by li.guolin 20170930
			,@ExpressNo  VARCHAR(250)=NULL --moidfy by li.guolin 20170811
			-- ShippingTypeCode 회전 매장간 '20' 본사 '24'
			SELECT @RequestNo = RequestNo
				,@ShippingTypeCode = ShippingTypeCode
				,@NormalProductType = NormalProductType
				,@SuppEmpID = SuppEmpID
				,@SuppEmpName = SuppEmpName
				,@ShippingCompanyCode=ShippingCompanyCode
				,@DeliveryID=DeliveryID --moidfy by li.guolin 20170811
				,@DeliveryOrderNo =DeliveryOrderNo --moidfy by li.guolin 20170811
				,@VolumeType =VolumeType --moidfy by li.guolin 20170811
				,@VolumesSize =VolumesSize --moidfy by li.guolin 20170811
				,@VolumesUnit =VolumesUnit --moidfy by li.guolin 20170811
				,@ProvinceCode =ProvinceCode --moidfy by li.guolin 20170823
				,@CityCode =CityCode --moidfy by li.guolin 20170823
				,@DistrictCode =DistrictCode --moidfy by li.guolin 20170823
				,@Area =Area --moidfy by li.guolin 20170823
				,@ShopManagerName =ShopManagerName --moidfy by li.guolin 20170823
				,@MobilePhone =MobilePhone --moidfy by li.guolin 20170823
				,@BoxType =BoxType --moidfy by li.guolin 20170930
			FROM RecvSuppMst
			WHERE RecvSuppNo = @RoundRecvSuppNo
			AND DelChk=0
		
		IF @ShippingCompanyCode=''
			SET @ShippingCompanyCode=NULL
		
			-- 출하유형 최종체크
			IF @ShippingTypeCode = '' or @ShippingTypeCode is null
			BEGIN
				IF @RequestNo = '' or @ShippingTypeCode is null
				SET @ShippingTypeCode = '20' -- 요청번호 없는경우 매장간
				ELSE
				SET @ShippingTypeCode = '24' -- 요청번호 없는경우 본사
			END
		
			SET @RecvSuppType = 'R'   -- 입출고구분 S:출고/ R:입고
			SET @TransTypeCode = '5'  -- 운송타입
			SET @RecvSuppStatusCode = 'F' -- 입고확정
		
			-- 물류센터코드
			SELECT @PlantCode = PlantCode
			FROM Brand WITH(NOLOCK)
			WHERE BrandCode = @BrandCode
		
			-- SAP관련세팅
			SET @SendFlag = 'R'
			SET @SendState = ''
			SET @SAPMenuType = '4'  -- 매장회전
		
		
			SELECT @EmpName = EmpName
			FROM Employee WITH(NOLOCK)
			WHERE EmpID = @EmpID
		
			SELECT @UserID = UserID
				FROM UserInfo
				WHERE EmpID = @EmpID
		
			-- 입출고 공통 (회전입고등록MASTER)
			INSERT INTO RecvSuppMst
			(
				RecvSuppNo
			,BrandCode
			,ShopCode
			,Dates
			,SeqNo
			,RecvSuppType
			,ShopSuppRecvDate
			,TransTypeCode
			,ShippingTypeCode
			,WayBillNo
			,RecvSuppStatusCode
			,BoxNo
			,PlantCode
			,RoundRecvSuppNo
			,TargetShopCode
			,RequestNo
			,InUserID
			,InDateTime
			,ModiUserID
			,ModiDateTime
			,RecvEmpID
			,RecvEmpName
			,SuppEmpID
			,SuppEmpName
			,BrandSuppRecvDate
			,SendState
			,SendFlag
			,SAPMenuType
			,NormalProductType
			,InvtBaseDate
			,ShippingCompanyCode
			,DeliveryID
			,DeliveryOrderNo
			,VolumeType
			,VolumesSize
			,VolumesUnit
			,ProvinceCode
			,CityCode
			,DistrictCode
			,Area
			,ShopManagerName
			,MobilePhone
			,BoxType
			,ExpressNo
			,Channel
			)
			VALUES
			(
				@NewRecvSuppNo
			,@BrandCode
			,@ShopCode
			,@InDate
			,@NewSeq     -- @SeqNo
			,@RecvSuppType
			,@InDate
			,@TransTypeCode
			,@ShippingTypeCode
			,@WayBillNo
			,@RecvSuppStatusCode
			,@BoxNo
			,@PlantCode
			,@RoundRecvSuppNo
			,@TargetShopCode
			,@RequestNo
			,@UserID     -- 등록자
			,GETDATE()
			,@UserID     -- 수정자
			,GETDATE()
			,@EmpID
			,@EmpName
			,@SuppEmpID       -- SuppEmpID
			,@SuppEmpName       -- SuppEmpName
			,''
			,@SendState
			,@SendFlag
			,@SAPMenuType
			,@NormalProductType
			,@InDate   -- 재고기준일자.
			,@ShippingCompanyCode
			,@DeliveryID
			,@DeliveryOrderNo
			,@VolumeType
			,@VolumesSize
			,@VolumesUnit
			,@ProvinceCode
			,@CityCode
			,@DistrictCode
			,@Area
			,@ShopManagerName
			,@MobilePhone
			,@BoxType
			,@ExpressNo
			,'Clearance'
			)
		
			-- 입고유효성 체크
			-- 출고된 내역이 입고되었는데 또 입고시도할려는 경우
			--modified by weile
			IF (SELECT count(*) FROM RecvSuppMst WITH(NOLOCK) WHERE RoundRecvSuppNo = @RoundRecvSuppNo AND Delchk='0') = 2    --modify by zhai.weihao 20150930 add delchk 判断
			BEGIN
			SET @ErrorCode = 'IOM0304_1'
			SET @ErrorParam1 = ''
			EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2
			END
			-- 회전출고에 회전 입고번호 적용
			UPDATE RecvSuppMst
				SET RoundRecvSuppNo = @NewRecvSuppNo
				,RecvEmpID = @EmpID
				,RecvEmpName = @EmpName
				,RecvSuppStatusCode = 'F'   -- 진행상태(입고확정)
				,SendFlag = @SendFlag
				,InvtBaseDate = @InDate  --  재고기준일자.
				,PlantCode = @PlantCode			--20190104 wang.wanyue 出库卖场和入库卖场PlantCode一致
			WHERE RecvSuppNo = @RoundRecvSuppNo
		
		
			-- 插入Detail
			EXEC [dbo].[up_CSLK_IOM_InsertRotationEnterConfirm_RecvSuppDtl_C1_Clearance]
				@RecvSuppNo    = @NewRecvSuppNo
				,@RecvSuppSeqNo  = NULL
				,@BrandCode     = @BrandCode
				,@ShopCode     = @ShopCode
				,@InDate       = @InDate
				,@RoundRecvSuppNo   = @RoundRecvSuppNo  -- 调货出库单的RecvSuppNo
				,@EmpID   = @EmpID
		
		
			-- 입출고번호 리턴
			SELECT @NewRecvSuppNo AS RecvSuppNo
		
		
			END TRY
			BEGIN CATCH
		
		
			EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2
			END CATCH
		END
	`

	if _, err := session.Exec(sql); err != nil {
		log.Printf("createTransferConfirmMasterSP error: %v", err.Error())
		log.Println()
	}
}

func createTransferConfirmDetailSP() {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE CSL;"); err != nil {
		log.Printf("createTransferMasterSP error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("DROP PROCEDURE IF EXISTS dbo.up_CSLK_IOM_InsertRotationEnterConfirm_RecvSuppDtl_C1_Clearance"); err != nil {
		log.Printf("createTransferMasterSP error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE PROCEDURE [dbo].[up_CSLK_IOM_InsertRotationEnterConfirm_RecvSuppDtl_C1_Clearance]
			@RecvSuppNo    char(14) = null  
			,@RecvSuppSeqNo    int   = null      
			,@BrandCode     varchar(4) = null  
			,@ShopCode     char(4) = null  
			,@InDate      char(8) = null -- 入库日期
			,@RoundRecvSuppNo   char(14) = null  -- 调货出库单的RecvSuppNo
			,@EmpID    char(10) = null
			
		AS     
		SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED        
		SET XACT_ABORT ON     
		SET NOCOUNT ON  
		BEGIN
			DECLARE @SendState  VARCHAR(2)  
			DECLARE @SendFlag  CHAR(1)  
			DECLARE @SENDF   CHAR(1)  
			DECLARE @ErrorCode  NVARCHAR(50) ='';  
			DECLARE @ErrorParam1 NVARCHAR(4000)='';  
			DECLARE @ErrorParam2 NVARCHAR(4000)='';
			DECLARE @UserID    VARCHAR(20)
			
			BEGIN TRY   
		
			SET @SendState = ''  
			SET @SendFlag = 'R'
		
			IF @InDate IS NULL OR @InDate = ''
			SET @InDate = CONVERT(CHAR(8),GETDATE(),112)
			
			-- 인터페이스실시간처리확인    
			EXEC [up_CSLK_IF_CHK_RecvSupp] @p_RECVSUPPNO = @RecvSuppNo, @o_SENDF  = @SENDF OUTPUT  
			
			/* 수신시스템체크  
			R : 등록/수정된상태, 전송전( 수정불가능)  
			I : 전송중( 수정불가능)  
			S : 전송후( 수정가능)  
			*/    
			IF (@SENDF = 'S')  
			--IF (1=1)  
			BEGIN  
			
			-- SAP진행상황체크  
			SELECT @SendFlag = SendFlag  
				FROM RecvSuppDtl   
			WHERE RecvSuppNo = @RecvSuppNo  
				AND RecvSuppSeqNo = @RecvSuppSeqNo  
			
			-- R:등록 I:전송중 S:완료    
			IF ( @SendFlag = 'I')   
			BEGIN   
				
				-- 데이터가 전송중이므로 잠시후 처리하십시오  
				SET @ErrorCode = 'COM100'  
				SET @ErrorParam1 = ''  
				SET @ErrorParam2 = ''       
				EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2  
				
			END  
			ELSE  
			
			BEGIN    
			
			SELECT @UserID = UserID
				FROM UserInfo
				WHERE EmpID = @EmpID
				
			
			-- 입출고 상세 (회전출고상세)  
			INSERT INTO RecvSuppDtl  
			(  
				RecvSuppNo  
				,RecvSuppSeqNo  
				,BrandCode  
				,ShopCode  
				,Dates  
				,SeqNo  
				,RoundRecvSuppNo  
				,RoundRecvSuppDtSeq  
				,ProdCode  
				,RecvSuppQty  
				,RecvSuppFixedQty  
				,SalePrice  
				,InUserID  
				,InDateTime  
				,ModiUserID  
				,ModiDateTime  
				,SendState  
				,SendFlag  
				,DelChk  
			)  
			SELECT   
				@RecvSuppNo
				,RecvSuppSeqNo  
				,@BrandCode  
				,@ShopCode  
				,@InDate
				,SeqNo  
				,RecvSuppNo   -- 회전출고의 입출고번호  
				,RecvSuppSeqNo  -- 회전출고의 입출고상세번호  
				,ProdCode  
				,RecvSuppQty  
				,RecvSuppQty  -- 확정수량 RecvSuppFixedQty    
				,SalePrice   --@SalePrice  
				,@UserID  
				,GETDATE()  
				,@UserID  
				,GETDATE()  
				,@SendState  
				,@SendFlag  
				,DelChk  
			FROM RecvSuppDtl  
			WHERE RecvSuppNo = @RoundRecvSuppNo  
		
		
			-- 회전출고된것 수량 확정  
			UPDATE RecvSuppDtl  
				SET RecvSuppFixedQty = RecvSuppQty  
				,RoundRecvSuppNo = @RecvSuppNo  
				,RoundRecvSuppDtSeq = RecvSuppSeqNo  
				,SendFlag = @SendFlag  
			WHERE RecvSuppNo = @RoundRecvSuppNo     
			
			END  
			END  
			ELSE  
			BEGIN  
			
			-- 본사수정중입니다.  
			SET @ErrorCode = 'IOM132'  
			SET @ErrorParam1 = ''  
			SET @ErrorParam2 = ''       
			EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2  
			
			END  
		
			
			END TRY  
			BEGIN CATCH  
			
			EXEC [up_CSLK_ComonRaiseError] @ErrorCode,@ErrorParam1,@ErrorParam2  
			
			END CATCH   
		END
	`

	if _, err := session.Exec(sql); err != nil {
		log.Printf("createTransferMasterSP error: %v", err.Error())
		log.Println()
	}
}
