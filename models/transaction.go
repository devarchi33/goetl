package models

// TransactionMaster ...
type TransactionMaster struct {
	ID              int64  `xorm:"bigint not null autoincr pk 'id'"`
	Date            string `xorm:"varchar(8) not null"`
	PlantCode       string `xorm:"varchar(50) null default null"`
	WaybillNo       string `xorm:"varchar(50) null default null"`
	OrderNo         string `xorm:"varchar(50) null default null"`
	TransactionCode string `xorm:"varchar(50) null default null"`
	Channel         string `xorm:"varchar(20) null default null"`
}

// TableName 设置对应的表名
func (TransactionMaster) TableName() string {
	return "transaction_masters"
}

// TransactionDetail ...
type TransactionDetail struct {
	ID                  int64  `xorm:"bigint not null autoincr pk 'id'"`
	TransactionMasterID int64  `xorm:"bigint not null 'transaction_master_id'"`
	SkuCode             string `json:"sku_code" xorm:"varchar(200) not null"`
	Qty                 int    `json:"qty" xorm:"int not null"`
}

// TableName 设置对应的表名
func (TransactionDetail) TableName() string {
	return "transaction_details"
}

// Transaction ...
type Transaction struct {
	ID            int64  `xorm:"bigint not null autoincr pk 'id'"`
	TransactionID string `xorm:"varchar(14) not null 'transaction_id'"`
	WaybillNo     string `xorm:"varchar(13) not null"`
	BoxNo         string `xorm:"varchar(20) not null"`
	SkuCode       string `xorm:"varchar(18) not null"`
	Qty           int    `xorm:"int not null"`
}

// TableName 设置对应的表名
func (Transaction) TableName() string {
	return "transactions"
}
