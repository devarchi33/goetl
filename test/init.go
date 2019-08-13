package test

import (
	"clearance-adapter/factory"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	factory.Init()
	initCSL()
	initP2BrandDB()
}

func initP2BrandDB() {
	initStore()
	initProduct()
	initLocation()
}

// InitTestData 初始化测试数据
func InitTestData() {
	initCSL()
	initP2BrandDB()
}
