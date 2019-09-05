package models

import "time"

// StockDistributionError 入库错误信息
type StockDistributionError struct {
	ID             int64     `xorm:"bigint not null autoincr pk 'id'"`
	DistributionID int64     `xorm:"bigint 'distribution_id' not null default '0'"`
	WaybillNo      string    `xorm:"varchar(30) not null default ''"`
	Type           string    `xorm:"varchar(20) not null default ''"`
	ErrorMessage   string    `xorm:"varchar(500) not null default ''"`
	IsProcessed    bool      `xorm:"tinyint(1) not null default '0'"`
	CreatedAt      time.Time `xorm:"datetime not null default now()"`
	CreatedBy      string    `xorm:"varchar(20) not null default ''"`
}
