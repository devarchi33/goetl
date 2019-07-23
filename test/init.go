package test

import (
	"clearance-adapter/factory"
	"clearance-adapter/models"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	factory.Init()
	setUpCSLDB()
}

func setUpCSLDB() {
	truncateTables()
	initRecvSuppData()
}

func truncateTables() {
	sql := `
		TRUNCATE TABLE CSL.dbo.RecvSuppMst;
	`
	factory.GetCSLEngine().Exec(sql)
}

func initRecvSuppData() {
	masters := make([]models.RecvSuppMst, 0)
	masters = append(masters, models.RecvSuppMst{
		RecvSuppNo:         "CDVQ1907220001",
		BrandCode:          "EE",
		ShopCode:           "CDVQ",
		Dates:              "20190722",
		RecvSuppType:       "R",
		ShippingTypeCode:   "01",
		WayBillNo:          "SR20190722001",
		RecvSuppStatusCode: "R",
	})
	masters = append(masters, models.RecvSuppMst{
		RecvSuppNo:         "CDVQ1907220002",
		BrandCode:          "EE",
		ShopCode:           "CDVQ",
		Dates:              "20190722",
		RecvSuppType:       "R",
		ShippingTypeCode:   "01",
		WayBillNo:          "SR20190722002",
		RecvSuppStatusCode: "R",
	})
	factory.GetCSLEngine().Insert(&masters)
}
