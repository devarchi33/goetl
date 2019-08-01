package services

import (
	"clearance-adapter/domain/entities"
	"clearance-adapter/repositories"
	_ "clearance-adapter/test"
	"context"
	"log"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInStorageETLBuildTransactions(t *testing.T) {
	Convey("测试buildTransactions方法", t, func() {
		Convey("源数据中包含多个运单号的数据，应该能够根据运单号生成Transaction", func() {
			data := []map[string]string{
				map[string]string{
					"brand_code": "SA",
					"shop_code":  "CEGP",
					"waybill_no": "1010590009008",
					"box_no":     "1010590009008",
					"sku_code":   "SPWJ948S2255070",
					"qty":        "2",
					"user_id":    "shi.yanxun",
				}, map[string]string{
					"brand_code": "SA",
					"shop_code":  "CEGP",
					"waybill_no": "1010590009008",
					"box_no":     "1010590009008",
					"sku_code":   "SPWJ948S2256070",
					"qty":        "3",
					"user_id":    "shi.yanxun",
				},
				// map[string]string{
				// 	"brand_code": "SA",
				// 	"shop_code":  "CEGP",
				// 	"waybill_no": "1010590009009",
				// 	"box_no":     "1010590009009",
				// 	"sku_code":   "SPYC949H2130095",
				// 	"qty":        "4",
				// 	"user_id":    "shi.yanxun",
				// }, map[string]string{
				// 	"brand_code": "SA",
				// 	"shop_code":  "CEGP",
				// 	"waybill_no": "1010590009009",
				// 	"box_no":     "1010590009009",
				// 	"sku_code":   "SPYC949H2130100",
				// 	"qty":        "5",
				// 	"user_id":    "shi.yanxun",
				// },
			}

			result, err := InStorageETL{}.buildTransactions(context.Background(), data)
			if err != nil {
				log.Printf(err.Error())
			}
			So(err, ShouldBeNil)
			transactions, ok := result.([]entities.Transaction)
			if !ok {
				log.Printf("Convert Failed")
			}
			// 第一条
			So(transactions[0].BrandCode, ShouldEqual, "SA")
			So(transactions[0].ShopCode, ShouldEqual, "CEGP")
			So(transactions[0].WaybillNo, ShouldEqual, "1010590009008")
			So(transactions[0].BoxNo, ShouldEqual, "1010590009008")
			So(transactions[0].UserID, ShouldEqual, "shi.yanxun")
			So(transactions[0].Items, ShouldNotBeNil)
			So(len(transactions[0].Items), ShouldEqual, 2)
			So(transactions[0].Items[0].SkuCode, ShouldEqual, "SPWJ948S2255070")
			So(transactions[0].Items[0].Qty, ShouldEqual, 2)
			So(transactions[0].Items[1].SkuCode, ShouldEqual, "SPWJ948S2256070")
			So(transactions[0].Items[1].Qty, ShouldEqual, 3)
			// 第二条
			// So(transactions[1].BrandCode, ShouldEqual, "SA")
			// So(transactions[1].ShopCode, ShouldEqual, "CEGP")
			// So(transactions[1].WaybillNo, ShouldEqual, "1010590009009")
			// So(transactions[1].BoxNo, ShouldEqual, "1010590009009")
			// So(transactions[1].UserID, ShouldEqual, "shi.yanxun")
			// So(transactions[1].Items, ShouldNotBeNil)
			// So(len(transactions[1].Items), ShouldEqual, 2)
			// So(transactions[1].Items[0].SkuCode, ShouldEqual, "SPYC949H2130095")
			// So(transactions[1].Items[0].Qty, ShouldEqual, 4)
			// So(transactions[1].Items[1].SkuCode, ShouldEqual, "SPYC949H2130100")
			// So(transactions[1].Items[1].Qty, ShouldEqual, 5)
		})
	})
}

func TestInStorageETL(t *testing.T) {
	// SELECT * FROM tansactions
	Convey("测试InStorageETL的Run方法", t, func() {
		Convey("某个时间段没有入库运单的话，应该没有数据在CSL入库", func() {
			etl := InStorageETL{}.New("2019-07-01 00:00:00", "2019-07-01 00:00:00")
			err := etl.Run(context.Background())
			So(err, ShouldBeNil)
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", "CEGP", "1010590009008")
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 14)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, false)
			}
		})
		Convey("运单号为1010590009008的运单应该在CSL入库", func() {
			etl := InStorageETL{}.New("2019-07-01 00:00:00", "2019-07-31 23:59:59")
			err := etl.Run(context.Background())
			So(err, ShouldBeNil)
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", "CEGP", "1010590009008")
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 14)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, true)
				So(recvSupp.RecvEmpID, ShouldEqual, "7000028260")
				So(recvSupp.RecvEmpName, ShouldEqual, "史妍珣")
			}
		})
		Convey("运单号为1010590009014的运单应该在CSL入库", func() {
			etl := InStorageETL{}.New("2019-07-01 00:00:00", "2019-07-31 23:59:59")
			err := etl.Run(context.Background())
			So(err, ShouldBeNil)
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", "CFGY", "1010590009014")
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 11)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, true)
				So(recvSupp.RecvEmpID, ShouldEqual, "7000028260")
				So(recvSupp.RecvEmpName, ShouldEqual, "史妍珣")
			}
		})
	})
}
