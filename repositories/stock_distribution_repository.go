package repositories

import (
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"log"
)

// StockDistributionRepository P2-brand 物流分配入库单仓库
type StockDistributionRepository struct{}

// GetUnsyncedDistributionOrders 获取同步过的分配入库数据
func (StockDistributionRepository) GetUnsyncedDistributionOrders() ([]map[string]string, error) {
	sql := `
		SELECT
			sdi.brand_code,
			store.code AS shop_code,
			sd.waybill_no,
			sd.box_no,
			sku.code AS sku_code,
			sdi.quantity AS qty,
			emp.employee_no AS emp_id
		FROM pangpang_brand_sku_location.stock_distribute AS sd
			JOIN pangpang_brand_sku_location.stock_distribute_item AS sdi
				ON sd.id = sdi.stock_distribute_id
			JOIN pangpang_brand_product.sku AS sku
				ON sku.id = sdi.sku_id
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = sd.receipt_location_id
			JOIN pangpang_common_colleague_employee.employees AS emp
				ON emp.id = sd.colleague_id
		WHERE sd.tenant_code = 'pangpang'
			AND sd.synced = false
			;
	`

	result, err := factory.GetP2BrandEngine().Query(sql)
	if err != nil {
		return nil, err
	}

	return infra.ConvertByteResult(result), nil
}

// MarkWaybillSynced 标记某个运单已经同步过
func (StockDistributionRepository) MarkWaybillSynced(receiptLocationCode, waybillNo string) error {
	sql := `
		UPDATE pangpang_brand_sku_location.stock_distribute sd
		JOIN pangpang_brand_place_management.store AS store
			ON store.id = sd.receipt_location_id
		SET sd.synced = true,
			sd.last_updated_at = now()
			WHERE sd.tenant_code = 'pangpang'
			AND store.code = ?
			AND sd.waybill_no = ?
		;
	`
	_, err := factory.GetP2BrandEngine().Exec(sql, receiptLocationCode, waybillNo)
	if err != nil {
		log.Printf("MarkWaybillSynced error: %v", err)
		return err
	}

	return nil
}
