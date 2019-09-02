package repositories

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"clearance-adapter/domain/entities"
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	_ "clearance-adapter/test"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetUnsyncedDistributionOrders(t *testing.T) {
	Convey("测试GetUnsyncedDistributionOrders", t, func() {
		requiredKeys := append(entities.DistributionOrder{}.RequiredKeys(), entities.DistributionOrderItem{}.RequiredKeys()...)
		title := fmt.Sprintf("应该返回[]map[string]string类型的结果, 并且包含%v字段", strings.Join(requiredKeys, " ,"))
		Convey(title, func() {
			result, err := StockDistributionRepository{}.GetUnsyncedDistributionOrders()
			So(err, ShouldBeNil)
			for _, item := range result {
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

func TestPutInStorage(t *testing.T) {
	order := entities.DistributionOrder{
		BrandCode:           "SA",
		ReceiptLocationCode: "CEGP",
		WaybillNo:           "20190829003",
		BoxNo:               "20190829003",
		Version:             "1",
		Items:               make([]entities.DistributionOrderItem, 0),
	}
	order.Items = append(order.Items, entities.DistributionOrderItem{
		SkuCode: "SPYC949S1139085",
		Qty:     1,
	})
	title := fmt.Sprintf("应该在%v卖场，生成运单号为%v的入库单", order.ReceiptLocationCode, order.WaybillNo)
	Convey(title, t, func() {
		err := StockDistributionRepository{}.PutInStorage(order)
		So(err, ShouldBeNil)
	})
}
