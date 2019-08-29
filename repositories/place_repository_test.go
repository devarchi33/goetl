package repositories

import (
	_ "clearance-adapter/test"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetStoreByCode(t *testing.T) {
	Convey("测试GetStoreByCode", t, func() {
		code := "CEGP"
		id := 2
		title := fmt.Sprintf("卖场%v的id应该为%v", code, id)
		Convey(title, func() {
			store, has, err := PlaceRepository{}.GetStoreByCode(code)
			So(err, ShouldBeNil)
			So(has, ShouldEqual, true)
			So(store.ID, ShouldEqual, id)
		})
	})
}
