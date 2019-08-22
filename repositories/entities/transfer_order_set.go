package entities

// TransferOrderSet 用于在CSL同时创建调货出库单和调货入库单的实体
type TransferOrderSet struct {
	BrandCode            string
	ShipmentLocationCode string // 发货卖场主卖场
	ReceiptLocationCode  string // 收货卖场主卖场
	ShipmentShopCode     string // 发货卖场子卖场
	ReceiptShopCode      string // 收货卖场子卖场
	WaybillNo            string
	BoxNo                string
	ShippingCompanyCode  string
	DeliveryOrderNo      string
	OutDate              string
	InDate               string
	OutEmpID             string
	InEmpID              string
	Items                []TransferOrderSetItem
}

// TransferOrderSetItem 调出单中的商品
type TransferOrderSetItem struct {
	SkuCode string
	Qty     int
}
