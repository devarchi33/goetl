package proxy

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTokenByUsernameAndPassword(t *testing.T) {
	Convey("Test GetTokenByUsernameAndPassword [ok] ---------->", t, func() {
		username := "staging.pangpang"
		password := "1111"
		header := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
		Convey(fmt.Sprintf("输入正确的用户名：%v，密码：%v，token的header部分应该为：%v", username, password, header), func() {
			token, err := ColleagueServiceProxy{}.GetTokenByUsernameAndPassword(username, password)
			So(err, ShouldBeNil)
			So(strings.Split(token, ".")[0], ShouldEqual, header)
		})
	})

	Convey("Test GetTokenByUsernameAndPassword [error] ---------->", t, func() {
		username := "wrong username"
		password := "wrong password"
		Convey(fmt.Sprintf("输入错误的用户名：%v，密码：%v，应该返回error", username, password), func() {
			token, err := ColleagueServiceProxy{}.GetTokenByUsernameAndPassword(username, password)
			So(err, ShouldNotBeNil)
			So(token, ShouldEqual, "")
		})
		Convey("没有输入用户名或密码，应该返回error", func() {
			token, err := ColleagueServiceProxy{}.GetTokenByUsernameAndPassword("", "")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "用户名或密码不能为空")
			So(token, ShouldEqual, "")
		})
	})
}
