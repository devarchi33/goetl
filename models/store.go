package models

// Store 卖场
type Store struct {
	ID         int64  `xorm:"bigint not null autoincr pk 'id'"`
	TenantCode string `xorm:"varchar(50)"`
	Code       string `xorm:"varchar(50)"`
	Enable     bool   `xorm:"tinyint(1)"`
}

// TableName 设置对应的表名
func (Store) TableName() string {
	return "pangpang_brand_place_management.store"
}
