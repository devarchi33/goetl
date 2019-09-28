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

	p2bConst "clearance-adapter/domain/p2brand-constants"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetUnsyncedDirectDistributionOrders(t *testing.T) {
	Convey("测试GetUnsyncedDistributionOrders", t, func() {
		requiredKeys := append(entities.DistributionOrder{}.RequiredKeys(), entities.DistributionOrderItem{}.RequiredKeys()...)
		title := fmt.Sprintf("应该返回[]map[string]string类型的结果, 并且包含%v字段", strings.Join(requiredKeys, " ,"))
		Convey(title, func() {
			result, err := DirectDistributionRepository{}.GetUnsyncedDistributionOrders()
			So(err, ShouldBeNil)
			for _, item := range result {
				isOk := true
				for _, key := range requiredKeys {
					_, ok := item[key]
					isOk = isOk && ok
				}
				So(isOk, ShouldEqual, true)
				So(item["type"], ShouldEqual, p2bConst.TypFactoryToShop)
			}
		})
	})
}

func TestDirectMarkWaybillSynced(t *testing.T) {
	recptLocCode := "CEGP"
	waybillNo := "20190909001"
	title := fmt.Sprintf("%v卖场的%v运单，应该标记为已同步", recptLocCode, waybillNo)
	Convey(title, t, func() {
		err := DirectDistributionRepository{}.MarkWaybillSynced(recptLocCode, waybillNo)
		So(err, ShouldBeNil)
		sql := `
			SELECT
				dd.synced,
				dd.last_updated_at
			FROM pangpang_brand_sku_location.direct_distribution AS dd
				JOIN pangpang_brand_place_management.store AS store
					ON store.id = dd.receipt_location_id
			WHERE dd.tenant_code = 'pangpang'
				AND store.code = ?
				AND dd.waybill_no = ?
		`
		result, err := factory.GetP2BrandEngine().Query(sql, recptLocCode, waybillNo)
		So(err, ShouldBeNil)
		So(len(result), ShouldEqual, 1)
		distList := infra.ConvertByteResult(result)
		So(distList[0]["synced"], ShouldEqual, "1")
		now := time.Now()
		utc, _ := time.LoadLocation("")
		So(distList[0]["last_updated_at"], ShouldEqual, now.In(utc).Format("2006-01-02T15:04:05Z"))
	})
}

func TestDirectPutInStorage(t *testing.T) {
	order := entities.DistributionOrder{
		BrandCode:           "SA",
		ReceiptLocationCode: "CEGP",
		WaybillNo:           "20190909001",
		BoxNo:               "20190909001",
		Version:             "1",
		Items:               make([]entities.DistributionOrderItem, 0),
	}
	order.Items = append(order.Items, entities.DistributionOrderItem{
		SkuCode: "SPWJ948S2255070",
		Qty:     4,
	}, entities.DistributionOrderItem{
		SkuCode: "SPWJ948S2255075",
		Qty:     3,
	})

	title := fmt.Sprintf("应该在%v卖场，生成运单号为%v的入库单", order.ReceiptLocationCode, order.WaybillNo)
	Convey(title, t, func() {
		err := DirectDistributionRepository{}.PutInStorage(order, true)
		So(err, ShouldBeNil)
	})
}
