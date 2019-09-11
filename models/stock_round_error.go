package models

import "time"

// StockRoundError 调货错误信息
type StockRoundError struct {
	ID                   int64     `xorm:"bigint not null autoincr pk 'id'"`
	BrandCode            string    `xorm:"varchar(20) not null default ''"`
	ShipmentLocationCode string    `xorm:"varchar(20) not null default ''"`
	ReceiptLocationCode  string    `xorm:"varchar(20) not null default ''"`
	WaybillNo            string    `xorm:"varchar(30) not null default ''"`
	Type                 string    `xorm:"varchar(20) not null default ''"`
	ErrorMessage         string    `xorm:"varchar(500) not null default ''"`
	IsProcessed          bool      `xorm:"tinyint(1) not null default '0'"`
	CreatedAt            time.Time `xorm:"created"`
	CreatedBy            string    `xorm:"varchar(20) not null default ''"`
}
