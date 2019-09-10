package clrConstants

const (
	// TypStockDistributionError 手动入库，包括工厂直送和物流分配入库错误
	TypStockDistributionError string = "STOCK_DISTRIBUTION_ERROR"

	// TypAutoStockDistributionError 自动入库，包括工厂直送和物流分配入库错误
	TypAutoStockDistributionError string = "AUTO_STOCK_DISTRIBUTION_ERROR"

	// TypReturnToWarehouseError 退仓错误
	TypReturnToWarehouseError string = "RETURN_TO_WAREHOUSE_ERROR"
)
