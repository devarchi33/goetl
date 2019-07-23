package factory

import (
	"clearance-adapter/config"
	"fmt"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var (
	// cslEngine CSL 数据库
	cslEngine *xorm.Engine
	// invEngine MSL v1.0 数据库
	invEngine *xorm.Engine
	once      sync.Once
)

// Init 初始化 数据库引擎
func Init() {
	cslEngine = CreateMSSQLEngine(config.GetCSLConnString())
	SetCSLEngine(cslEngine)

	invEngine = CreateMySQLEngine(config.GetINVConnString())
	SetINVEngine(invEngine)
}

// SetCSLEngine 设置CSL数据库引擎
func SetCSLEngine(engine *xorm.Engine) {
	once.Do(func() {
		cslEngine = engine
	})
}

// SetINVEngine 设置INV数据库引擎
func SetINVEngine(engine *xorm.Engine) {
	once.Do(func() {
		invEngine = engine
	})
}

// CreateMSSQLEngine 创建SQLServer数据库引擎
func CreateMSSQLEngine(connString string) *xorm.Engine {
	engine, err := xorm.NewEngine("mssql", connString)
	if err != nil {
		fmt.Println("createCSLEngine error")
	}
	engine.TZLocation, _ = time.LoadLocation("UTC")
	engine.SetTableMapper(core.SameMapper{})
	engine.SetColumnMapper(core.SameMapper{})

	return engine
}

// GetCSLEngine 获取CSL数据库引擎
func GetCSLEngine() *xorm.Engine {
	return cslEngine
}

// CreateMySQLEngine 创建MySQL数据库引擎
func CreateMySQLEngine(connString string) *xorm.Engine {
	var err error
	engine, err := xorm.NewEngine("mysql", connString)
	if err != nil {
		fmt.Println("createINVEngine error")
	}
	engine.SetTableMapper(core.SnakeMapper{})
	engine.SetColumnMapper(core.SnakeMapper{})

	return engine
}

// GetINVEngine 获取MSL v1.0数据库引擎
func GetINVEngine() *xorm.Engine {
	return invEngine
}
