package models

// Sku 单位库存商品
type Sku struct {
	ID        int64  `xorm:"bigint not null autoincr pk 'id'"`
	ProductID int64  `xorm:"bigint not null 'product_id'"`
	Code      string `xorm:"varchar(50)"`
	Enable    bool   `xorm:"tinyint(1)"`
}

// TableName 设置对应的表名
func (Sku) TableName() string {
	return "pangpang_brand_product.sku"
}
