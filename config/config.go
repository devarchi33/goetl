package config

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var congfigEnv = flag.String("env", os.Getenv("env"), "")

func readConfig(env string) {
	viper.SetConfigName("config." + env)
	viper.AddConfigPath(".")
	viper.AddConfigPath(getCurrPath())

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

func getCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	filepath.Abs(filepath.Dir(os.Args[0]))
	return ret
}

// GetCSLConnString CSL 数据库连接字符串
func GetCSLConnString() string {
	v := os.Getenv("CslConnString")
	if v != "" {
		return v
	}

	if *congfigEnv == "" {
		defaultAppEnv := "mslv2-dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)
	return viper.Get("cslConnString").(string)
}

// GetClrConnString Clearance 数据库连接字符串
func GetClrConnString() string {
	v := os.Getenv("ClrConnString")
	if v != "" {
		return v
	}

	if *congfigEnv == "" {
		defaultAppEnv := "mslv2-dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)
	return viper.Get("clrConnString").(string)
}

// GetP2BrandConnString pangpang-brand 数据库连接字符串
func GetP2BrandConnString() string {
	v := os.Getenv("P2brandConnString")
	if v != "" {
		return v
	}

	if *congfigEnv == "" {
		defaultAppEnv := "mslv2-dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)
	return viper.Get("p2brandConnString").(string)
}

// GetP2BrandSkuLocationAPIRoot pangpang-brand sku-location API Root
func GetP2BrandSkuLocationAPIRoot() string {
	v := os.Getenv("P2brandSkuLocationAPIRoot")
	if v != "" {
		return v
	}

	if *congfigEnv == "" {
		defaultAppEnv := "mslv2-dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)
	return viper.Get("p2brandSkuLocationAPIRoot").(string)
}

// GetP2BrandColleagueAPIRoot pangpang-brand colleague API Root
func GetP2BrandColleagueAPIRoot() string {
	v := os.Getenv("P2brandColleagueAPIRoot")
	if v != "" {
		return v
	}

	if *congfigEnv == "" {
		defaultAppEnv := "mslv2-dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)
	return viper.Get("p2brandColleagueAPIRoot").(string)
}

func GetBatchJobAPIRoot() string {
	v := os.Getenv("BatchjobAPIRoot")
	if v != "" {
		return v
	}
	if *congfigEnv == "" {
		defaultAppEnv := "mslv2-dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)
	return viper.Get("batchjobAPIRoot").(string)
}

// GetAutoDistributeDeadlineDays 自动入库截止日期天数
func GetAutoDistributeDeadlineDays() int {
	v, err := strconv.Atoi(os.Getenv("AutoDistributeDeadlineDays"))
	if err == nil {
		return v
	}

	if *congfigEnv == "" {
		defaultAppEnv := "mslv2-dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)
	return viper.Get("autoDistributeDeadlineDays").(int)
}

// GetAutoTransferInDeadlineDays 自动入库截止日期天数
func GetAutoTransferInDeadlineDays() int {
	v, err := strconv.Atoi(os.Getenv("AutoTransferInDeadlineDays"))
	if err == nil {
		return v
	}

	if *congfigEnv == "" {
		defaultAppEnv := "mslv2-dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)
	return viper.Get("autoTransferInDeadlineDays").(int)
}

// GetColleagueClearanceUserInfo 获取Clearance使用的用户名的信息，租户代码，用户名，密码
func GetColleagueClearanceUserInfo() (string, string, string) {
	tenantCode := os.Getenv("ColleagueClearanceTenantCode")
	username := os.Getenv("ColleagueClearanceUsername")
	password := os.Getenv("ColleagueClearancePassword")
	if tenantCode != "" && username != "" && password != "" {
		return tenantCode, username, password
	}

	if *congfigEnv == "" {
		defaultAppEnv := "mslv2-dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)

	return viper.Get("colleagueClearanceTenantCode").(string),
		viper.Get("colleagueClearanceUsername").(string),
		viper.Get("colleagueClearancePassword").(string)
}

//GetTenantCode ...
func GetTenantCode() string {
	v := os.Getenv("ColleagueClearanceTenantCode")
	if len(v) > 0 {
		return v
	}
	if *congfigEnv == "" {
		defaultAppEnv := "mslv2-dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)
	return viper.Get("colleagueClearanceTenantCode").(string)
}
