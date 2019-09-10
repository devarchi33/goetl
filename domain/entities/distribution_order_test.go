package entities

import (
	p2bConst "clearance-adapter/domain/p2brand-constants"
	"clearance-adapter/models"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateDistributionOrder(t *testing.T) {
	Convey("测试DistributionOrder.Create(), 根据map类型的数据转换成DistributionOrder", t, func() {
		Convey("如果源数据类型正确，应该能成确转换", func() {
			data := []map[string]string{
				map[string]string{
					"brand_code":            "SA",
					"receipt_location_code": "CEGP",
					"waybill_no":            "1010590009008",
					"box_no":                "1010590009008",
					"sku_code":              "SPWJ948S2255070",
					"qty":                   "2",
					"in_emp_id":             "7000028260",
					"in_date":               "2019-08-18T19:03:13Z",
					"type":                  "DIRECT_DISTRIBUTION",
				}, map[string]string{
					"brand_code":            "SA",
					"receipt_location_code": "CEGP",
					"waybill_no":            "1010590009008",
					"box_no":                "1010590009008",
					"sku_code":              "SPWJ948S2256070",
					"qty":                   "3",
					"in_emp_id":             "7000028260",
					"in_date":               "2019-08-18T19:03:13Z",
					"type":                  "DIRECT_DISTRIBUTION",
				},
			}
			order, err := DistributionOrder{}.Create(data)
			So(err, ShouldBeNil)
			So(order.BrandCode, ShouldEqual, "SA")
			So(order.ReceiptLocationCode, ShouldEqual, "CEGP")
			So(order.WaybillNo, ShouldEqual, "1010590009008")
			So(order.BoxNo, ShouldEqual, "1010590009008")
			So(order.InEmpID, ShouldEqual, "7000028260")
			So(order.Type, ShouldEqual, p2bConst.TypFactoryToShop)
			So(order.Items, ShouldNotBeNil)
			So(len(order.Items), ShouldEqual, 2)
			So(order.Items[0].SkuCode, ShouldEqual, "SPWJ948S2255070")
			So(order.Items[0].Qty, ShouldEqual, 2)
			So(order.Items[1].SkuCode, ShouldEqual, "SPWJ948S2256070")
			So(order.Items[1].Qty, ShouldEqual, 3)
		})
		title := fmt.Sprintf("如果源数据中如果不包含%v其中的一个或多个key，应该转换失败", strings.Join(DistributionOrder{}.RequiredKeys(), ", "))
		Convey(title, func() {
			data := []map[string]string{
				map[string]string{
					"brand_code": "SA",
					"waybill_no": "1010590009008",
					"box_no":     "1010590009008",
					"sku_code":   "SPWJ948S2255070",
					"qty":        "2",
					"in_emp_id":  "7000028260",
					"in_date":    "2019-08-18T19:03:13Z",
				},
			}
			_, err := DistributionOrder{}.Create(data)
			So(err, ShouldNotBeNil)
			So(err, ShouldResemble, errors.New("receipt_location_code is required"))
		})

		title = fmt.Sprintf("如果源数据中如果不包含%v其中的一个或多个key，应该转换失败", strings.Join(DistributionOrderItem{}.RequiredKeys(), ", "))
		Convey(title, func() {
			data := []map[string]string{
				map[string]string{
					"brand_code":            "SA",
					"receipt_location_code": "CEGP",
					"waybill_no":            "1010590009008",
					"box_no":                "1010590009008",
					"sku_code":              "SPWJ948S2255070",
					"in_emp_id":             "7000028260",
					"in_date":               "2019-08-18T19:03:13Z",
				},
			}
			_, err := DistributionOrder{}.Create(data)
			So(err, ShouldNotBeNil)
			So(err, ShouldResemble, errors.New("qty is required"))
		})
	})
}

func TestCreateByRecvSupp(t *testing.T) {
	brandCode := "SA"
	receiptLocationCode := "CEGP"
	waybillNo := "20190902001"
	boxNo := "20190902001"
	sku1 := "SPWJ948S2255070"
	qty1 := 1
	sku2 := "SPWJ948S2256070"
	qty2 := 2
	Convey("如果源数据类型正确，应该能成确转换", t, func() {
		recvSupp1 := models.RecvSupp{
			RecvSuppMst: models.RecvSuppMst{
				BrandCode:        brandCode,
				ShopCode:         receiptLocationCode,
				WayBillNo:        waybillNo,
				BoxNo:            boxNo,
				ShippingTypeCode: "16",
			},
			RecvSuppDtl: models.RecvSuppDtl{
				ProdCode:    sku1,
				RecvSuppQty: qty1,
			},
		}
		recvSupp2 := models.RecvSupp{
			RecvSuppMst: models.RecvSuppMst{
				BrandCode:        brandCode,
				ShopCode:         receiptLocationCode,
				WayBillNo:        waybillNo,
				BoxNo:            boxNo,
				ShippingTypeCode: "16",
			},
			RecvSuppDtl: models.RecvSuppDtl{
				ProdCode:    sku2,
				RecvSuppQty: qty2,
			},
		}
		data := []models.RecvSupp{recvSupp1, recvSupp2}
		order, err := DistributionOrder{}.CreateByRecvSupp(data)
		So(err, ShouldBeNil)
		So(order.BrandCode, ShouldEqual, brandCode)
		So(order.ReceiptLocationCode, ShouldEqual, receiptLocationCode)
		So(order.WaybillNo, ShouldEqual, waybillNo)
		So(order.BoxNo, ShouldEqual, boxNo)
		So(order.Type, ShouldEqual, p2bConst.TypFactoryToShop)
		So(order.Version, ShouldEqual, strconv.Itoa(len(data)))
		So(order.Items[0].SkuCode, ShouldEqual, sku1)
		So(order.Items[0].Qty, ShouldEqual, qty1)
		So(order.Items[1].SkuCode, ShouldEqual, sku2)
		So(order.Items[1].Qty, ShouldEqual, qty2)
	})
}
