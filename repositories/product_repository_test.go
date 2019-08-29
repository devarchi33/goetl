package repositories

import (
	_ "clearance-adapter/test"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetSkueByCode(t *testing.T) {
	Convey("测试GetSkuByCode", t, func() {
		code := "SPYC949S1139085"
		id := 1
		productID := 1
		title := fmt.Sprintf("sku代码为%v的sku，id应该为%v，producID应该为%v", code, id, productID)
		Convey(title, func() {
			sku, has, err := ProductRepository{}.GetSkuByCode(code)
			So(err, ShouldBeNil)
			So(has, ShouldEqual, true)
			So(sku.ID, ShouldEqual, id)
			So(sku.ProductID, ShouldEqual, productID)
		})
	})
}
