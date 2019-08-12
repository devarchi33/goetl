package repositories

import (
	"errors"
	"testing"
	"time"

	_ "clearance-adapter/test"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func checkRequirement(data map[string]string, requiredKeys ...string) error {
	for _, key := range requiredKeys {
		if _, ok := data[key]; !ok {
			return errors.New(key + " is required")
		}
	}
	return nil
}

func TestGetInStorageByCreateAt(t *testing.T) {
	Convey("测试GetInStorageByCreateAt", t, func() {
		Convey("应该返回[]map[string]string类型的结果, 并且包含brand_code, shop_code, waybill_no, box_no, emp_id, sku_code, qty字段", func() {
			local, _ := time.LoadLocation("Local")
			start, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-07-01 00:00:00", local)
			end, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-07-31 23:59:59", local)
			result, err := StockDistributionRepository{}.GetInStorageByCreateAt(start, end)
			So(err, ShouldBeNil)
			for _, item := range result {
				err := checkRequirement(item, "brand_code", "shop_code", "waybill_no", "box_no", "emp_id", "sku_code", "qty")
				So(err, ShouldBeNil)
			}
		})
	})
}
