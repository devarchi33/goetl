package config

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	if *congfigEnv == "" {
		defaultAppEnv := "dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)
	return viper.Get("cslConnString").(string)
}

// GetINVConnString MSL v1.0 数据库连接字符串
func GetINVConnString() string {
	if *congfigEnv == "" {
		defaultAppEnv := "dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)
	return viper.Get("invConnString").(string)
}

// GetClrConnString Clearance 数据库连接字符串
func GetClrConnString() string {
	if *congfigEnv == "" {
		defaultAppEnv := "dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)
	return viper.Get("clrConnString").(string)
}

// GetMSLConnString MSL v2.0 数据库连接字符串
func GetMSLConnString() string {
	if *congfigEnv == "" {
		defaultAppEnv := "dev"
		congfigEnv = &defaultAppEnv
	}
	readConfig(*congfigEnv)
	return viper.Get("mslConnString").(string)
}
