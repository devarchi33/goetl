package services

import (
	clrConst "clearance-adapter/domain/clr-constants"
	cslConst "clearance-adapter/domain/csl-constants"
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"clearance-adapter/repositories"
	_ "clearance-adapter/test"
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func rtwDefectiveSyncedShouldBeTrueOrFalse(shipmentLocationCode, waybillNo string, shouldBeTrue bool) {
	utc, _ := time.LoadLocation("")
	startDate := time.Now().In(utc).Format("2006-01-02T15:04:05Z")
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
	endDate := time.Now().In(utc).Format("2006-01-02T15:04:05Z")
	So(err, ShouldBeNil)
	So(len(result), ShouldEqual, 1)
	distList := infra.ConvertByteResult(result)
	if shouldBeTrue {
		So(distList[0]["synced"], ShouldEqual, "1")
	} else {
		So(distList[0]["synced"], ShouldEqual, "0")
	}
	So(distList[0]["last_updated_at"], ShouldBeGreaterThanOrEqualTo, startDate)
	So(distList[0]["last_updated_at"], ShouldBeLessThanOrEqualTo, endDate)
}

func TestReturnToWarehouseDefectiveETL(t *testing.T) {
	Convey("测试ReturnToWarehouseDefectiveETL的Run方法", t, func() {
		etl := ReturnToWarehouseDefectiveETL{}.New()
		err := etl.Run(context.Background())
		So(err, ShouldBeNil)
		waybillNo := "D20190819001"
		Convey("运单号为"+waybillNo+"的出库单synced=true，所以应该不会同步到CSL入库", func() {
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", "CEGP", waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 0)
		})
		waybillNo = "D20190813001"
		Convey("运单号为"+waybillNo+"的退仓出库单应该在CSL有记录", func() {
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("SA", "CEGP", waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 2)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, false)
				So(recvSupp.ShippingTypeCode, ShouldEqual, cslConst.TypRTWDefective)
				So(recvSupp.RecvSuppType, ShouldEqual, cslConst.TypSentOut)
				So(recvSupp.RecvSuppStatusCode, ShouldEqual, cslConst.StsSentOut)
				So(recvSupp.SuppEmpID, ShouldEqual, "7000028260")
				So(recvSupp.SuppEmpName, ShouldEqual, "史妍珣")
				So(recvSupp.RecvSuppMst.InUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppDtl.InUserID, ShouldEqual, "shi.yanxun")
				if recvSupp.ProdCode == "SPWJ948S2255070" {
					So(recvSupp.RecvSuppQty, ShouldEqual, 4)
				}
				if recvSupp.ProdCode == "SPYC949S1139095" {
					So(recvSupp.RecvSuppQty, ShouldEqual, 1)
				}
			}
			rtwDefectiveSyncedShouldBeTrueOrFalse("CEGP", waybillNo, true)
		})

		waybillNo = "D20190814001"
		Convey("CEGP的子卖场CJC1运单号为"+waybillNo+"的退仓出库单应该在CSL有记录", func() {
			recvSuppList, err := repositories.RecvSuppRepository{}.GetByWaybillNo("Q3", "CJC1", waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSuppList), ShouldEqual, 2)
			for _, recvSupp := range recvSuppList {
				So(recvSupp.RecvChk, ShouldEqual, false)
				So(recvSupp.ShippingTypeCode, ShouldEqual, cslConst.TypRTWDefective)
				So(recvSupp.RecvSuppType, ShouldEqual, cslConst.TypSentOut)
				So(recvSupp.RecvSuppStatusCode, ShouldEqual, cslConst.StsSentOut)
				So(recvSupp.SuppEmpID, ShouldEqual, "7000028260")
				So(recvSupp.SuppEmpName, ShouldEqual, "史妍珣")
				So(recvSupp.RecvSuppMst.InUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppDtl.InUserID, ShouldEqual, "shi.yanxun")
				So(recvSupp.RecvSuppMst.Dates, ShouldEqual, "20190814")
				So(recvSupp.ShopSuppRecvDate, ShouldEqual, "20190814")
				So(recvSupp.InvtBaseDate, ShouldEqual, "20190814")
				if recvSupp.ProdCode == "Q3AFAFDU6S2100230" {
					So(recvSupp.RecvSuppQty, ShouldEqual, 2)
				}
				if recvSupp.ProdCode == "Q3AFAFDU6S2100240" {
					So(recvSupp.RecvSuppQty, ShouldEqual, 3)
				}
			}
			rtwDefectiveSyncedShouldBeTrueOrFalse("CEGP", waybillNo, true)
		})

		brandCode := "SA"
		shipLocCode := "CFGY"
		waybillNo = "D20190910001"
		title := fmt.Sprintf("%v品牌，%v卖场，运单号为：%v的运单应该同步失败并且在error表中有记录", brandCode, shipLocCode, waybillNo)
		Convey(title, func() {
			has, distError, err := repositories.ReturnToWarehouseErrorRepository{}.GetByWaybillNo("XX", shipLocCode, waybillNo)
			So(has, ShouldEqual, true)
			So(err, ShouldBeNil)
			So(distError.Type, ShouldEqual, clrConst.TypReturnToWarehouseError)
			So(len(distError.ErrorMessage), ShouldBeGreaterThan, 0)
		})
	})
}
