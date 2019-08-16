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

func TestDistributionETLBuildDistributions(t *testing.T) {
	Convey("测试buildDistributions方法", t, func() {
		Convey("源数据中包含多个运单号的数据，应该能够根据运单号生成DistributionOrder", func() {
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

			result, err := DistributionETL{}.buildDistributionOrders(context.Background(), data)
			if err != nil {
				log.Printf(err.Error())
			}
			So(err, ShouldBeNil)
			orders, ok := result.([]entities.DistributionOrder)
			if !ok {
				log.Printf("Convert Failed")
			}
			So(len(orders), ShouldEqual, 2)
			for _, order := range orders {
				if order.BrandCode == "SA" && order.ShopCode == "CEGP" && order.WaybillNo == "1010590009008" {
					So(order.BoxNo, ShouldEqual, "1010590009008")
					So(order.EmpID, ShouldEqual, "7000028260")
					So(order.Items, ShouldNotBeNil)
					So(len(order.Items), ShouldEqual, 2)
					So(order.Items[0].SkuCode, ShouldEqual, "SPWJ948S2255070")
					So(order.Items[0].Qty, ShouldEqual, 2)
					So(order.Items[1].SkuCode, ShouldEqual, "SPWJ948S2256070")
					So(order.Items[1].Qty, ShouldEqual, 3)
				} else if order.BrandCode == "SA" && order.ShopCode == "CEGP" && order.WaybillNo == "1010590009009" {
					So(order.BoxNo, ShouldEqual, "1010590009009")
					So(order.EmpID, ShouldEqual, "7000028260")
					So(order.Items, ShouldNotBeNil)
					So(len(order.Items), ShouldEqual, 2)
					So(order.Items[0].SkuCode, ShouldEqual, "SPYC949H2130095")
					So(order.Items[0].Qty, ShouldEqual, 4)
					So(order.Items[1].SkuCode, ShouldEqual, "SPYC949H2130100")
					So(order.Items[1].Qty, ShouldEqual, 5)
				}
			}
		})
	})
}

func syncedShouldBeTrue(receiptLocationCode, waybillNo string) {
	sql := `
		SELECT
			sd.synced
		FROM pangpang_brand_sku_location.stock_distribute AS sd
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = sd.receipt_location_id
		WHERE sd.tenant_code = 'pangpang'
			AND store.code = ?
			AND sd.waybill_no = ?
	`
	result, err := factory.GetP2BrandEngine().Query(sql, receiptLocationCode, waybillNo)
	So(err, ShouldBeNil)
	So(len(result), ShouldEqual, 1)
	distList := infra.ConvertByteResult(result)
	So(distList[0]["synced"], ShouldEqual, "1")
}

func TestDistributionETL(t *testing.T) {
	Convey("测试DistributionETL的Run方法", t, func() {
		Convey("运单号为1010590009007的出库单synced=true，所以应该不会同步到CSL入库", func() {
			etl := DistributionETL{}.New()
			err := etl.Run(context.Background())
			So(err, ShouldBeNil)
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", "CEGP", "1010590009007")
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 2)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, false)
			}
		})
		Convey("运单号为1010590009008的运单应该在CSL入库", func() {
			receiptLocationCode := "CEGP"
			waybillNo := "1010590009008"
			etl := DistributionETL{}.New()
			err := etl.Run(context.Background())
			So(err, ShouldBeNil)
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", receiptLocationCode, waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 14)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, true)
				So(recvSupp.RecvEmpID, ShouldEqual, "7000028260")
				So(recvSupp.RecvSuppMst.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppDtl.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvEmpName, ShouldEqual, "史妍珣")
			}
			syncedShouldBeTrue(receiptLocationCode, waybillNo)
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
			receiptLocationCode := "CFGY"
			waybillNo := "1010590009014"
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", receiptLocationCode, waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 11)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, true)
				So(recvSupp.RecvEmpID, ShouldEqual, "7000028260")
				So(recvSupp.RecvSuppMst.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppDtl.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvEmpName, ShouldEqual, "史妍珣")
			}
			syncedShouldBeTrue(receiptLocationCode, waybillNo)
		})

		Convey("CEGP的子卖场CJC1运单号为1010590009009的运单应该在CSL入库", func() {
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("Q3", "CJC1", "1010590009009")
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 5)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, true)
				So(recvSupp.RecvEmpID, ShouldEqual, "7000028260")
				So(recvSupp.RecvSuppMst.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppDtl.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvEmpName, ShouldEqual, "史妍珣")
			}
			receiptLocationCode := "CEGP"
			waybillNo := "1010590009009"
			syncedShouldBeTrue(receiptLocationCode, waybillNo)
		})
	})
}
