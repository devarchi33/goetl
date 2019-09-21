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

// TestMarkAnytimeWaybillSynced 随时退仓单标记为已同步
func TestMarkAnytimeWaybillSynced(t *testing.T) {
	Convey("CEGP卖场的A20190813001运单，应该标记为已同步", t, func() {
		shipmentLocationCode := "CEGP"
		waybillNo := "A20190813001"
		err := ReturnToWarehouseRepository{}.MarkAnytimeWaybillSynced(shipmentLocationCode, waybillNo)
		So(err, ShouldBeNil)
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
		So(err, ShouldBeNil)
		So(len(result), ShouldEqual, 1)
		distList := infra.ConvertByteResult(result)
		So(distList[0]["synced"], ShouldEqual, "1")
		now := time.Now()
		utc, _ := time.LoadLocation("")
		So(distList[0]["last_updated_at"], ShouldEqual, now.In(utc).Format("2006-01-02T15:04:05Z"))
	})
}

// TestMarkSeasonalWaybillSynced 季节退仓单标记为已同步
func TestMarkSeasonalWaybillSynced(t *testing.T) {
	Convey("CEGP卖场的S20190813001运单，应该标记为已同步", t, func() {
		shipmentLocationCode := "CEGP"
		waybillNo := "S20190813001"
		err := ReturnToWarehouseRepository{}.MarkSeasonalWaybillSynced(shipmentLocationCode, waybillNo)
		So(err, ShouldBeNil)
		sql := `
			SELECT
				rtw.synced,
				rtw.last_updated_at
			FROM pangpang_brand_sku_location.return_to_warehouse_seasonal AS rtw
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

// TestMarkDefectiveWaybillSynced 随时退仓单标记为已同步
func TestMarkDefectiveWaybillSynced(t *testing.T) {
	Convey("CEGP卖场的D20190813001运单，应该标记为已同步", t, func() {
		shipmentLocationCode := "CEGP"
		waybillNo := "D20190813001"
		err := ReturnToWarehouseRepository{}.MarkDefectiveWaybillSynced(shipmentLocationCode, waybillNo)
		So(err, ShouldBeNil)
		sql := `
			SELECT
				rtw.synced,
				rtw.last_updated_at
			FROM pangpang_brand_sku_location.return_to_warehouse_defective AS rtw
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

func TestGetUnsyncedReturnToWarehouseAnytimeOrders(t *testing.T) {
	Convey("测试TestGetUnsyncedReturnToWarehouseAnytimeOrders", t, func() {
		Convey("应该返回[]map[string]string类型的结果, 并且包含brand_code, shipment_location_code, waybill_no, status_code, emp_id, sku_code, qty, out_date字段", func() {
			result, err := ReturnToWarehouseRepository{}.GetUnsyncedReturnToWarehouseAnytimeOrders()
			So(err, ShouldBeNil)
			So(len(result), ShouldEqual, 3)
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

func TestGetUnsyncedReturnToWarehouseSeasonalOrders(t *testing.T) {
	Convey("测试TestGetUnsyncedReturnToWarehouseSeasonalOrders", t, func() {
		Convey("应该返回[]map[string]string类型的结果, 并且包含brand_code, shipment_location_code, waybill_no, status_code, emp_id, sku_code, qty, out_date字段", func() {
			result, err := ReturnToWarehouseRepository{}.GetUnsyncedReturnToWarehouseSeasonalOrders()
			So(err, ShouldBeNil)
			So(len(result), ShouldEqual, 3)
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

func TestGetUnsyncedReturnToWarehouseDefectiveOrders(t *testing.T) {
	Convey("测试TestGetUnsyncedReturnToWarehouseDefectiveOrders", t, func() {
		Convey("应该返回[]map[string]string类型的结果, 并且包含brand_code, shipment_location_code, waybill_no, status_code, emp_id, sku_code, qty, out_date字段", func() {
			result, err := ReturnToWarehouseRepository{}.GetUnsyncedReturnToWarehouseDefectiveOrders()
			So(err, ShouldBeNil)
			So(len(result), ShouldEqual, 3)
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
