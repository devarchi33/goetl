package models

import "time"

// StockTransaction ...
type StockTransaction struct {
	ID                 int64     `xorm:"bigint not null autoincr pk 'id'"`
	Type               string    `xorm:"varchar(8) not null"`
	CreatedAt          time.Time `xorm:"datetime"`
	TenantCode         string    `xorm:"varchar(50)"`
	WaybillNo          string    `xorm:"varchar(50) not null default ''"`
	ColleagueID        int64     `xorm:"bigint 'colleague_id'"`
	ShipmentLocationID int64     `xorm:"bigint 'shipment_location_Id'"`
	ReceiptLocationID  int64     `xorm:"bigint 'receipt_location_id'"`
	BrandCode          string    `xorm:"varchar(50)"`
	BoxNo              string    `xorm:"varchar(50) not null default ''"`
}

// StockTransactionItem ...
type StockTransactionItem struct {
	ID                 int64     `xorm:"bigint not null autoincr pk 'id'"`
	StockTransactionID int64     `xorm:"bigint not null 'stock_transaction_id'"`
	ProductID          int64     `xorm:"bigint not null 'product_id'"`
	SkuID              int64     `xorm:"bigint not null 'sku_id'"`
	Quantity           int       `xorm:"int not null default 0"`
	CreatedAt          time.Time `xorm:"datetime"`
	Barcode            string    `xorm:"varchar(50) not null"`
	BrandCode          string    `xorm:"varchar(50)"`
}

//StockTransactionGroup ...
type StockTransactionGroup struct {
	StockTransaction     `xorm:"extends"`
	StockTransactionItem `xorm:"extends"`
}

// TableName 设置对应的表名
func (StockTransactionGroup) TableName() string {
	return "stock_transaction_item"
}
