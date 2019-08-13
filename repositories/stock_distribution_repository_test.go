package repositories

import (
	"testing"
	"time"

	_ "clearance-adapter/test"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetDistributionOrdersByCreateAt(t *testing.T) {
	Convey("测试GetDistributionOrdersByCreateAt", t, func() {
		Convey("应该返回[]map[string]string类型的结果, 并且包含brand_code, shop_code, waybill_no, box_no, emp_id, sku_code, qty字段", func() {
			local, _ := time.LoadLocation("Local")
			start, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-07-01 00:00:00", local)
			end, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-07-31 23:59:59", local)
			result, err := StockDistributionRepository{}.GetDistributionOrdersByCreateAt(start, end)
			So(err, ShouldBeNil)
			for _, item := range result {
				requiredKeys := [7]string{"brand_code", "shop_code", "waybill_no", "box_no", "emp_id", "sku_code", "qty"}
				isOk := true
				for _, key := range requiredKeys {
					_, ok := item[key]
					isOk = isOk && ok
				}
				So(isOk, ShouldEqual, true)
			}
		})
	})
}
