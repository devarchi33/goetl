package services

import (
	"clearance-adapter/domain/entities"
	"clearance-adapter/factory"
	"clearance-adapter/infra"
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
					"emp_id":     "7000028260",
				}, map[string]string{
					"brand_code": "SA",
					"shop_code":  "CEGP",
					"waybill_no": "1010590009008",
					"box_no":     "1010590009008",
					"sku_code":   "SPWJ948S2256070",
					"qty":        "3",
					"emp_id":     "7000028260",
				},
				map[string]string{
					"brand_code": "SA",
					"shop_code":  "CEGP",
					"waybill_no": "1010590009009",
					"box_no":     "1010590009009",
					"sku_code":   "SPYC949H2130095",
					"qty":        "4",
					"emp_id":     "7000028260",
				}, map[string]string{
					"brand_code": "SA",
					"shop_code":  "CEGP",
					"waybill_no": "1010590009009",
					"box_no":     "1010590009009",
					"sku_code":   "SPYC949H2130100",
					"qty":        "5",
					"emp_id":     "7000028260",
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
			So(len(transactions), ShouldEqual, 2)
			for _, txn := range transactions {
				if txn.BrandCode == "SA" && txn.ShopCode == "CEGP" && txn.WaybillNo == "1010590009008" {
					So(txn.BoxNo, ShouldEqual, "1010590009008")
					So(txn.EmpID, ShouldEqual, "7000028260")
					So(txn.Items, ShouldNotBeNil)
					So(len(txn.Items), ShouldEqual, 2)
					So(txn.Items[0].SkuCode, ShouldEqual, "SPWJ948S2255070")
					So(txn.Items[0].Qty, ShouldEqual, 2)
					So(txn.Items[1].SkuCode, ShouldEqual, "SPWJ948S2256070")
					So(txn.Items[1].Qty, ShouldEqual, 3)
				} else if txn.BrandCode == "SA" && txn.ShopCode == "CEGP" && txn.WaybillNo == "1010590009009" {
					So(txn.BoxNo, ShouldEqual, "1010590009009")
					So(txn.EmpID, ShouldEqual, "7000028260")
					So(txn.Items, ShouldNotBeNil)
					So(len(txn.Items), ShouldEqual, 2)
					So(txn.Items[0].SkuCode, ShouldEqual, "SPYC949H2130095")
					So(txn.Items[0].Qty, ShouldEqual, 4)
					So(txn.Items[1].SkuCode, ShouldEqual, "SPYC949H2130100")
					So(txn.Items[1].Qty, ShouldEqual, 5)
				}
			}
		})
	})
}

func TestInStorageETL(t *testing.T) {
	// SELECT * FROM tansactions
	Convey("测试InStorageETL的Run方法", t, func() {
		Convey("某个时间段没有入库运单的话，应该没有数据在CSL入库", func() {
			etl := InStorageETL{}.New("2019-07-01 00:00:00", "2019-07-01 00:01:00")
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
				So(recvSupp.RecvSuppMst.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppDtl.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvEmpName, ShouldEqual, "史妍珣")
			}
		})
		Convey("运单号为1010590009008的运单应该存在有误差的商品", func() {
			// SPYC949H2159095 4, 5 +1
			// SPYC949H2159100 3, 2 -1
			sql := `
				SELECT * FROM CSL.dbo.StockMisDtl
				WHERE BrandCode = ?
				AND ShopCode = ?
				AND WayBillNo01 = ?
			`
			result, _ := factory.GetCSLEngine().Query(sql, "SA", "CEGP", "1010590009008")
			So(len(result), ShouldEqual, 3)
			stockMissList := infra.ConvertByteResult(result)
			for _, stockMiss := range stockMissList {
				skuCode := stockMiss["ProdCode"]
				if skuCode == "SPYC949H2159095" {
					So(stockMiss["RecvSuppQty"], ShouldEqual, "4")
					So(stockMiss["StockMisQty"], ShouldEqual, "1")
					So(stockMiss["RecvSuppType"], ShouldEqual, "R")
				} else if skuCode == "SPYC949H2159100" {
					So(stockMiss["RecvSuppQty"], ShouldEqual, "3")
					So(stockMiss["StockMisQty"], ShouldEqual, "1")
					So(stockMiss["RecvSuppType"], ShouldEqual, "S")
				} else if skuCode == "SPWJ948S2256070" {
					So(stockMiss["RecvSuppQty"], ShouldEqual, "4")
					So(stockMiss["StockMisQty"], ShouldEqual, "4")
					So(stockMiss["RecvSuppType"], ShouldEqual, "S")
				}
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
				So(recvSupp.RecvSuppMst.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppDtl.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvEmpName, ShouldEqual, "史妍珣")
			}
		})
	})
}
