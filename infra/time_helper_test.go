package infra

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParse8BitsDate(t *testing.T) {

	intput := "2019-08-23T13:48:20Z"
	output := "20190823"
	Convey(fmt.Sprintf("输入%s，应该输出%s", intput, output), t, func() {
		date := Parse8BitsDate("2019-08-23T13:48:20Z", nil)
		So(date, ShouldEqual, "20190823")
	})

	intput = "2019-08-23T03:48:20Z"
	output = "20190823"
	Convey(fmt.Sprintf("输入%s，应该输出%s", intput, output), t, func() {
		date := Parse8BitsDate("2019-08-23T03:48:20Z", nil)
		So(date, ShouldEqual, "20190823")
	})

	intput = "2019-08-23T23:48:20Z"
	output = "20190824"
	Convey(fmt.Sprintf("输入%s，应该输出%s", intput, output), t, func() {
		date := Parse8BitsDate("2019-08-23T23:48:20Z", nil)
		So(date, ShouldEqual, "20190824")
	})
}
