package cslConstants

// StockMissStatus 误差状态
type StockMissStatus = string

const (
	// StsSentOut 已出库状态
	StsSentOut string = "R"

	// StsConfirmed 已入库确认状态
	StsConfirmed string = "F"

	// StsShopRegistered 卖场登记
	StsShopRegistered StockMissStatus = "10"

	// StsLogisticsReceived 物流接受
	StsLogisticsReceived StockMissStatus = "11"

	// StsStockMissFinished 处理结束
	StsStockMissFinished StockMissStatus = "12"
)
