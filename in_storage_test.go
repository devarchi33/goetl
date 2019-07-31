package main

import (
	"clearance-adapter/domain/entities"
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
				}, map[string]string{
					"brand_code": "SA",
					"shop_code":  "CEGP",
					"waybill_no": "1010590009009",
					"box_no":     "1010590009009",
					"sku_code":   "SPYC949H2130095",
					"qty":        "4",
					"user_id":    "shi.yanxun",
				}, map[string]string{
					"brand_code": "SA",
					"shop_code":  "CEGP",
					"waybill_no": "1010590009009",
					"box_no":     "1010590009009",
					"sku_code":   "SPYC949H2130100",
					"qty":        "5",
					"user_id":    "shi.yanxun",
				},
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
			So(transactions[1].BrandCode, ShouldEqual, "SA")
			So(transactions[1].ShopCode, ShouldEqual, "CEGP")
			So(transactions[1].WaybillNo, ShouldEqual, "1010590009009")
			So(transactions[1].BoxNo, ShouldEqual, "1010590009009")
			So(transactions[1].UserID, ShouldEqual, "shi.yanxun")
			So(transactions[1].Items, ShouldNotBeNil)
			So(len(transactions[1].Items), ShouldEqual, 2)
			So(transactions[1].Items[0].SkuCode, ShouldEqual, "SPYC949H2130095")
			So(transactions[1].Items[0].Qty, ShouldEqual, 4)
			So(transactions[1].Items[1].SkuCode, ShouldEqual, "SPYC949H2130100")
			So(transactions[1].Items[1].Qty, ShouldEqual, 5)
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

// func TestInStorageETLTransform(t *testing.T) {
// 	Convey("测试InStorageETL的Transform方法", t, func() {
// 		Convey("将Transctions匹配成RecvSuppMst和RecvSuppDtl", func() {
// 			transactions := []models.Transaction{
// 				models.Transaction{
// 					ID:            1,
// 					TransactionID: "CFW51907236000",
// 					WaybillNo:     "1010590009017",
// 					BoxNo:         "1010590009017",
// 					SkuCode:       "SPWH937G8951075",
// 					Qty:           2,
// 				},
// 				models.Transaction{
// 					ID:            1,
// 					TransactionID: "CFW51907236000",
// 					WaybillNo:     "1010590009017",
// 					BoxNo:         "1010590009017",
// 					SkuCode:       "SPAC937D0219999",
// 					Qty:           2,
// 				},
// 			}
// 			mstDtlMap, err := InStorageETL{}.Transform(context.Background(), transactions)
// 			So(err, ShouldBeNil)
// 			masters := mstDtlMap.(map[string]interface{})["RecvSuppMst"].([]models.RecvSuppMst)
// 			So(masters, ShouldNotBeNil)
// 			So(masters[0].RecvSuppNo, ShouldEqual, "CFW51907236000")
// 			So(masters[0].WayBillNo, ShouldEqual, "1010590009017")
// 			So(masters[0].BoxNo, ShouldEqual, "1010590009017")
// 			details := mstDtlMap.(map[string]interface{})["RecvSuppDtl"].([]models.RecvSuppDtl)
// 			So(details, ShouldNotBeNil)
// 			So(details[0].ProdCode, ShouldEqual, "SPWH937G8951075")
// 			So(details[0].RecvSuppFixedQty, ShouldEqual, 2)
// 			So(details[1].ProdCode, ShouldEqual, "SPAC937D0219999")
// 			So(details[1].RecvSuppFixedQty, ShouldEqual, 2)
// 		})
// 	})
// }
