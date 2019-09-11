package services

import (
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"clearance-adapter/repositories"
	_ "clearance-adapter/test"
	"context"
	"fmt"
	"testing"
	"time"

	clrConst "clearance-adapter/domain/clr-constants"
	cslConst "clearance-adapter/domain/csl-constants"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func transferOrderSyncedShouldBeTrue(shipmentLocationCode, receiptLocationCode, waybillNo string) {
	utc, _ := time.LoadLocation("")
	startDate := time.Now().Add(-1000).In(utc).Format("2006-01-02T15:04:05Z")
	sql := `
		SELECT
			sr.synced,
			sr.last_updated_at
		FROM pangpang_brand_sku_location.stock_round AS sr
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
	endDate := time.Now().In(utc).Format("2006-01-02T15:04:05Z")
	So(err, ShouldBeNil)
	So(len(result), ShouldEqual, 1)
	tranList := infra.ConvertByteResult(result)
	So(tranList[0]["synced"], ShouldEqual, "1")
	So(tranList[0]["last_updated_at"], ShouldBeGreaterThanOrEqualTo, startDate)
	So(tranList[0]["last_updated_at"], ShouldBeLessThanOrEqualTo, endDate)
}

func TestTransferETL(t *testing.T) {
	Convey("测试TransferETL的Run方法", t, func() {
		etl := TransferETL{}.New()
		err := etl.Run(context.Background())
		So(err, ShouldBeNil)
		brandCode := "SA"
		shipmentLocationCode := "CEGP"
		receiptLocationCode := "CFGY"
		waybillNo := "20190821001"

		Convey("CEGP卖场调到CFGY卖场的运单为20190821001的调货单同步后，应该在CSL有调出数据，出库日期为20190821", func() {
			recvSupp, err := repositories.RecvSuppRepository{}.GetByWaybillNo(brandCode, shipmentLocationCode, waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSupp), ShouldEqual, 2)
			for _, v := range recvSupp {
				So(v.RecvSuppStatusCode, ShouldEqual, cslConst.StsConfirmed)
				So(v.RecvSuppMst.Dates, ShouldEqual, "20190821")
				So(v.RecvSuppDtl.Dates, ShouldEqual, "20190821")
				So(v.ShopSuppRecvDate, ShouldEqual, "20190821")
				if v.ProdCode == "SPWJ948S2255070" {
					So(v.RecvSuppQty, ShouldEqual, 4)
				}
				if v.ProdCode == "SPYC949S1139095" {
					So(v.RecvSuppQty, ShouldEqual, 1)
				}
			}
		})
		Convey("CEGP卖场调到CFGY卖场的运单为20190821001的调货单同步后，应该在CSL有调入数据，入库日期为20190822", func() {
			recvSupp, err := repositories.RecvSuppRepository{}.GetByWaybillNo(brandCode, receiptLocationCode, waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSupp), ShouldEqual, 2)
			for _, v := range recvSupp {
				So(v.RecvSuppStatusCode, ShouldEqual, cslConst.StsConfirmed)
				So(v.RecvSuppMst.Dates, ShouldEqual, "20190822")
				So(v.RecvSuppDtl.Dates, ShouldEqual, "20190822")
				So(v.ShopSuppRecvDate, ShouldEqual, "20190822")
				if v.ProdCode == "SPWJ948S2255070" {
					So(v.RecvSuppFixedQty, ShouldEqual, 4)
				}
				if v.ProdCode == "SPYC949S1139095" {
					So(v.RecvSuppFixedQty, ShouldEqual, 1)
				}
			}
		})
		Convey("CEGP卖场调到CFGY卖场的运单为20190821001的调货单同步后，synced应该为true", func() {
			transferOrderSyncedShouldBeTrue(shipmentLocationCode, receiptLocationCode, waybillNo)
		})

		brandCode = "SA"
		shipLocCode := "CEGP"
		recptLocCode := "CFGY"
		waybillNo = "20190911001"
		title := fmt.Sprintf("%v品牌，从%v卖场调货到%v卖场，运单号为：%v的运单应该同步失败并且在error表中有记录", brandCode, shipLocCode, recptLocCode, waybillNo)
		Convey(title, func() {
			has, tranError, err := repositories.StockRoundErrorRepository{}.GetByWaybillNo("XX", shipLocCode, recptLocCode, waybillNo)
			So(has, ShouldEqual, true)
			So(err, ShouldBeNil)
			So(tranError.Type, ShouldEqual, clrConst.TypTransferError)
			So(len(tranError.ErrorMessage), ShouldBeGreaterThan, 0)
		})
	})
}
