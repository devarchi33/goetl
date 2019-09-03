package services

import (
	"clearance-adapter/repositories"
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAutoTransferETL(t *testing.T) {
	etl := AutoTransferETL{}.New()
	etl.Run(context.Background())
	waybillNo := "20190903001"
	brandCode := "Q3"
	shopCode := "CJC1"
	targetShopCode := "CHTM"
	shipmentLocationCode := "CEGP"
	receiptLocationCode := "CFGY"
	title := fmt.Sprintf("复合卖场子卖场%v(品牌%v)调到%v卖场的%v运单应该在P2Brand中入库，并且ShipmentLocationCode为%v, ReceiptLocationCode为%v", shopCode, brandCode, targetShopCode, waybillNo, shipmentLocationCode, receiptLocationCode)
	Convey(title, t, func() {
		orders, err := repositories.StockRoundRepository{}.GetUnsyncedTransferInOrders()
		So(err, ShouldBeNil)

		has := false
		for _, v := range orders {
			if v["waybill_no"] == waybillNo {
				has = true
				So(v["brand_code"], ShouldEqual, brandCode)
				So(v["shipment_location_code"], ShouldEqual, shipmentLocationCode)
				So(v["receipt_location_code"], ShouldEqual, receiptLocationCode)
			}
		}
		So(has, ShouldEqual, true)
	})

	waybillNo = "20190903002"
	brandCode = "SA"
	shopCode = "CJUT"
	targetShopCode = "CEGP"
	title = fmt.Sprintf("非复合卖场子卖场%v(品牌%v)调到%v卖场的%v运单应该在P2Brand中入库，并且ShipmentLocationCode为%v, ReceiptLocationCode为%v", shopCode, brandCode, targetShopCode, waybillNo, shipmentLocationCode, receiptLocationCode)
	Convey(title, t, func() {
		orders, err := repositories.StockRoundRepository{}.GetUnsyncedTransferInOrders()
		So(err, ShouldBeNil)
		has := false
		for _, v := range orders {
			if v["waybill_no"] == waybillNo {
				has = true
				So(v["brand_code"], ShouldEqual, brandCode)
				So(v["shipment_location_code"], ShouldEqual, shopCode)
				So(v["receipt_location_code"], ShouldEqual, targetShopCode)
			}
		}
		So(has, ShouldEqual, true)
	})
}
