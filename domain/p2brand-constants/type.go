package p2brandConstants

// DistributionType 入库单类型
type DistributionType = string

const (
	// TypLogisticsToShop 从物流到卖场，卖场分配
	TypLogisticsToShop DistributionType = "STOCK_DISTRIBUTION"

	// TypFactoryToShop 从工厂到卖场，工厂直送
	TypFactoryToShop DistributionType = "DIRECT_DISTRIBUTION"
)
