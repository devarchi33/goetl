package services

import (
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"clearance-adapter/repositories"
	_ "clearance-adapter/test"
	"context"
	"testing"
	"time"

	cslConst "clearance-adapter/domain/csl-constants"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func transferOrderSyncedShouldBeTrue(shipmentLocationCode, receiptLocationCode, waybillNo string) {
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
	So(err, ShouldBeNil)
	So(len(result), ShouldEqual, 1)
	distList := infra.ConvertByteResult(result)
	So(distList[0]["synced"], ShouldEqual, "1")
	now := time.Now()
	utc, _ := time.LoadLocation("")
	So(distList[0]["last_updated_at"], ShouldEqual, now.In(utc).Format("2006-01-02T15:04:05Z"))
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

		Convey("CEGP卖场调到CFGY卖场的运单为20190821001的调货单同步后，应该在CSL有调出数据", func() {
			recvSupp, err := repositories.RecvSuppRepository{}.GetByWaybillNo(brandCode, shipmentLocationCode, waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSupp), ShouldEqual, 2)
			for _, v := range recvSupp {
				So(v.RecvSuppStatusCode, ShouldEqual, cslConst.StsConfirmed)
				if v.ProdCode == "SPWJ948S2255070" {
					So(v.RecvSuppQty, ShouldEqual, 4)
				}
				if v.ProdCode == "SPYC949S1139095" {
					So(v.RecvSuppQty, ShouldEqual, 1)
				}
			}
		})
		Convey("CEGP卖场调到CFGY卖场的运单为20190821001的调货单同步后，应该在CSL有调入数据", func() {
			recvSupp, err := repositories.RecvSuppRepository{}.GetByWaybillNo(brandCode, receiptLocationCode, waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSupp), ShouldEqual, 2)
			for _, v := range recvSupp {
				So(v.RecvSuppStatusCode, ShouldEqual, cslConst.StsConfirmed)
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
	})
}
