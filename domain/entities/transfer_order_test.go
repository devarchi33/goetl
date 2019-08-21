package entities

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateTransferOrder(t *testing.T) {
	Convey("测试TransferOrder.Create(), 根据map类型的数据转换成TransferOrder", t, func() {
		Convey("如果源数据类型正确，应该能成确转换", func() {
			data := []map[string]string{
				map[string]string{
					"brand_code":             "SA",
					"shipment_location_code": "CEGP",
					"receipt_location_code":  "CFGY",
					"waybill_no":             "20190821001",
					"box_no":                 "20190821001",
					"sku_code":               "SPWJ948S2255070",
					"qty":                    "2",
					"out_date":               "2019-08-20T20:02:09Z",
					"out_emp_id":             "7000000001",
					"in_date":                "2019-08-20T22:02:09Z",
					"in_emp_id":              "7000000002",
					"shipping_company_code":  "SR",
				}, map[string]string{
					"brand_code":             "SA",
					"shipment_location_code": "CEGP",
					"receipt_location_code":  "CFGY",
					"waybill_no":             "20190821001",
					"box_no":                 "20190821001",
					"sku_code":               "SPWJ948S2256070",
					"qty":                    "3",
					"out_date":               "2019-08-20T20:02:09Z",
					"out_emp_id":             "7000000001",
					"in_date":                "2019-08-20T22:02:09Z",
					"in_emp_id":              "7000000002",
					"shipping_company_code":  "SR",
				},
			}
			order, err := TransferOrder{}.Create(data)
			So(err, ShouldBeNil)
			So(order.BrandCode, ShouldEqual, "SA")
			So(order.ShipmentLocationCode, ShouldEqual, "CEGP")
			So(order.ReceiptLocationCode, ShouldEqual, "CFGY")
			So(order.WaybillNo, ShouldEqual, "20190821001")
			So(order.BoxNo, ShouldEqual, "20190821001")
			So(order.OutDate, ShouldEqual, "20190821")
			So(order.OutEmpID, ShouldEqual, "7000000001")
			So(order.InDate, ShouldEqual, "20190821")
			So(order.InEmpID, ShouldEqual, "7000000002")
			So(order.ShippingCompanyCode, ShouldEqual, "SR")
			So(order.Items, ShouldNotBeNil)
			So(len(order.Items), ShouldEqual, 2)
			So(order.Items[0].SkuCode, ShouldEqual, "SPWJ948S2255070")
			So(order.Items[0].Qty, ShouldEqual, 2)
			So(order.Items[1].SkuCode, ShouldEqual, "SPWJ948S2256070")
			So(order.Items[1].Qty, ShouldEqual, 3)
		})
	})
}
