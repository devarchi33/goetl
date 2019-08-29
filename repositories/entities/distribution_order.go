package entities

// DistributionOrder 以StockDistribute为主的信息，会把id替换成code
type DistributionOrder struct {
	BrandCode           string
	ReceiptLocationCode string
	WaybillNo           string
	BoxNo               string
	Version             string
	Items               []DistributionOrderItem
}

// DistributionOrderItem 以 StockDistributeItem 为主的信息，会把id替换成code
type DistributionOrderItem struct {
	SkuCode string
	Qty     int
}
