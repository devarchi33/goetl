package test

import (
	"clearance-adapter/factory"
	"log"
)

func initProduct() {
	createProdcutDB()
	setProductData()
	setSkuData()
}

func createProdcutDB() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("DROP DATABASE IF EXISTS pangpang_brand_product;"); err != nil {
		log.Printf("createProdcutDB error: %v", err.Error())
		log.Println()
	}

	if _, err := session.Exec("CREATE DATABASE pangpang_brand_product;"); err != nil {
		log.Printf("createProdcutDB error: %v", err.Error())
		log.Println()
	}
}

func setProductData() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_product;"); err != nil {
		log.Printf("setProductData error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE product
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			tenant_code VARCHAR(16),
			code VARCHAR(64),
			name VARCHAR(255),
			brand_id BIGINT(20),
			title_image VARCHAR(255),
			list_price DOUBLE,
			has_digital TINYINT(1),
			enable TINYINT(1),
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setProductData error: %v", err.Error())
		log.Println()
	}
	sql = `
		INSERT INTO pangpang_brand_product.product 
		(tenant_code, code, name, brand_id, title_image, list_price, has_digital, enable, created_at, updated_at, deleted_at) 
		VALUES 
		('pangpang', 'SPYC949S11', '女式格子衬衫', 2, '', 259, 0, 1, '2019-07-25 00:37:52', '2019-07-25 00:37:52', null),
		('pangpang', 'SPWJ948S22', '女式牛仔裙', 2, '', 259, 0, 1, '2019-07-25 00:37:59', '2019-07-25 00:37:59', null),
		('pangpang', 'SPWJ948S23', '女式牛仔裙', 2, '', 199, 0, 1, '2019-07-25 00:37:59', '2019-07-25 00:37:59', null),
		('pangpang', 'SPYC949H21', '大格纹衬衫', 2, '', 259, 0, 1, '2019-07-25 00:37:59', '2019-07-25 00:37:59', null),
		('pangpang', 'SPYS949H22', '条纹衬衫', 2, '', 259, 0, 1, '2019-07-25 00:38:00', '2019-07-25 00:38:00', null),
		('pangpang', 'Q3AFAFDU6S21', '男女同款拖鞋', 40, '', 79, 0, 1, '2019-07-25 00:38:00', '2019-07-25 00:38:00', null);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setProductData error: %v", err.Error())
		log.Println()
	}
}

func setSkuData() {
	session := factory.GetP2BrandEngine().NewSession()
	defer session.Close()

	if _, err := session.Exec("USE pangpang_brand_product;"); err != nil {
		log.Printf("setSkuData error: %v", err.Error())
		log.Println()
	}

	sql := `
		CREATE TABLE sku
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			product_id BIGINT(20),
			code VARCHAR(64),
			name VARCHAR(255),
			image VARCHAR(255),
			enable TINYINT(1),
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setSkuData error: %v", err.Error())
		log.Println()
	}

	sql = `
		INSERT INTO pangpang_brand_product.sku (product_id, code, name, image, created_at, updated_at, deleted_at, enable) 
		VALUES 
			(1, 'SPYC949S1139085', '女式格子衬衫, (39)Ivory, 160/84A(S)', '', '2019-07-25 00:39:04', '2019-07-25 00:39:04', null, 1),
			(1, 'SPYC949S1139090', '女式格子衬衫, (39)Ivory, 165/88A(M)', '', '2019-07-25 00:39:04', '2019-07-25 00:39:04', null, 1),
			(1, 'SPYC949S1139095', '女式格子衬衫, (39)Ivory, 170/92A(L)', '', '2019-07-25 00:39:04', '2019-07-25 00:39:04', null, 1),
			(1, 'SPYC949S1159085', '女式格子衬衫, (59)Navy, 160/84A(S)', '', '2019-07-25 00:39:04', '2019-07-25 00:39:04', null, 1),
			(1, 'SPYC949S1159090', '女式格子衬衫, (59)Navy, 165/88A(M)', '', '2019-07-25 00:39:04', '2019-07-25 00:39:04', null, 1),
			(1, 'SPYC949S1159095', '女式格子衬衫, (59)Navy, 170/92A(L)', '', '2019-07-25 00:39:04', '2019-07-25 00:39:04', null, 1),
			(1, 'SPYC949S11NAONA', '女式格子衬衫, (NA)I.Color, 통합사이즈', '', '2019-07-25 00:39:04', '2019-07-25 00:39:04', null, 1),
			(2, 'SPWJ948S2255070', '女式牛仔裙, (55)Indigo, 160/66A(S)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(2, 'SPWJ948S2255075', '女式牛仔裙, (55)Indigo, 165/70A(M)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(2, 'SPWJ948S2255080', '女式牛仔裙, (55)Indigo, 170/74A(L)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(2, 'SPWJ948S2256070', '女式牛仔裙, (56)Light Indigo, 160/66A(S)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(2, 'SPWJ948S2256075', '女式牛仔裙, (56)Light Indigo, 165/70A(M)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(2, 'SPWJ948S2256080', '女式牛仔裙, (56)Light Indigo, 170/74A(L)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(2, 'SPWJ948S22NAONA', '女式牛仔裙, (NA)I.Color, 통합사이즈', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(3, 'SPWJ948S2355070', '女式牛仔裙, (55)Indigo, 160/66A(S)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(3, 'SPWJ948S2355075', '女式牛仔裙, (55)Indigo, 165/70A(M)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(3, 'SPWJ948S2355080', '女式牛仔裙, (55)Indigo, 170/74A(L)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(3, 'SPWJ948S2356070', '女式牛仔裙, (56)Light Indigo, 160/66A(S)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(3, 'SPWJ948S2356075', '女式牛仔裙, (56)Light Indigo, 165/70A(M)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(3, 'SPWJ948S2356080', '女式牛仔裙, (56)Light Indigo, 170/74A(L)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(3, 'SPWJ948S23NAONA', '女式牛仔裙, (NA)I.Color, 통합사이즈', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(4, 'SPYC949H2130095', '大格纹衬衫, (30)Yellow, 170/92A(M)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(4, 'SPYC949H2130100', '大格纹衬衫, (30)Yellow, 175/96A(L)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(4, 'SPYC949H2130105', '大格纹衬衫, (30)Yellow, 180/100A(XL)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(4, 'SPYC949H2185095', '大格纹衬衫, (85)Brown, 170/92A(M)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(4, 'SPYC949H2185100', '大格纹衬衫, (85)Brown, 175/96A(L)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(4, 'SPYC949H2185105', '大格纹衬衫, (85)Brown, 180/100A(XL)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(4, 'SPYC949H21NAONA', '大格纹衬衫, (NA)I.Color, 통합사이즈', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(4, 'SPYC949H2159095', '大格纹衬衫, (59)Navy, 170/92A(M)', '', '2019-07-25 00:39:18', '2019-07-25 00:39:18', null, 1),
			(4, 'SPYC949H2159100', '大格纹衬衫, (59)Navy, 175/96A(L)', '', '2019-07-25 00:39:18', '2019-07-25 00:39:18', null, 1),
			(4, 'SPYC949H2159105', '大格纹衬衫, (59)Navy, 180/100A(XL)', '', '2019-07-25 00:39:18', '2019-07-25 00:39:18', null, 1),
			(5, 'SPYS949H2230095', '条纹衬衫, (30)Yellow, 170/92A(M)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(5, 'SPYS949H2230100', '条纹衬衫, (30)Yellow, 175/96A(L)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(5, 'SPYS949H2230105', '条纹衬衫, (30)Yellow, 180/100A(XL)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(5, 'SPYS949H2250095', '条纹衬衫, (50)Blue, 170/92A(M)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(5, 'SPYS949H2250100', '条纹衬衫, (50)Blue, 175/96A(L)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(5, 'SPYS949H2250105', '条纹衬衫, (50)Blue, 180/100A(XL)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(5, 'SPYS949H2259095', '条纹衬衫, (59)Navy, 170/92A(M)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(5, 'SPYS949H2259100', '条纹衬衫, (59)Navy, 175/96A(L)', '', '2019-07-25 00:39:12', '2019-07-25 00:39:12', null, 1),
			(5, 'SPYS949H2259105', '条纹衬衫, (59)Navy, 180/100A(XL)', '', '2019-07-25 00:39:13', '2019-07-25 00:39:13', null, 1),
			(5, 'SPYS949H2285095', '条纹衬衫, (85)Brown, 170/92A(M)', '', '2019-07-25 00:39:13', '2019-07-25 00:39:13', null, 1),
			(5, 'SPYS949H2285100', '条纹衬衫, (85)Brown, 175/96A(L)', '', '2019-07-25 00:39:13', '2019-07-25 00:39:13', null, 1),
			(5, 'SPYS949H2285105', '条纹衬衫, (85)Brown, 180/100A(XL)', '', '2019-07-25 00:39:13', '2019-07-25 00:39:13', null, 1),
			(5, 'SPYS949H22NAONA', '条纹衬衫, (NA)I.Color, 통합사이즈', '', '2019-07-25 00:39:13', '2019-07-25 00:39:13', null, 1),
			(6, 'Q3AFAFDU6S2100230', '男女同款拖鞋, (00)生产代表颜色, 230MM', '', '2019-07-25 00:39:13', '2019-07-25 00:39:13', null, 1),
			(6, 'Q3AFAFDU6S2100240', '男女同款拖鞋, (00)生产代表颜色, 240MM', '', '2019-07-25 00:39:13', '2019-07-25 00:39:13', null, 1),
			(6, 'Q3AFAFDU6S2100250', '男女同款拖鞋, (00)生产代表颜色, 250MM', '', '2019-07-25 00:39:13', '2019-07-25 00:39:13', null, 1),
			(6, 'Q3AFAFDU6S2100260', '男女同款拖鞋, (00)生产代表颜色, 260MM', '', '2019-07-25 00:39:13', '2019-07-25 00:39:13', null, 1),
			(6, 'Q3AFAFDU6S2100270', '男女同款拖鞋, (00)生产代表颜色, 270MM', '', '2019-07-25 00:39:13', '2019-07-25 00:39:13', null, 1),
			(6, 'Q3AFAFDU6S2100280', '男女同款拖鞋, (00)生产代表颜色, 280MM', '', '2019-07-25 00:39:13', '2019-07-25 00:39:13', null, 1);
	`
	if _, err := session.Exec(sql); err != nil {
		log.Printf("setSkuData error: %v", err.Error())
		log.Println()
	}
}
