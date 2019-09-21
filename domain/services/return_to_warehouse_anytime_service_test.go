package services

import (
	clrConst "clearance-adapter/domain/clr-constants"
	cslConst "clearance-adapter/domain/csl-constants"
	"clearance-adapter/domain/entities"
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"clearance-adapter/repositories"
	_ "clearance-adapter/test"
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	p2bConst "clearance-adapter/domain/p2brand-constants"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTransform(t *testing.T) {
	Convey("测试Transform方法", t, func() {
		Convey("源数据中包含多个运单号的数据，应该能够根据运单号生成ReturnToWarehouseOrder", func() {
			data := []map[string]string{
				map[string]string{
					"brand_code":             "SA",
					"shipment_location_code": "CEGP",
					"waybill_no":             "1010590009008",
					"status_code":            p2bConst.StsSentOut,
					"sku_code":               "SPWJ948S2255070",
					"qty":                    "2",
					"emp_id":                 "7000028260",
					"out_date":               "2019-08-22T19:01:01Z",
				}, map[string]string{
					"brand_code":             "SA",
					"shipment_location_code": "CEGP",
					"waybill_no":             "1010590009008",
					"status_code":            p2bConst.StsSentOut,
					"sku_code":               "SPWJ948S2256070",
					"qty":                    "3",
					"emp_id":                 "7000028260",
					"out_date":               "2019-08-22T19:01:01Z",
				},
				map[string]string{
					"brand_code":             "SA",
					"shipment_location_code": "CEGP",
					"waybill_no":             "1010590009009",
					"status_code":            p2bConst.StsSentOut,
					"sku_code":               "SPYC949H2130095",
					"qty":                    "4",
					"emp_id":                 "7000028260",
					"out_date":               "2019-08-22T19:01:01Z",
				}, map[string]string{
					"brand_code":             "SA",
					"shipment_location_code": "CEGP",
					"waybill_no":             "1010590009009",
					"status_code":            p2bConst.StsSentOut,
					"sku_code":               "SPYC949H2130100",
					"qty":                    "5",
					"emp_id":                 "7000028260",
					"out_date":               "2019-08-22T19:01:01Z",
				},
			}

			result, err := ReturnToWarehouseAnytimeETL{}.Transform(context.Background(), data)
			if err != nil {
				log.Printf(err.Error())
			}
			So(err, ShouldBeNil)
			orders, ok := result.([]entities.ReturnToWarehouseOrder)
			if !ok {
				log.Printf("Convert Failed")
			}
			So(len(orders), ShouldEqual, 2)
			for _, order := range orders {
				if order.BrandCode == "SA" && order.ShipmentLocationCode == "CEGP" && order.WaybillNo == "1010590009008" {
					So(order.StatusCode, ShouldEqual, p2bConst.StsSentOut)
					So(order.EmpID, ShouldEqual, "7000028260")
					So(order.OutDate, ShouldEqual, "20190823")
					So(order.Items, ShouldNotBeNil)
					So(len(order.Items), ShouldEqual, 2)
					So(order.Items[0].SkuCode, ShouldEqual, "SPWJ948S2255070")
					So(order.Items[0].Qty, ShouldEqual, 2)
					So(order.Items[1].SkuCode, ShouldEqual, "SPWJ948S2256070")
					So(order.Items[1].Qty, ShouldEqual, 3)
				} else if order.BrandCode == "SA" && order.ShipmentLocationCode == "CEGP" && order.WaybillNo == "1010590009009" {
					So(order.StatusCode, ShouldEqual, p2bConst.StsSentOut)
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

func rtwAnytimeSyncedShouldBeTrueOrFalse(shipmentLocationCode, waybillNo string, shouldBeTrue bool) {
	utc, _ := time.LoadLocation("")
	startDate := time.Now().In(utc).Format("2006-01-02T15:04:05Z")
	sql := `
		SELECT
			rtw.synced,
			rtw.last_updated_at
		FROM pangpang_brand_sku_location.return_to_warehouse_anytime AS rtw
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = rtw.shipment_location_id
		WHERE rtw.tenant_code = 'pangpang'
			AND store.code = ?
			AND rtw.waybill_no = ?
	`
	result, err := factory.GetP2BrandEngine().Query(sql, shipmentLocationCode, waybillNo)
	endDate := time.Now().In(utc).Format("2006-01-02T15:04:05Z")
	So(err, ShouldBeNil)
	So(len(result), ShouldEqual, 1)
	distList := infra.ConvertByteResult(result)
	if shouldBeTrue {
		So(distList[0]["synced"], ShouldEqual, "1")
	} else {
		So(distList[0]["synced"], ShouldEqual, "0")
	}
	So(distList[0]["last_updated_at"], ShouldBeGreaterThanOrEqualTo, startDate)
	So(distList[0]["last_updated_at"], ShouldBeLessThanOrEqualTo, endDate)
}

func TestReturnToWarehouseAnytimeETL(t *testing.T) {
	Convey("测试ReturnToWarehouseAnytimeETL的Run方法", t, func() {
		etl := ReturnToWarehouseAnytimeETL{}.New()
		err := etl.Run(context.Background())
		So(err, ShouldBeNil)
		waybillNo := "A20190819001"
		Convey("运单号为"+waybillNo+"的出库单synced=true，所以应该不会同步到CSL入库", func() {
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", "CEGP", waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 0)
		})
		waybillNo = "A20190813001"
		Convey("运单号为"+waybillNo+"的退仓出库单应该在CSL有记录", func() {
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", "CEGP", waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 2)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, false)
				So(recvSupp.ShippingTypeCode, ShouldEqual, cslConst.TypRTWAnytime)
				So(recvSupp.RecvSuppType, ShouldEqual, cslConst.TypSentOut)
				So(recvSupp.RecvSuppStatusCode, ShouldEqual, cslConst.StsSentOut)
				So(recvSupp.SuppEmpID, ShouldEqual, "7000028260")
				So(recvSupp.SuppEmpName, ShouldEqual, "史妍珣")
				So(recvSupp.RecvSuppMst.InUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppDtl.InUserID, ShouldEqual, "shi.yanxun")
				if recvSupp.ProdCode == "SPWJ948S2255070" {
					So(recvSupp.RecvSuppQty, ShouldEqual, 4)
				}
				if recvSupp.ProdCode == "SPYC949S1139095" {
					So(recvSupp.RecvSuppQty, ShouldEqual, 1)
				}
			}
			rtwAnytimeSyncedShouldBeTrueOrFalse("CEGP", waybillNo, true)
		})

		waybillNo = "A20190814001"
		Convey("CEGP的子卖场CJC1运单号为"+waybillNo+"的退仓出库单应该在CSL有记录", func() {
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("Q3", "CJC1", waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 2)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, false)
				So(recvSupp.ShippingTypeCode, ShouldEqual, cslConst.TypRTWAnytime)
				So(recvSupp.RecvSuppType, ShouldEqual, cslConst.TypSentOut)
				So(recvSupp.RecvSuppStatusCode, ShouldEqual, cslConst.StsSentOut)
				So(recvSupp.SuppEmpID, ShouldEqual, "7000028260")
				So(recvSupp.SuppEmpName, ShouldEqual, "史妍珣")
				So(recvSupp.RecvSuppMst.InUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppDtl.InUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppMst.Dates, ShouldEqual, "20190814")
				So(recvSupp.ShopSuppRecvDate, ShouldEqual, "20190814")
				So(recvSupp.InvtBaseDate, ShouldEqual, "20190814")
				if recvSupp.ProdCode == "Q3AFAFDU6S2100230" {
					So(recvSupp.RecvSuppQty, ShouldEqual, 2)
				}
				if recvSupp.ProdCode == "Q3AFAFDU6S2100240" {
					So(recvSupp.RecvSuppQty, ShouldEqual, 3)
				}
			}
			rtwAnytimeSyncedShouldBeTrueOrFalse("CEGP", waybillNo, true)
		})

		brandCode := "SA"
		shipLocCode := "CFGY"
		waybillNo = "A20190910001"
		title := fmt.Sprintf("%v品牌，%v卖场，运单号为：%v的运单应该同步失败并且在error表中有记录", brandCode, shipLocCode, waybillNo)
		Convey(title, func() {
			has, distError, err := repositories.ReturnToWarehouseErrorRepository{}.GetByWaybillNo("XX", shipLocCode, waybillNo)
			So(has, ShouldEqual, true)
			So(err, ShouldBeNil)
			So(distError.Type, ShouldEqual, clrConst.TypReturnToWarehouseError)
			So(len(distError.ErrorMessage), ShouldBeGreaterThan, 0)
		})
	})
}