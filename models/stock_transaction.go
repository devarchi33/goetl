package models

// StockTransaction ...
type StockTransaction struct {
	ID                 int64  `xorm:"bigint not null autoincr pk 'id'"`
	Type               string `xorm:"varchar(8) not null"`
	OrderID            int64  `xorm:"bigint null default null"`
	TenantCode         string `xorm:"varchar(50) null default null"`
	WaybillNo          string `xorm:"varchar(50) not null"`
	BoxNo              string `xorm:"varchar(50) not null"`
	ColleagueID        int64  `xorm:"bigint"`
	ShipmentLocationID int64  `xorm:"bigint"`
	ReceiptLocationID  int64  `xorm:"bigint"`
}

func (StockTransaction) TableName() string {
	return "stock_transaction"
}

// StockTransactionItem ...
type StockTransactionItem struct {
	ID                 int64  `xorm:"bigint not null autoincr pk 'id'"`
	StockTransactionID int64  `xorm:"bigint not null 'stock_transaction_id'"`
	ProductID          int64  `xorm:"bigint not null 'product_id'"`
	SkuID              int64  `xorm:"bigint not null 'sku_id'"`
	Quantity           int    `xorm:"int not null 'quantity'"`
	Barcode            string `xorm:"varchar(50) not null 'barcode'"`
	BoxNo              string `xorm:"varchar(50) not null 'box_no'"`
}

//StockTrans ...
type StockTrans struct {
	StockTransaction     `xorm:"extends"`
	StockTransactionItem `xorm:"extends"`
}

// TableName 设置对应的表名
func (StockTrans) TableName() string {
	return "stock_transaction"
}
