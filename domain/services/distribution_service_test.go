package services

import (
	clrConst "clearance-adapter/domain/clr-constants"
	"clearance-adapter/domain/entities"
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"clearance-adapter/repositories"
	_ "clearance-adapter/test"
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDistributionETLBuildDistributions(t *testing.T) {
	Convey("测试buildDistributions方法", t, func() {
		Convey("源数据中包含多个运单号的数据，应该能够根据运单号生成DistributionOrder", func() {
			data := []map[string]string{
				map[string]string{
					"brand_code":            "SA",
					"receipt_location_code": "CEGP",
					"waybill_no":            "1010590009008",
					"box_no":                "1010590009008",
					"sku_code":              "SPWJ948S2255070",
					"qty":                   "2",
					"in_date":               "2019-08-23T13:29:05Z",
					"in_emp_id":             "7000028260",
				}, map[string]string{
					"brand_code":            "SA",
					"receipt_location_code": "CEGP",
					"waybill_no":            "1010590009008",
					"box_no":                "1010590009008",
					"sku_code":              "SPWJ948S2256070",
					"qty":                   "3",
					"in_date":               "2019-08-23T13:29:05Z",
					"in_emp_id":             "7000028260",
				},
				map[string]string{
					"brand_code":            "SA",
					"receipt_location_code": "CEGP",
					"waybill_no":            "1010590009009",
					"box_no":                "1010590009009",
					"sku_code":              "SPYC949H2130095",
					"qty":                   "4",
					"in_date":               "2019-08-23T13:29:05Z",
					"in_emp_id":             "7000028260",
				}, map[string]string{
					"brand_code":            "SA",
					"receipt_location_code": "CEGP",
					"waybill_no":            "1010590009009",
					"box_no":                "1010590009009",
					"sku_code":              "SPYC949H2130100",
					"qty":                   "5",
					"in_date":               "2019-08-23T13:29:05Z",
					"in_emp_id":             "7000028260",
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
				if order.BrandCode == "SA" && order.ReceiptLocationCode == "CEGP" && order.WaybillNo == "1010590009008" {
					So(order.BoxNo, ShouldEqual, "1010590009008")
					So(order.InEmpID, ShouldEqual, "7000028260")
					So(order.Items, ShouldNotBeNil)
					So(len(order.Items), ShouldEqual, 2)
					So(order.Items[0].SkuCode, ShouldEqual, "SPWJ948S2255070")
					So(order.Items[0].Qty, ShouldEqual, 2)
					So(order.Items[1].SkuCode, ShouldEqual, "SPWJ948S2256070")
					So(order.Items[1].Qty, ShouldEqual, 3)
				} else if order.BrandCode == "SA" && order.ReceiptLocationCode == "CEGP" && order.WaybillNo == "1010590009009" {
					So(order.BoxNo, ShouldEqual, "1010590009009")
					So(order.InEmpID, ShouldEqual, "7000028260")
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
	utc, _ := time.LoadLocation("")
	startDate := time.Now().Add(-1000).In(utc).Format("2006-01-02T15:04:05Z")
	sql := `
		SELECT
			sd.synced,
			sd.last_updated_at
		FROM pangpang_brand_sku_location.stock_distribute AS sd
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = sd.receipt_location_id
		WHERE sd.tenant_code = 'pangpang'
			AND store.code = ?
			AND sd.waybill_no = ?
	`
	result, err := factory.GetP2BrandEngine().Query(sql, receiptLocationCode, waybillNo)
	endDate := time.Now().In(utc).Format("2006-01-02T15:04:05Z")
	So(err, ShouldBeNil)
	So(len(result), ShouldEqual, 1)
	distList := infra.ConvertByteResult(result)
	So(distList[0]["synced"], ShouldEqual, "1")
	So(distList[0]["last_updated_at"], ShouldBeGreaterThanOrEqualTo, startDate)
	So(distList[0]["last_updated_at"], ShouldBeLessThanOrEqualTo, endDate)
}

func directSyncedShouldBeTrue(receiptLocationCode, waybillNo string) {
	utc, _ := time.LoadLocation("")
	startDate := time.Now().Add(-2000).In(utc).Format("2006-01-02T15:04:05Z")
	sql := `
		SELECT
			dd.synced,
			dd.last_updated_at
		FROM pangpang_brand_sku_location.direct_distribution AS dd
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = dd.receipt_location_id
		WHERE dd.tenant_code = 'pangpang'
			AND store.code = ?
			AND dd.waybill_no = ?
	`
	result, err := factory.GetP2BrandEngine().Query(sql, receiptLocationCode, waybillNo)
	endDate := time.Now().In(utc).Format("2006-01-02T15:04:05Z")
	So(err, ShouldBeNil)
	So(len(result), ShouldEqual, 1)
	distList := infra.ConvertByteResult(result)
	So(distList[0]["synced"], ShouldEqual, "1")
	So(distList[0]["last_updated_at"], ShouldBeGreaterThanOrEqualTo, startDate)
	So(distList[0]["last_updated_at"], ShouldBeLessThanOrEqualTo, endDate)
}

func TestDistributionETL(t *testing.T) {
	Convey("测试DistributionETL的Run方法", t, func() {
		etl := DistributionETL{}.New()
		err := etl.Run(context.Background())
		So(err, ShouldBeNil)
		Convey("运单号为1010590009007的出库单synced=true，所以应该不会同步到CSL入库", func() {
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

		Convey("CFGY卖场，运单号为1010590009015的运单应该在CSL入库", func() {
			receiptLocationCode := "CFGY"
			waybillNo := "1010590009015"
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", receiptLocationCode, waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 11)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, true)
				So(recvSupp.RecvEmpID, ShouldEqual, "7000028260")
				So(recvSupp.RecvSuppMst.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppDtl.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvEmpName, ShouldEqual, "史妍珣")
				So(recvSupp.ShopSuppRecvDate, ShouldEqual, "20190820")
				So(recvSupp.InvtBaseDate, ShouldEqual, "20190820")
			}
			syncedShouldBeTrue(receiptLocationCode, waybillNo)
		})

		Convey("CFGY卖场，运单号为1010590009016的运单存在两个master，并且master中的商品是相同的，如果该商品存在误差应该合并数量后再登记", func() {
			/*
				sku					出			入			误差
				SPWJ948S2255075		2（1+1）	6（3+3）	4
			*/
			sql := `
				SELECT * FROM CSL.dbo.StockMisDtl
				WHERE BrandCode = ?
				AND ShopCode = ?
				AND WayBillNo01 = ?
			`
			result, _ := factory.GetCSLEngine().Query(sql, "SA", "CFGY", "1010590009016")
			So(len(result), ShouldBeGreaterThanOrEqualTo, 1)
			stockMissList := infra.ConvertByteResult(result)
			outQty := 0
			missQty := 0
			for _, stockMiss := range stockMissList {
				skuCode := stockMiss["ProdCode"]
				if skuCode == "SPWJ948S2255075" {
					recvSuppQty, _ := strconv.Atoi(stockMiss["RecvSuppQty"])
					outQty += recvSuppQty
					stockMisQty, _ := strconv.Atoi(stockMiss["StockMisQty"])
					missQty += stockMisQty
				}
			}
			So(outQty, ShouldEqual, 1)
			So(missQty, ShouldEqual, 4)
		})

		brandCode := "SA"
		recptLocCode := "CFGY"
		waybillNo := "20190906001"
		title := fmt.Sprintf("%v品牌，%v卖场，运单号为：%v的运单应该同步失败并且在error表中有记录", brandCode, recptLocCode, waybillNo)
		Convey(title, func() {
			has, distError, err := repositories.StockDistributionErrorRepository{}.GetByWaybillNo(brandCode, recptLocCode, waybillNo)
			So(has, ShouldEqual, true)
			So(err, ShouldBeNil)
			So(distError.Type, ShouldEqual, clrConst.TypStockDistributionError)
			So(len(distError.ErrorMessage), ShouldBeGreaterThan, 0)
		})

		brandCode = "SA"
		recptLocCode = "CEGP"
		waybillNo = "20190909001"
		Convey(fmt.Sprintf("【工厂直送】运单号为%v的运单应该在CSL入库", waybillNo), func() {
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", recptLocCode, waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 2)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, true)
				So(recvSupp.RecvEmpID, ShouldEqual, "7000028260")
				So(recvSupp.RecvSuppMst.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppDtl.ModiUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvEmpName, ShouldEqual, "史妍珣")
			}
			directSyncedShouldBeTrue(recptLocCode, waybillNo)
		})
	})
}
