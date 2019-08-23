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

// TesRTWtMarkWaybillSynced 退仓单标记为已同步
func TesRTWtMarkWaybillSynced(t *testing.T) {
	Convey("CEGP卖场的20190819001运单，应该标记为已同步", t, func() {
		shipmentLocationCode := "CEGP"
		waybillNo := "1010590009008"
		err := ReturnToWarehouseRepository{}.MarkWaybillSynced(shipmentLocationCode, waybillNo)
		So(err, ShouldBeNil)
		sql := `
			SELECT
				rtw.synced,
				rtw.last_updated_at
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
		So(distList[0]["synced"], ShouldEqual, "1")
		now := time.Now()
		utc, _ := time.LoadLocation("")
		So(distList[0]["last_updated_at"], ShouldEqual, now.In(utc).Format("2006-01-02T15:04:05Z"))
	})
}

func TestGetUnsyncedReturnToWarehouseOrders(t *testing.T) {
	Convey("测试GetUnsyncedReturnToWarehouseOrders", t, func() {
		Convey("应该返回[]map[string]string类型的结果, 并且包含brand_code, shipment_location_code, waybill_no, status_code, emp_id, sku_code, qty, out_date字段", func() {
			result, err := ReturnToWarehouseRepository{}.GetUnsyncedReturnToWarehouseOrders()
			So(err, ShouldBeNil)
			So(len(result), ShouldEqual, 4)
			for _, item := range result {
				requiredKeys := [8]string{"brand_code", "shipment_location_code", "waybill_no", "status_code", "emp_id", "sku_code", "qty", "out_date"}
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
