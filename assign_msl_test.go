package main

import (
	_ "clearance-adapter/test"
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pangpanglabs/goetl"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAssignMsl(t *testing.T) {
	Convey("测试AssignETL的Run方法", t, func() {
		Convey("可以把入库预约从MSL导入到clearance", func() {
			etl := goetl.New(AssignMslETL{})
			err := etl.Run(context.Background())
			So(err, ShouldBeNil)
		})
	})
}
