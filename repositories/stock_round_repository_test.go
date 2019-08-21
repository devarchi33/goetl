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

// TesMarkTransferOrderSynced 标记调货入库单为已同步
func TesMarkTransferOrderSynced(t *testing.T) {
	Convey("CEGP卖场到CFGY卖场的运单号为20190821001的调货入库单，应该标记为已同步", t, func() {
		shipmentLocationCode := "CEGP"
		receiptLocationCode := "CFGY"
		waybillNo := "20190821001"
		err := StockRoundRepository{}.MarkWaybillSynced(shipmentLocationCode, receiptLocationCode, waybillNo)
		So(err, ShouldBeNil)
		sql := `
			SELECT
				sr.synced,
				sr.last_updated_at
			FROM pangpang_brand_sku_location.stock_round sr
			JOIN pangpang_brand_place_management.store AS shipmentStore
				ON shipmentStore.id = sr.shipment_location_id
			JOIN pangpang_brand_place_management.store AS receitpStore
				ON receitpStore.id = sr.receipt_location_id
			WHERE sr.tenant_code = 'pangpang'
				AND shipmentStore.code = ?
				AND receitpStore.code = ?
				AND sr.waybill_no = ?
			;
		`
		result, err := factory.GetP2BrandEngine().Query(sql, shipmentLocationCode, receiptLocationCode, waybillNo)
		So(err, ShouldBeNil)
		So(len(result), ShouldEqual, 1)
		distList := infra.ConvertByteResult(result)
		So(distList[0]["synced"], ShouldEqual, "1")
		now := time.Now()
		utc, _ := time.LoadLocation("")
		So(distList[0]["last_updated_at"], ShouldEqual, now.In(utc).Format("2006-01-02T15:04:05Z"))
	})
}

func TestGetUnsyncedTransferInOrders(t *testing.T) {
	Convey("测试GetUnsyncedTransferInOrders", t, func() {
		Convey("应该返回[]map[string]string类型的结果, 并且包含brand_code, shipment_location_code, receipt_location_code, waybill_no, box_no, out_date, out_emp_id, in_date, in_emp_id, sku_code, qty, shipping_company_code字段", func() {
			result, err := StockRoundRepository{}.GetUnsyncedTransferInOrders()
			So(err, ShouldBeNil)
			So(len(result), ShouldEqual, 2)
			for _, item := range result {
				requiredKeys := [12]string{"brand_code", "shipment_location_code", "receipt_location_code", "waybill_no", "box_no", "out_date", "out_emp_id", "in_date", "in_emp_id", "sku_code", "qty", "shipping_company_code"}
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
