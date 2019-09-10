package services

import (
	"clearance-adapter/repositories"
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAutoDistributionETL(t *testing.T) {
	etl := AutoDistributionETL{}.New()
	etl.Run(context.Background())
	waybillNo := "20190902001"
	brandCode := "Q3"
	shopCode := "CJC1"
	receiptLocationCode := "CEGP"
	title := fmt.Sprintf("复合卖场子卖场%v(品牌%v)的%v运单应该在P2Brand中入库，并且ReceiptLocationCode为%v", shopCode, brandCode, waybillNo, receiptLocationCode)
	Convey(title, t, func() {
		orders, err := repositories.StockDistributionRepository{}.GetUnsyncedDistributionOrders()
		So(err, ShouldBeNil)

		has := false
		for _, v := range orders {
			if v["waybill_no"] == waybillNo {
				has = true
				So(v["brand_code"], ShouldEqual, brandCode)
				So(v["receipt_location_code"], ShouldEqual, receiptLocationCode)
			}
		}
		So(has, ShouldEqual, true)
	})

	waybillNo = "20190902002"
	brandCode = "SA"
	shopCode = "CJUT"
	title = fmt.Sprintf("非复合卖场卖场%v(品牌%v)的%v运单应该在P2Brand中入库，并且ReceiptLocationCode为%v", shopCode, brandCode, waybillNo, shopCode)
	Convey(title, t, func() {
		orders, err := repositories.StockDistributionRepository{}.GetUnsyncedDistributionOrders()
		So(err, ShouldBeNil)
		has := false
		for _, v := range orders {
			if v["waybill_no"] == waybillNo {
				has = true
				So(v["brand_code"], ShouldEqual, brandCode)
				So(v["receipt_location_code"], ShouldEqual, shopCode)
			}
		}
		So(has, ShouldEqual, true)
	})

	waybillNo = "20190909002"
	brandCode = "SA"
	shopCode = "CEGP"
	title = fmt.Sprintf("【工厂直送】卖场%v(品牌%v)的%v运单应该在P2Brand中入库，并且ReceiptLocationCode为%v", shopCode, brandCode, waybillNo, shopCode)
	Convey(title, t, func() {
		orders, err := repositories.DirectDistributionRepository{}.GetUnsyncedDistributionOrders()
		So(err, ShouldBeNil)
		has := false
		for _, v := range orders {
			if v["waybill_no"] == waybillNo {
				has = true
				So(v["brand_code"], ShouldEqual, brandCode)
				So(v["receipt_location_code"], ShouldEqual, shopCode)
			}
		}
		So(has, ShouldEqual, true)
	})
}
