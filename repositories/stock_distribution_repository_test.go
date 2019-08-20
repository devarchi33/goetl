package repositories

import (
	"testing"
	"time"

	"clearance-adapter/factory"
	"clearance-adapter/infra"
	_ "clearance-adapter/test"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetUnsyncedDistributionOrders(t *testing.T) {
	Convey("测试GetUnsyncedDistributionOrders", t, func() {
		Convey("应该返回[]map[string]string类型的结果, 并且包含brand_code, shop_code, waybill_no, box_no, emp_id, sku_code, qty字段", func() {
			result, err := StockDistributionRepository{}.GetUnsyncedDistributionOrders()
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

func TestMarkWaybillSynced(t *testing.T) {
	Convey("CEGP卖场的1010590009008运单，应该标记为已同步", t, func() {
		receiptLocationCode := "CEGP"
		waybillNo := "1010590009008"
		err := StockDistributionRepository{}.MarkWaybillSynced(receiptLocationCode, waybillNo)
		So(err, ShouldBeNil)
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
		So(err, ShouldBeNil)
		So(len(result), ShouldEqual, 1)
		distList := infra.ConvertByteResult(result)
		So(distList[0]["synced"], ShouldEqual, "1")
		now := time.Now()
		utc, _ := time.LoadLocation("")
		So(distList[0]["last_updated_at"], ShouldEqual, now.In(utc).Format("2006-01-02T15:04:05Z"))
	})
}
