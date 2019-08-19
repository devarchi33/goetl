package services

import (
	cslConst "clearance-adapter/domain/csl-constants"
	"clearance-adapter/domain/entities"
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"clearance-adapter/repositories"
	_ "clearance-adapter/test"
	"context"
	"log"
	"testing"

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
				}, map[string]string{
					"brand_code":             "SA",
					"shipment_location_code": "CEGP",
					"waybill_no":             "1010590009008",
					"status_code":            p2bConst.StsSentOut,
					"sku_code":               "SPWJ948S2256070",
					"qty":                    "3",
					"emp_id":                 "7000028260",
				},
				map[string]string{
					"brand_code":             "SA",
					"shipment_location_code": "CEGP",
					"waybill_no":             "1010590009009",
					"status_code":            p2bConst.StsSentOut,
					"sku_code":               "SPYC949H2130095",
					"qty":                    "4",
					"emp_id":                 "7000028260",
				}, map[string]string{
					"brand_code":             "SA",
					"shipment_location_code": "CEGP",
					"waybill_no":             "1010590009009",
					"status_code":            p2bConst.StsSentOut,
					"sku_code":               "SPYC949H2130100",
					"qty":                    "5",
					"emp_id":                 "7000028260",
				},
			}

			result, err := ReturnToWarehouseETL{}.Transform(context.Background(), data)
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

func rtwSyncedShouldBeTrueOrFalse(shipmentLocationCode, waybillNo string, shouldBeTrue bool) {
	sql := `
		SELECT
			rtw.synced
		FROM pangpang_brand_sku_location.return_to_warehouse AS rtw
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = rtw.shipment_location_id
		WHERE rtw.tenant_code = 'pangpang'
			AND store.code = ?
			AND rtw.waybill_no = ?
	`
	result, err := factory.GetP2BrandEngine().Query(sql, shipmentLocationCode, waybillNo)
	So(err, ShouldBeNil)
	So(len(result), ShouldEqual, 1)
	distList := infra.ConvertByteResult(result)
	if shouldBeTrue {
		So(distList[0]["synced"], ShouldEqual, "1")
	} else {
		So(distList[0]["synced"], ShouldEqual, "0")
	}
}

func TestReturnToWarehouseETL(t *testing.T) {
	Convey("测试ReturnToWarehouseETL的Run方法", t, func() {
		etl := ReturnToWarehouseETL{}.New()
		err := etl.Run(context.Background())
		So(err, ShouldBeNil)
		Convey("运单号为20190819001的出库单synced=true，所以应该不会同步到CSL入库", func() {
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", "CEGP", "20190819001")
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 0)
		})
		Convey("运单号为20190813001的退仓出库单应该在CSL有记录", func() {
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", "CEGP", "20190813001")
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 2)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, false)
				So(recvSupp.ShippingTypeCode, ShouldEqual, "41")
				So(recvSupp.RecvSuppType, ShouldEqual, "S")
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
			rtwSyncedShouldBeTrueOrFalse("CEGP", "20190813001", true)
		})

		Convey("CEGP的子卖场CJC1运单号为20190814001的退仓出库单应该在CSL有记录", func() {
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("Q3", "CJC1", "20190814001")
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 2)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, false)
				So(recvSupp.ShippingTypeCode, ShouldEqual, "41")
				So(recvSupp.RecvSuppType, ShouldEqual, "S")
				So(recvSupp.RecvSuppStatusCode, ShouldEqual, cslConst.StsSentOut)
				So(recvSupp.SuppEmpID, ShouldEqual, "7000028260")
				So(recvSupp.SuppEmpName, ShouldEqual, "史妍珣")
				So(recvSupp.RecvSuppMst.InUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppDtl.InUserID, ShouldEqual, "shi.yanxun")
				if recvSupp.ProdCode == "Q3AFAFDU6S2100230" {
					So(recvSupp.RecvSuppQty, ShouldEqual, 2)
				}
				if recvSupp.ProdCode == "Q3AFAFDU6S2100240" {
					So(recvSupp.RecvSuppQty, ShouldEqual, 3)
				}
			}
			rtwSyncedShouldBeTrueOrFalse("CEGP", "20190814001", true)
		})
	})
}
