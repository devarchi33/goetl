package models

import "time"

// StockDistribute ...
type StockDistribute struct {
	ID                 int64     `xorm:"bigint not null autoincr pk 'id'"`
	TenantCode         string    `xorm:"varchar(50)"`
	BrandCode          string    `xorm:"varchar(50)"`
	BoxNo              string    `xorm:"varchar(50) not null default ''"`
	WaybillNo          string    `xorm:"varchar(50) not null default ''"`
	ShipmentLocationID int64     `xorm:"bigint 'shipment_location_Id'"`
	ReceiptLocationID  int64     `xorm:"bigint 'receipt_location_id'"`
	CreatedAt          time.Time `xorm:"datetime"`
	ColleagueID        int64     `xorm:"bigint 'colleague_id'"`
}

// StockDistributeItem ...
type StockDistributeItem struct {
	ID                int64     `xorm:"bigint not null autoincr pk 'id'"`
	StockDistributeID int64     `xorm:"bigint not null 'stock_distribute_id'"`
	ProductID         int64     `xorm:"bigint not null 'product_id'"`
	SkuID             int64     `xorm:"bigint not null 'sku_id'"`
	Barcode           string    `xorm:"varchar(50) not null"`
	Quantity          int       `xorm:"int not null default 0"`
	CreatedAt         time.Time `xorm:"datetime"`
	BrandCode         string    `xorm:"varchar(50)"`
}

//StockDistributeGroup ...
type StockDistributeGroup struct {
	StockDistribute     `xorm:"extends"`
	StockDistributeItem `xorm:"extends"`
}

// TableName 设置对应的表名
func (StockDistributeGroup) TableName() string {
	return "stock_distribute_item"
}
