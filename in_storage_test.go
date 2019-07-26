package main

import (
	"clearance-adapter/models"
	_ "clearance-adapter/test"
	"context"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInStorageETLTransform(t *testing.T) {
	Convey("测试InStorageETL的Transform方法", t, func() {
		Convey("将Transctions匹配成RecvSuppMst和RecvSuppDtl", func() {
			transactions := []models.Transaction{
				models.Transaction{
					ID:            1,
					TransactionID: "CFW51907236000",
					WaybillNo:     "1010590009017",
					BoxNo:         "1010590009017",
					SkuCode:       "SPWH937G8951075",
					Qty:           2,
				},
				models.Transaction{
					ID:            1,
					TransactionID: "CFW51907236000",
					WaybillNo:     "1010590009017",
					BoxNo:         "1010590009017",
					SkuCode:       "SPAC937D0219999",
					Qty:           2,
				},
			}
			mstDtlMap, err := InStorageETL{}.Transform(context.Background(), transactions)
			So(err, ShouldBeNil)
			masters := mstDtlMap.(map[string]interface{})["RecvSuppMst"].([]models.RecvSuppMst)
			So(masters, ShouldNotBeNil)
			So(masters[0].RecvSuppNo, ShouldEqual, "CFW51907236000")
			So(masters[0].WayBillNo, ShouldEqual, "1010590009017")
			So(masters[0].BoxNo, ShouldEqual, "1010590009017")
			details := mstDtlMap.(map[string]interface{})["RecvSuppDtl"].([]models.RecvSuppDtl)
			So(details, ShouldNotBeNil)
			So(details[0].ProdCode, ShouldEqual, "SPWH937G8951075")
			So(details[0].RecvSuppFixedQty, ShouldEqual, 2)
			So(details[1].ProdCode, ShouldEqual, "SPAC937D0219999")
			So(details[1].RecvSuppFixedQty, ShouldEqual, 2)
		})
	})
}

func TestInStorageETL(t *testing.T) {
	// SELECT * FROM tansactions
	Convey("测试InStorageETL的Run方法", t, func() {
		Convey("可以把入库数据从Clearance导入到CSL", func() {
			etl := InStorageETL{}.New()
			err := etl.Run(context.Background())
			So(err, ShouldBeNil)
		})
	})
}
