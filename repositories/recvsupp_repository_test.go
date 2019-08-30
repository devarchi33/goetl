package repositories

import (
	cslConst "clearance-adapter/domain/csl-constants"
	"fmt"
	"strings"
	"testing"

	"clearance-adapter/factory"
	"clearance-adapter/models"
	"clearance-adapter/test"
	_ "clearance-adapter/test"

	"clearance-adapter/repositories/entities"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInit(t *testing.T) {
	Convey("测试初始化数据库", t, func() {
		So(1, ShouldEqual, 1)
	})
}

func TestGetUnconfirmedDistributionOrdersByDeadline(t *testing.T) {
	deadline := "20190801"
	title := fmt.Sprintf("%v日以前的，为入库的出库单应该有4个", deadline)
	Convey(title, t, func() {
		result, err := RecvSuppRepository{}.GetUnconfirmedDistributionOrdersByDeadline(deadline)
		So(err, ShouldBeNil)
		orders := make(map[string]string)
		for _, v := range result {
			orders[v.RecvSuppMst.WayBillNo] = v.RecvSuppMst.WayBillNo
		}
		So(len(orders), ShouldEqual, 4)
	})
}

func TestGetByWaybillNo(t *testing.T) {
	Convey("测试GetByWaybillNo", t, func() {
		Convey("SA品牌的CEGP卖场的1010590009008运单, 应该有14个商品", func() {
			result, err := RecvSuppRepository{}.GetByWaybillNo("SA", "CEGP", "1010590009008")
			So(err, ShouldBeNil)
			So(len(result), ShouldEqual, 14)
		})
	})
}

func TestCreateReturnToWarehouseOrder(t *testing.T) {
	Convey("测试CreateReturnToWarehouseOrder", t, func() {
		recvSuppNo, err := RecvSuppRepository{}.CreateReturnToWarehouseOrder("SA", "CEGP", "2010590009008", "20190823", "7000028260", "456")
		Convey("SA-CEGP卖场应该有一个运单号为2010590009008的退仓订单", func() {
			So(err, ShouldBeNil)
			So(strings.HasPrefix(recvSuppNo, "CEGP"), ShouldEqual, true)
		})
		Convey("Waybill表中应该存在运单号为2010590009008的数据", func() {
			sql := `
				SELECT * FROM WayBillNo WHERE ShippingCompanyCode= 'SR' AND WayBillNo = '2010590009008'
			`
			result, err := factory.GetCSLEngine().Query(sql)
			So(err, ShouldBeNil)
			So(len(result), ShouldEqual, 1)
		})
	})
}

func TestAddReturnToWarehouseOrderItem(t *testing.T) {
	test.InitTestData()

	brandCode := "SA"
	shopCode := "CEGP"
	waybillNo := "2010590009008"
	outDate := "20190823"
	empID := "7000028260"
	skuCode := "SPWJ948S2255070"
	Convey("测试AddReturnToWarehouseOrderItem", t, func() {
		recvSuppNo, err := RecvSuppRepository{}.CreateReturnToWarehouseOrder(brandCode, shopCode, waybillNo, outDate, empID, "456")
		So(err, ShouldBeNil)
		err = RecvSuppRepository{}.AddReturnToWarehouseOrderItem(brandCode, shopCode, outDate, recvSuppNo, skuCode, 1, empID)
		So(err, ShouldBeNil)
		Convey("SA-CEGP卖场运单号为2010590009008的退仓订单中应该有一个商品代码为SPWJ948S2255070的商品，并且出库数量为1", func() {
			recvSupp, err := RecvSuppRepository{}.GetByWaybillNo(brandCode, shopCode, waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSupp), ShouldEqual, 1)
			So(recvSupp[0].RecvSuppMst.Dates, ShouldEqual, outDate)
			So(recvSupp[0].ProdCode, ShouldEqual, skuCode)
			So(recvSupp[0].RecvSuppQty, ShouldEqual, 1)
		})
	})
}

var TransferOrderRecvSuppNp string

func TestCreateTransferOrder(t *testing.T) {
	Convey("测试CreateTransferOrder", t, func() {
		brandCode := "SA"
		shipmentLocationCode := "CEGP"
		receiptLocationCode := "CJ2F"
		waybillNo := "20190815001"
		boxNo := "20190815001-1"
		shippingCompanyCode := "SR"
		deliveryOrderNo := "456"
		empID := "7000028260"
		outDate := "20190826"

		Convey("SA-CEGP卖场应该有一个运单号为20190815001的调货单，调给SA-CJ2F，并且不是跨广域支社调货", func() {
			recvSuppNo, err := RecvSuppRepository{}.CreateTransferOrder(brandCode, shipmentLocationCode, receiptLocationCode, outDate, waybillNo, boxNo, shippingCompanyCode, deliveryOrderNo, empID)
			So(err, ShouldBeNil)
			So(strings.HasPrefix(recvSuppNo, shipmentLocationCode), ShouldEqual, true)
			masters := make([]models.RecvSuppMst, 0)
			err = factory.GetCSLEngine().Where("BrandCode = ? AND ShopCode = ? AND WayBillNo = ?",
				brandCode, shipmentLocationCode, waybillNo).
				Find(&masters)
			So(err, ShouldBeNil)
			So(len(masters), ShouldEqual, 1)
			So(masters[0].TargetShopCode, ShouldEqual, receiptLocationCode)
			So(masters[0].RecvSuppStatusCode, ShouldEqual, cslConst.StsSentOut)
			TransferOrderRecvSuppNp = recvSuppNo
		})
		Convey("Waybill表中应该存在运单号为20190815001的数据", func() {
			sql := `
				SELECT * FROM WayBillNo WHERE ShippingCompanyCode= ? AND WayBillNo = ?
			`
			result, err := factory.GetCSLEngine().Query(sql, shippingCompanyCode, waybillNo)
			So(err, ShouldBeNil)
			So(len(result), ShouldEqual, 1)
		})

		Convey("SA-CEGP卖场应该有一个运单号为20190815002的调货单，调给SA-CFGY，并且是跨广域支社调货", func() {
			waybillNo := "20190815002"
			boxNo := "20190815002-1"
			receiptLocationCode = "CFGY"
			outDate := "20190826"
			recvSuppNo, err := RecvSuppRepository{}.CreateTransferOrder(brandCode, shipmentLocationCode, receiptLocationCode, outDate, waybillNo, boxNo, shippingCompanyCode, deliveryOrderNo, empID)
			So(err, ShouldBeNil)
			So(strings.HasPrefix(recvSuppNo, shipmentLocationCode), ShouldEqual, true)
			masters := make([]models.RecvSuppMst, 0)
			err = factory.GetCSLEngine().Where("BrandCode = ? AND ShopCode = ? AND WayBillNo = ?",
				brandCode, shipmentLocationCode, waybillNo).
				Find(&masters)
			So(err, ShouldBeNil)
			So(len(masters), ShouldEqual, 1)
			So(masters[0].TargetShopCode, ShouldEqual, receiptLocationCode)
			So(masters[0].RecvSuppStatusCode, ShouldEqual, "W")
		})
	})
}

func TestAddTransferOrderItem(t *testing.T) {
	Convey("测试AddTransferOrderItem", t, func() {
		brandCode := "SA"
		shipmentLocationCode := "CEGP"
		recvSuppNo := TransferOrderRecvSuppNp
		empID := "7000028260"
		waybillNo := "20190815001"
		outDate := "20190826"
		Convey("SA-CEGP卖场运单号为20190815001的调货单，添加两个商品，SPWJ948S2255070的数量应该为1，SPYS949H2250100的商品数量应该为2", func() {
			skuCode := "SPWJ948S2255070"
			qty := 1
			err := RecvSuppRepository{}.AddTransferOrderItem(brandCode, shipmentLocationCode, outDate, recvSuppNo, skuCode, qty, empID)
			So(err, ShouldBeNil)
			skuCode = "SPYS949H2250100"
			qty = 2
			err = RecvSuppRepository{}.AddTransferOrderItem(brandCode, shipmentLocationCode, outDate, recvSuppNo, skuCode, qty, empID)
			So(err, ShouldBeNil)
			recvSupp, err := RecvSuppRepository{}.GetByWaybillNo(brandCode, shipmentLocationCode, waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSupp), ShouldEqual, 2)
			for _, v := range recvSupp {
				if v.ProdCode == "SPWJ948S2255070" {
					So(v.RecvSuppQty, ShouldEqual, 1)
				}
				if v.ProdCode == "SPYS949H2250100" {
					So(v.RecvSuppQty, ShouldEqual, 2)
				}
			}
		})
	})
}

func TestConfirmTransferOrder(t *testing.T) {
	Convey("测试ConfirmTransferOrder", t, func() {
		brandCode := "SA"
		receiptLocationCode := "CJ2F"
		shipmentLocationCode := "CEGP"
		roundRecvSuppNo := TransferOrderRecvSuppNp
		empID := "7000028260"
		waybillNo := "20190815001"
		boxNo := "20190815001-1"
		inDate := "20190826"
		Convey("SA-CJ2F卖场应该有对运单号为20190815001的调货单有确认记录，并且SPWJ948S2255070的数量应该为1，SPYS949H2250100的商品数量应该为2", func() {
			recvSuppNo, err := RecvSuppRepository{}.ConfirmTransferOrder(brandCode, receiptLocationCode, shipmentLocationCode, inDate, waybillNo, boxNo, roundRecvSuppNo, empID)
			So(err, ShouldBeNil)
			So(strings.HasPrefix(recvSuppNo, receiptLocationCode), ShouldEqual, true)
			masters := make([]models.RecvSuppMst, 0)
			err = factory.GetCSLEngine().Where("RecvSuppNo = ?", recvSuppNo).Find(&masters)
			So(err, ShouldBeNil)
			So(len(masters), ShouldEqual, 1)
			So(masters[0].ShopCode, ShouldEqual, receiptLocationCode)
			So(masters[0].TargetShopCode, ShouldEqual, shipmentLocationCode)
			So(masters[0].WayBillNo, ShouldEqual, waybillNo)
			So(masters[0].RecvSuppStatusCode, ShouldEqual, cslConst.StsConfirmed)
			details := make([]models.RecvSuppDtl, 0)
			err = factory.GetCSLEngine().Where("RecvSuppNo = ?", recvSuppNo).Find(&details)
			So(err, ShouldBeNil)
			for _, v := range details {
				if v.ProdCode == "SPWJ948S2255070" {
					So(v.RecvSuppQty, ShouldEqual, 1)
				}
				if v.ProdCode == "SPYS949H2250100" {
					So(v.RecvSuppQty, ShouldEqual, 2)
				}
			}
		})
	})
}

func TestCreateTransferOrderSet(t *testing.T) {
	test.InitTestData()
	Convey("测试CreateTransferOrderSet", t, func() {
		brandCode := "SA"
		shipmentLocationCode := "CEGP"
		receiptLocationCode := "CJ2F"
		empID := "7000028260"
		waybillNo := "20190822001"
		boxNo := "20190822001-1"
		outDate := "20190822"
		inDate := "20190823"
		shippingCompanyCode := "SR"
		orderSet := entities.TransferOrderSet{
			BrandCode:            brandCode,
			ShipmentLocationCode: shipmentLocationCode,
			ShipmentShopCode:     shipmentLocationCode,
			ReceiptLocationCode:  receiptLocationCode,
			ReceiptShopCode:      receiptLocationCode,
			WaybillNo:            waybillNo,
			BoxNo:                boxNo,
			ShippingCompanyCode:  shippingCompanyCode,
			DeliveryOrderNo:      "",
			OutDate:              outDate,
			InDate:               inDate,
			OutEmpID:             empID,
			InEmpID:              empID,
			Items: []entities.TransferOrderSetItem{
				entities.TransferOrderSetItem{
					SkuCode: "SPWJ948S2255070",
					Qty:     1,
				},
				entities.TransferOrderSetItem{
					SkuCode: "SPYS949H2250100",
					Qty:     2,
				},
			},
		}
		Convey("SA-CEGP卖场应该有一个运单号为20190822001的调货单，调给SA-CJ2F，并且不是跨广域支社调货", func() {
			err := RecvSuppRepository{}.CreateTransferOrderSet(orderSet)
			So(err, ShouldBeNil)
			masters := make([]models.RecvSuppMst, 0)
			err = factory.GetCSLEngine().Where("BrandCode = ? AND ShopCode = ? AND WayBillNo = ?",
				brandCode, shipmentLocationCode, waybillNo).
				Find(&masters)
			So(err, ShouldBeNil)
			So(len(masters), ShouldEqual, 1)
			So(masters[0].Dates, ShouldEqual, outDate)
			So(masters[0].ShopSuppRecvDate, ShouldEqual, outDate)
			So(masters[0].TargetShopCode, ShouldEqual, receiptLocationCode)
			So(masters[0].RecvSuppStatusCode, ShouldEqual, cslConst.StsConfirmed)
		})
		Convey("Waybill表中应该存在运单号为20190822001的数据", func() {
			sql := `
				SELECT * FROM WayBillNo WHERE ShippingCompanyCode= ? AND WayBillNo = ?
			`
			result, err := factory.GetCSLEngine().Query(sql, shippingCompanyCode, waybillNo)
			So(err, ShouldBeNil)
			So(len(result), ShouldEqual, 1)
		})
		Convey("SA-CJ2F卖场应该有对运单号为20190822001的调货单有确认记录，并且SPWJ948S2255070的数量应该为1，SPYS949H2250100的商品数量应该为2", func() {
			masters := make([]models.RecvSuppMst, 0)
			err := factory.GetCSLEngine().Where("BrandCode = ? AND ShopCode = ? AND WayBillNo = ?",
				brandCode, receiptLocationCode, waybillNo).
				Find(&masters)
			So(err, ShouldBeNil)
			So(len(masters), ShouldEqual, 1)
			So(masters[0].ShopCode, ShouldEqual, receiptLocationCode)
			So(masters[0].TargetShopCode, ShouldEqual, shipmentLocationCode)
			So(masters[0].WayBillNo, ShouldEqual, waybillNo)
			So(masters[0].RecvSuppStatusCode, ShouldEqual, cslConst.StsConfirmed)
			So(masters[0].Dates, ShouldEqual, inDate)
			So(masters[0].ShopSuppRecvDate, ShouldEqual, inDate)
			details := make([]models.RecvSuppDtl, 0)
			err = factory.GetCSLEngine().Where("RecvSuppNo = ?", masters[0].RecvSuppNo).Find(&details)
			So(err, ShouldBeNil)
			for _, v := range details {
				if v.ProdCode == "SPWJ948S2255070" {
					So(v.RecvSuppQty, ShouldEqual, 1)
				}
				if v.ProdCode == "SPYS949H2250100" {
					So(v.RecvSuppQty, ShouldEqual, 2)
				}
			}
		})
	})
}
