package proxy

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetInstance(t *testing.T) {
	tokenProxy := TokenProxy{}.GetInstance()
	Convey("Test GetInstance [ok] ---------->", t, func() {
		So(tokenProxy.Username, ShouldEqual, "staging.pangpang")
		So(tokenProxy.Password, ShouldEqual, "1111")
		So(tokenProxy.Token, ShouldEqual, "")
		So(tokenProxy.Expiration, ShouldEqual, 0)
	})
}

func TestGetToken(t *testing.T) {
	Convey("Test TestGetToken [ok] ---------->", t, func() {
		tokenProxy := TokenProxy{}.GetInstance()
		tokenProxy.Username = "staging.pangpang"
		tokenProxy.Password = "1111"
		header := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
		Convey(fmt.Sprintf("输入正确的用户名：%v，密码：%v，token的header部分应该为：%v", tokenProxy.Username, tokenProxy.Password, header), func() {
			token, err := tokenProxy.GetToken()
			So(err, ShouldBeNil)
			So(strings.Split(token, ".")[0], ShouldEqual, header)
		})
	})

	Convey("Test TestGetToken [error] ---------->", t, func() {
		tokenProxy := TokenProxy{}.GetInstance()
		tokenProxy.Username = "wrong username"
		tokenProxy.Password = "wrong username"
		tokenProxy.Expiration = 0
		Convey(fmt.Sprintf("输入错误的用户名：%v，密码：%v，应该返回error", tokenProxy.Username, tokenProxy.Password), func() {
			token, err := tokenProxy.GetToken()
			So(err, ShouldNotBeNil)
			So(token, ShouldEqual, "")
		})
		Convey("没有输入用户名或密码，应该返回error", func() {
			tokenProxy := TokenProxy{}.GetInstance()
			tokenProxy.Username = ""
			tokenProxy.Password = ""
			tokenProxy.Expiration = 0
			token, err := tokenProxy.GetToken()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "用户名或密码不能为空")
			So(token, ShouldEqual, "")
		})
	})
}
