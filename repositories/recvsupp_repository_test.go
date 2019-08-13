package repositories

import (
	"strings"
	"testing"
	"time"

	"clearance-adapter/factory"
	"clearance-adapter/test"
	_ "clearance-adapter/test"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInit(t *testing.T) {
	Convey("测试初始化数据库", t, func() {
		So(1, ShouldEqual, 1)
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
		recvSuppNo, err := RecvSuppRepository{}.CreateReturnToWarehouseOrder("SA", "CEGP", "2010590009008", "7000028260", "456")
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
	empID := "7000028260"
	skuCode := "SPWJ948S2255070"
	Convey("测试AddReturnToWarehouseOrderItem", t, func() {
		recvSuppNo, err := RecvSuppRepository{}.CreateReturnToWarehouseOrder(brandCode, shopCode, waybillNo, empID, "456")
		So(err, ShouldBeNil)
		err = RecvSuppRepository{}.AddReturnToWarehouseOrderItem(brandCode, shopCode, recvSuppNo, skuCode, 1, empID)
		So(err, ShouldBeNil)
		Convey("SA-CEGP卖场运单号为2010590009008的退仓订单中应该有一个商品代码为SPWJ948S2255070的商品，并且出库数量为1", func() {
			recvSupp, err := RecvSuppRepository{}.GetByWaybillNo(brandCode, shopCode, waybillNo)
			So(err, ShouldBeNil)
			So(len(recvSupp), ShouldEqual, 1)
			So(recvSupp[0].RecvSuppMst.Dates, ShouldEqual, time.Now().Format("20060102"))
			So(recvSupp[0].ProdCode, ShouldEqual, skuCode)
			So(recvSupp[0].RecvSuppQty, ShouldEqual, 1)
		})
	})
}
