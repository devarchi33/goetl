package clrConstants

// ClrErrorType Clearance 错误类型
type ClrErrorType = string

const (
	// TypStockDistributionError 手动入库，包括工厂直送和物流分配入库错误
	TypStockDistributionError ClrErrorType = "STOCK_DISTRIBUTION_ERROR"

	// TypAutoStockDistributionError 自动入库，包括工厂直送和物流分配入库错误
	TypAutoStockDistributionError ClrErrorType = "AUTO_STOCK_DISTRIBUTION_ERROR"

	// TypReturnToWarehouseError 退仓错误
	TypReturnToWarehouseError ClrErrorType = "RETURN_TO_WAREHOUSE_ERROR"

	// TypTransferError 调货错误，不区分调入还是调出，因为只有在调入后才会同步到CSL
	TypTransferError ClrErrorType = "TRANSFER_ERROR"

	// TypAutoTransferInError 自动调货入库错误
	TypAutoTransferInError ClrErrorType = "AUTO_TRANSFER_IN_ERROR"
)
