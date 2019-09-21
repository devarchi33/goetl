package cslConstants

const (
	// TypSentOut 已出库
	TypSentOut string = "S"

	// TypLogisticsToShop 从物流到卖场，卖场分配
	TypLogisticsToShop string = "01"

	// TypFactoryToShop 从工厂到卖场，工厂直送
	TypFactoryToShop string = "16"

	// TypAncillaryProduct 销售辅助用品入库
	TypAncillaryProduct string = "66"

	// TypShopToShop 调货
	TypShopToShop string = "20"

	// TypRTWAnytime 随时退仓
	TypRTWAnytime string = "41"

	// TypRTWSeasonal 季节退仓
	TypRTWSeasonal string = "42"

	// TypRTWDefective 次品退仓
	TypRTWDefective string = "47"
)
