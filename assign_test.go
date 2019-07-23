package goetl

import (
	"clearance-adapter/models"
	_ "clearance-adapter/test"
	"context"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTransform(t *testing.T) {
	Convey("测试AssignETL的Transform方法", t, func() {
		Convey("如果map中的Key都是正确的，应该得到正确的TransactionMaster数组", func() {
			source := []map[string]string{
				map[string]string{
					"RecvSuppNo":         "EECDVQ201907220001",
					"BrandCode":          "EE",
					"ShopCode":           "CDVQ",
					"Dates":              "20190722",
					"ShippingTypeCode":   "01",
					"WayBillNo":          "SR201907220001",
					"RecvSuppStatusCode": "R",
				},
			}
			masters, err := AssignETL{}.Transform(context.Background(), source)
			So(err, ShouldBeNil)
			master := masters.([]models.TransactionMaster)[0]
			So(master.ID, ShouldEqual, 0)
			So(master.Date, ShouldEqual, "20190722")
			So(master.PlantCode, ShouldEqual, "EE-CDVQ")
			So(master.WaybillNo, ShouldEqual, "SR201907220001")
			So(master.OrderNo, ShouldEqual, "EECDVQ201907220001")
			So(master.TransactionCode, ShouldEqual, "OS100")
			So(master.Channel, ShouldEqual, "CLEARANCE")
		})
	})
}

func TestAssignETL(t *testing.T) {
	Convey("测试AssignETL的Run方法", t, func() {
		Convey("可以把入库预约从CSL导入到MSL", func() {
			etl := New(AssignETL{})
			etl.AfterTransform(AssignETL{}.ReadyToLoad)
			err := etl.Run(context.Background())
			So(err, ShouldBeNil)
		})
	})
}
