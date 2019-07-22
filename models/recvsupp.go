package models

// RecvSuppMst ...
type RecvSuppMst struct {
	RecvSuppNo         string `xorm:"char(14) not null pk "`
	BrandCode          string `xorm:"varchar(4) null default null"`
	ShopCode           string `xorm:"char(4) null default null"`
	Dates              string `xorm:"char(8) not null"`
	RecvSuppType       string `xorm:"char(1) null default null"`
	ShippingTypeCode   string `xorm:"char(2) null default null"`
	WayBillNo          string `xorm:"varchar(50) not null"`
	RecvSuppStatusCode string `xorm:"char(1) not null"`
}

// // RecvSuppDtl ...
// type RecvSuppDtl struct {
// 	RecvSuppNo       string `xorm:"char(14) not null"`
// 	BrandCode        string `xorm:"varchar(4) null default null"`
// 	ShopCode         string `xorm:"char(4) null default null"`
// 	Dates            string `xorm:"char(8) not null"`
// 	ProdCode         string `xorm:"varchar(18) null default null"`
// 	RecvSuppQty      int    `xorm:"int not null default 0"`
// 	RecvSuppFixedQty int    `xorm:"int not null default 0"`
// }
