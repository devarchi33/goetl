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
			'CEGP' AS shop_code,
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
		WHERE tx.waybill_no = '1010590009008'
			AND tx.type = 'IN';
	`
	result, err := factory.GetP2BrandEngine().Query(sql)
	if err != nil {
		return nil, err
	}

	return infra.ConvertByteResult(result), nil
}
