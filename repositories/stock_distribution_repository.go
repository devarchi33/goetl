package repositories

import (
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"time"
)

// StockDistributionRepository P2-brand 物流分配入库单仓库
type StockDistributionRepository struct{}

// GetDistributionOrdersByCreateAt 获取某一时间段的分配入库数据
func (StockDistributionRepository) GetDistributionOrdersByCreateAt(start, end time.Time) ([]map[string]string, error) {
	sql := `
		SELECT
			sdi.brand_code,
			store.code AS shop_code,
			sd.waybill_no,
			sd.box_no,
			sku.code AS sku_code,
			sdi.quantity AS qty,
			'7000028260' AS emp_id
		FROM pangpang_brand_sku_location.stock_distribute AS sd
			JOIN pangpang_brand_sku_location.stock_distribute_item AS sdi
				ON sd.id = sdi.stock_distribute_id
			JOIN pangpang_brand_product.sku AS sku
				ON sku.id = sdi.sku_id
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = sd.receipt_location_id
		WHERE sd.tenant_code = 'pangpang'
			AND sd.type = 'IN'
			AND sd.created_at BETWEEN ? AND ?
			;
	`

	result, err := factory.GetP2BrandEngine().Query(sql, start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}

	return infra.ConvertByteResult(result), nil
}
