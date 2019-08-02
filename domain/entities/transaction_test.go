package entities

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateTransaction(t *testing.T) {
	Convey("测试Transaction.Create(), 根据map类型的数据转换成Transaction", t, func() {
		Convey("如果源数据类型正确，应该能成确转换", func() {
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
			}
			transaction, err := Transaction{}.Create(data)
			So(err, ShouldBeNil)
			So(transaction.BrandCode, ShouldEqual, "SA")
			So(transaction.ShopCode, ShouldEqual, "CEGP")
			So(transaction.WaybillNo, ShouldEqual, "1010590009008")
			So(transaction.BoxNo, ShouldEqual, "1010590009008")
			So(transaction.EmpID, ShouldEqual, "7000028260")
			So(transaction.Items, ShouldNotBeNil)
			So(len(transaction.Items), ShouldEqual, 2)
			So(transaction.Items[0].SkuCode, ShouldEqual, "SPWJ948S2255070")
			So(transaction.Items[0].Qty, ShouldEqual, 2)
			So(transaction.Items[1].SkuCode, ShouldEqual, "SPWJ948S2256070")
			So(transaction.Items[1].Qty, ShouldEqual, 3)
		})
		Convey("如果源数据中如果不包含“brand_code, shop_code, waybill_no, box_no, emp_id其中的一个或多个key，应该转换失败", func() {
			data := []map[string]string{
				map[string]string{
					"brand_code": "SA",
					"waybill_no": "1010590009008",
					"box_no":     "1010590009008",
					"sku_code":   "SPWJ948S2255070",
					"qty":        "2",
					"emp_id":     "7000028260",
				},
			}
			_, err := Transaction{}.Create(data)
			So(err, ShouldNotBeNil)
			So(err, ShouldResemble, errors.New("shop_code is required"))
		})
		Convey("如果源数据中如果不包含“sku_code, qty其中的一个或多个key，应该转换失败", func() {
			data := []map[string]string{
				map[string]string{
					"brand_code": "SA",
					"shop_code":  "CEGP",
					"waybill_no": "1010590009008",
					"box_no":     "1010590009008",
					"sku_code":   "SPWJ948S2255070",
					"emp_id":     "7000028260",
				},
			}
			_, err := Transaction{}.Create(data)
			So(err, ShouldNotBeNil)
			So(err, ShouldResemble, errors.New("qty is required"))
		})
	})
}
