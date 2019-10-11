package config

import (
	"fmt"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"strconv"
)

func TestGetCSLConnString(t *testing.T) {
	Convey("测试CSLConnString的值", t, func(){
		cslConnString := "test csl connection string"
		Convey(fmt.Sprintf("没有设置环境变量的时候，应该为%v", cslConnString), func() {
			v := GetCSLConnString()
			So(v, ShouldEqual, cslConnString)
		})
	
		envCslConnString := "env test csl connection string"
		os.Setenv("CslConnString", envCslConnString)
		Convey(fmt.Sprintf("设置环境变量的后，应该为%v", envCslConnString), func() {
			v := GetCSLConnString()
			So(v, ShouldEqual, envCslConnString)
		})
	})
}

func TestGetClrConnString(t *testing.T) {
	Convey("测试ClrConnString的值", t, func(){
		clrConnString := "test clearance connection string"
		Convey(fmt.Sprintf("没有设置环境变量的时候，应该为%v", clrConnString), func() {
			v := GetClrConnString()
			So(v, ShouldEqual, clrConnString)
		})

		envClrConnString := "env test clearance connection string"
		os.Setenv("ClrConnString", envClrConnString)
		Convey(fmt.Sprintf("设置环境变量的后，应该为%v", envClrConnString), func() {
			v := GetClrConnString()
			So(v, ShouldEqual, envClrConnString)
		})
	})
}

func TestGetP2BrandConnString(t *testing.T) {
	Convey("测试P2BrandConnString的值", t, func(){
		p2brandConnString := "test pangpang brand connection string"
		Convey(fmt.Sprintf("没有设置环境变量的时候，应该为%v", p2brandConnString), func() {
			v := GetP2BrandConnString()
			So(v, ShouldEqual, p2brandConnString)
		})

		envP2brandConnString := "env test pangpang brand connection string"
		os.Setenv("P2brandConnString", envP2brandConnString)
		Convey(fmt.Sprintf("设置环境变量的后，应该为%v", envP2brandConnString), func() {
			v := GetP2BrandConnString()
			So(v, ShouldEqual, envP2brandConnString)
		})
	})
}

func TestGetP2BrandSkuLocationAPIRoot(t *testing.T) {
	Convey("测试P2brandSkuLocationAPIRoot的值", t, func(){
		p2brandSkuLocationAPIRoot := "test pangpang brand sku location api root"
		Convey(fmt.Sprintf("没有设置环境变量的时候，应该为%v", p2brandSkuLocationAPIRoot), func() {
			v := GetP2BrandSkuLocationAPIRoot()
			So(v, ShouldEqual, p2brandSkuLocationAPIRoot)
		})

		envP2brandSkuLocationAPIRoot := "env test pangpang brand sku location api root"
		os.Setenv("P2brandSkuLocationAPIRoot", envP2brandSkuLocationAPIRoot)
		Convey(fmt.Sprintf("设置环境变量的后，应该为%v", envP2brandSkuLocationAPIRoot), func() {
			v := GetP2BrandSkuLocationAPIRoot()
			So(v, ShouldEqual, envP2brandSkuLocationAPIRoot)
		})
	})
}

func TestGetP2BrandColleagueAPIRoot(t *testing.T) {
	Convey("测试P2brandColleagueAPIRoot的值", t, func(){
		p2brandColleagueAPIRoot := "test pangpang brand colleague root"
		Convey(fmt.Sprintf("没有设置环境变量的时候，应该为%v", p2brandColleagueAPIRoot), func() {
			v := GetP2BrandColleagueAPIRoot()
			So(v, ShouldEqual, p2brandColleagueAPIRoot)
		})

		envP2brandColleagueAPIRoot := "env test pangpang brand colleague root"
		os.Setenv("P2brandColleagueAPIRoot", envP2brandColleagueAPIRoot)
		Convey(fmt.Sprintf("设置环境变量的后，应该为%v", envP2brandColleagueAPIRoot), func() {
			v := GetP2BrandColleagueAPIRoot()
			So(v, ShouldEqual, envP2brandColleagueAPIRoot)
		})
	})
}

func TestGetAutoDistributeDeadlineDays(t *testing.T) {
	Convey("测试AutoDistributeDeadlineDays的值", t, func(){
		autoDistributeDeadlineDays := 7
		Convey(fmt.Sprintf("没有设置环境变量的时候，应该为%v", autoDistributeDeadlineDays), func() {
			v := GetAutoDistributeDeadlineDays()
			So(v, ShouldEqual, autoDistributeDeadlineDays)
		})

		envAutoDistributeDeadlineDays := 17
		os.Setenv("AutoDistributeDeadlineDays", strconv.Itoa(envAutoDistributeDeadlineDays))
		Convey(fmt.Sprintf("设置环境变量的后，应该为%v", envAutoDistributeDeadlineDays), func() {
			v := GetAutoDistributeDeadlineDays()
			So(v, ShouldEqual, envAutoDistributeDeadlineDays)
		})
	})
}

func TestGetAutoTransferInDeadlineDays(t *testing.T) {
	Convey("测试AutoTransferInDeadlineDays的值", t, func(){
		autoTransferInDeadlineDays := 6
		Convey(fmt.Sprintf("没有设置环境变量的时候，应该为%v", autoTransferInDeadlineDays), func() {
			v := GetAutoTransferInDeadlineDays()
			So(v, ShouldEqual, autoTransferInDeadlineDays)
		})

		envAutoTransferInDeadlineDays := 16
		os.Setenv("AutoTransferInDeadlineDays", strconv.Itoa(envAutoTransferInDeadlineDays))
		Convey(fmt.Sprintf("设置环境变量的后，应该为%v", envAutoTransferInDeadlineDays), func() {
			v := GetAutoTransferInDeadlineDays()
			So(v, ShouldEqual, envAutoTransferInDeadlineDays)
		})
	})
}

func TestGetColleagueClearanceUserInfo(t *testing.T) {
	Convey("测试ColleagueClearanceUserInfo的值", t, func(){
		tenantCode:= "test pangpang"
		username:= "test.pangpang"
		password:= "2222"
		Convey(fmt.Sprintf("没有设置环境变量的时候，tenant code为%v，username为%v，password为%v", tenantCode, username, password), func() {
			t, u, p := GetColleagueClearanceUserInfo()
			So(t, ShouldEqual, tenantCode)
			So(u, ShouldEqual, username)
			So(p, ShouldEqual, password)
		})

		envTenantCode:= "env test pangpang"
		envUsername:= "env test.pangpang"
		envPassword:= "env 2222"
		os.Setenv("ColleagueClearanceTenantCode", envTenantCode)
		os.Setenv("ColleagueClearanceUsername", envUsername)
		os.Setenv("ColleagueClearancePassword", envPassword)
		Convey(fmt.Sprintf("设置环境变量的后，tenant code为%v，username为%v，password为%v", envTenantCode, envUsername, envPassword), func() {
			t, u, p  := GetColleagueClearanceUserInfo()
			So(t, ShouldEqual, envTenantCode)
			So(u, ShouldEqual, envUsername)
			So(p, ShouldEqual, envPassword)
		})
	})
}