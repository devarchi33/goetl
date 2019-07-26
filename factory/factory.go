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
	// clrEngine Clearance 数据库
	clrEngine *xorm.Engine
	// mslEngine msl2.0 数据库
	mslEngine *xorm.Engine
	once      sync.Once
)

// Init 初始化 数据库引擎
func Init() {
	cslEngine = CreateMSSQLEngine(config.GetCSLConnString())
	SetCSLEngine(cslEngine)

	invEngine = CreateMySQLEngine(config.GetINVConnString())
	SetINVEngine(invEngine)

	clrEngine = CreateMySQLEngine(config.GetClrConnString())
	SetClrEngine(clrEngine)

	mslEngine = CreateMySQLEngine(config.GetMSLConnString())
	SetMslEngine(mslEngine)
}

// GetCSLEngine 获取CSL数据库引擎
func GetCSLEngine() *xorm.Engine {
	return cslEngine
}

// SetCSLEngine 设置CSL数据库引擎
func SetCSLEngine(engine *xorm.Engine) {
	once.Do(func() {
		cslEngine = engine
	})
}

// GetINVEngine 获取MSL v1.0数据库引擎
func GetINVEngine() *xorm.Engine {
	return invEngine
}

// SetINVEngine 设置INV数据库引擎
func SetINVEngine(engine *xorm.Engine) {
	once.Do(func() {
		invEngine = engine
	})
}

// GetClrEngine 获取Clearance数据库引擎
func GetClrEngine() *xorm.Engine {
	return clrEngine
}

// SetClrEngine 设置Clearance数据库引擎
func SetClrEngine(engine *xorm.Engine) {
	once.Do(func() {
		clrEngine = engine
	})
}

//GetMslEngine 获取MSL数据库引擎
func GetMslEngine() *xorm.Engine {
	return mslEngine
}

// SetMslEngine 设置MSL数据库引擎
func SetMslEngine(engine *xorm.Engine) {
	once.Do(func() {
		mslEngine = engine
	})
}

// CreateMSSQLEngine 创建SQLServer数据库引擎
func CreateMSSQLEngine(connString string) *xorm.Engine {
	engine, err := xorm.NewEngine("mssql", connString)
	if err != nil {
		fmt.Println("createMSSQLEngine error")
	}
	engine.TZLocation, _ = time.LoadLocation("UTC")
	engine.SetTableMapper(core.SameMapper{})
	engine.SetColumnMapper(core.SameMapper{})

	return engine
}

// CreateMySQLEngine 创建MySQL数据库引擎
func CreateMySQLEngine(connString string) *xorm.Engine {
	var err error
	engine, err := xorm.NewEngine("mysql", connString)
	if err != nil {
		fmt.Println("createMySQLEngine error")
	}
	engine.SetTableMapper(core.SnakeMapper{})
	engine.SetColumnMapper(core.SnakeMapper{})

	return engine
}
