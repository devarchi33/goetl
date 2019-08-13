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

func TestTransform(t *testing.T) {
	Convey("测试Transform方法", t, func() {
		Convey("源数据中包含多个运单号的数据，应该能够根据运单号生成ReturnToWarehouseOrder", func() {
			data := []map[string]string{
				map[string]string{
					"brand_code":             "SA",
					"shipment_location_code": "CEGP",
					"waybill_no":             "1010590009008",
					"status_code":            "R",
					"sku_code":               "SPWJ948S2255070",
					"qty":                    "2",
					"emp_id":                 "7000028260",
				}, map[string]string{
					"brand_code":             "SA",
					"shipment_location_code": "CEGP",
					"waybill_no":             "1010590009008",
					"status_code":            "R",
					"sku_code":               "SPWJ948S2256070",
					"qty":                    "3",
					"emp_id":                 "7000028260",
				},
				map[string]string{
					"brand_code":             "SA",
					"shipment_location_code": "CEGP",
					"waybill_no":             "1010590009009",
					"status_code":            "R",
					"sku_code":               "SPYC949H2130095",
					"qty":                    "4",
					"emp_id":                 "7000028260",
				}, map[string]string{
					"brand_code":             "SA",
					"shipment_location_code": "CEGP",
					"waybill_no":             "1010590009009",
					"status_code":            "R",
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
					So(order.StatusCode, ShouldEqual, "R")
					So(order.EmpID, ShouldEqual, "7000028260")
					So(order.Items, ShouldNotBeNil)
					So(len(order.Items), ShouldEqual, 2)
					So(order.Items[0].SkuCode, ShouldEqual, "SPWJ948S2255070")
					So(order.Items[0].Qty, ShouldEqual, 2)
					So(order.Items[1].SkuCode, ShouldEqual, "SPWJ948S2256070")
					So(order.Items[1].Qty, ShouldEqual, 3)
				} else if order.BrandCode == "SA" && order.ShipmentLocationCode == "CEGP" && order.WaybillNo == "1010590009009" {
					So(order.StatusCode, ShouldEqual, "R")
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

func TestReturnToWarehouseETL(t *testing.T) {
	Convey("测试ReturnToWarehouseETL的Run方法", t, func() {
		Convey("某个时间段没有退仓出库单的话，应该没有数据在CSL退仓记录", func() {
			etl := ReturnToWarehouseETL{}.New("2019-07-01 00:00:00", "2019-07-31 00:01:00")
			err := etl.Run(context.Background())
			So(err, ShouldBeNil)
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", "CEGP", "20190813001")
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 0)
		})
		Convey("运单号为20190813001的退仓出库单应该在CSL有记录", func() {
			etl := ReturnToWarehouseETL{}.New("2019-08-13 00:00:00", "2019-08-13 23:59:59")
			err := etl.Run(context.Background())
			So(err, ShouldBeNil)
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", "CEGP", "20190813001")
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 2)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, false)
				So(recvSupp.ShippingTypeCode, ShouldEqual, "41")
				So(recvSupp.RecvSuppType, ShouldEqual, "S")
				So(recvSupp.RecvSuppStatusCode, ShouldEqual, "R")
				So(recvSupp.SuppEmpID, ShouldEqual, "7000028260")
				So(recvSupp.SuppEmpName, ShouldEqual, "史妍珣")
				So(recvSupp.RecvSuppMst.InUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppDtl.InUserID, ShouldEqual, "shi.yanxun")
			}
		})
	})
}
