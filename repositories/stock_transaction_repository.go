package repositories

import (
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"time"
)

// StockTransactionRepository P2-brand-sku-location的仓库
type StockTransactionRepository struct{}

// GetInStorageByCreateAt 获取某一时间段的入库数据
func (StockTransactionRepository) GetInStorageByCreateAt(start, end time.Time) ([]map[string]string, error) {
	sql := `
			SELECT
			txi.brand_code,
			store.code AS shop_code,
			tx.waybill_no,
			tx.box_no,
			sku.code AS sku_code,
			txi.quantity AS qty,
			'shi.yanxun' AS user_id
		FROM pangpang_brand_sku_location.stock_transaction AS tx
			JOIN pangpang_brand_sku_location.stock_transaction_item AS txi
				ON tx.id = txi.stock_transaction_id
			JOIN pangpang_brand_product.sku AS sku
				ON sku.id = txi.sku_id
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = tx.receipt_location_id
		WHERE tx.type = 'IN'
			AND tx.created_at BETWEEN ? AND ?;
	`
	result, err := factory.GetP2BrandEngine().Query(sql, start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}

	return infra.ConvertByteResult(result), nil
}
