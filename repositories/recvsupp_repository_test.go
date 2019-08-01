package repositories

import (
	"testing"

	_ "clearance-adapter/test"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetByWaybillNo(t *testing.T) {
	Convey("测试GetByWaybillNo", t, func() {
		Convey("SA品牌的CEGP卖场的1010590009008运单, 应该有14个商品", func() {
			result, err := RecvSuppRepository{}.GetByWaybillNo("SA", "CEGP", "1010590009008")
			So(err, ShouldBeNil)
			So(len(result), ShouldEqual, 14)
		})
	})
}

func TestInit(t *testing.T) {
	Convey("测试初始化数据库", t, func() {
		So(1, ShouldEqual, 1)
	})
}
