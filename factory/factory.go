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
	// clrEngine Clearance 数据库
	clrEngine *xorm.Engine
	// p2brandEngine pangpang-brand 数据库
	p2brandEngine *xorm.Engine
	once          sync.Once
)

// Init 初始化 数据库引擎
func Init() {
	cslEngine = CreateMSSQLEngine(config.GetCSLConnString())
	SetCSLEngine(cslEngine)

	clrEngine = CreateMySQLEngine(config.GetClrConnString())
	SetClrEngine(clrEngine)

	p2brandEngine = CreateMySQLEngine(config.GetP2BrandConnString())
	SetP2BrandEngine(p2brandEngine)
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

//GetP2BrandEngine 获取pangpang-brand数据库引擎
func GetP2BrandEngine() *xorm.Engine {
	return p2brandEngine
}

// SetP2BrandEngine 设置MSL数据库引擎
func SetP2BrandEngine(engine *xorm.Engine) {
	once.Do(func() {
		p2brandEngine = engine
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
