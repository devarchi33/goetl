package repositories

import (
	"clearance-adapter/domain/entities"
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"log"
)

// DirectDistributionRepository P2-brand 工厂直送分配入库单仓库
type DirectDistributionRepository struct{}

// GetUnsyncedDistributionOrders 获取未同步过的分配入库数据
func (DirectDistributionRepository) GetUnsyncedDistributionOrders() ([]map[string]string, error) {
	sql := `
		SELECT
			dd.id AS distribution_id,
			ddi.brand_code,
			store.code AS receipt_location_code,
			dd.waybill_no,
			dd.box_no,
			sku.code AS sku_code,
			ddi.quantity AS qty,
			dd.created_at AS in_date,
			emp.emp_id AS in_emp_id,
			'DIRECT_DISTRIBUTION' AS type,
			dd.is_auto_receipt
		FROM pangpang_brand_sku_location.direct_distribution AS dd
			JOIN pangpang_brand_sku_location.direct_distribution_item AS ddi
				ON dd.id = ddi.stock_distribute_id
			JOIN pangpang_brand_product.sku AS sku
				ON sku.id = ddi.sku_id
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = dd.receipt_location_id
			JOIN pangpang_common_colleague_auth.employees AS emp
				ON emp.colleague_id = dd.colleague_id
		WHERE dd.tenant_code = ?
			AND dd.synced = false
			AND dd.is_auto_receipt = 0
			;
	`

	result, err := factory.GetP2BrandEngine().Query(sql, getTenantCode())
	if err != nil {
		return nil, err
	}

	return infra.ConvertByteResult(result), nil
}

// GetUnsyncedAutoDistributionOrders 获取未同步过的分配入库数据
func (DirectDistributionRepository) GetUnsyncedAutoDistributionOrders() ([]map[string]string, error) {
	sql := `
		SELECT
			dd.id AS distribution_id,
			ddi.brand_code,
			store.code AS receipt_location_code,
			dd.waybill_no,
			dd.box_no,
			sku.code AS sku_code,
			ddi.quantity AS qty,
			dd.created_at AS in_date,
			'auto' AS in_emp_id,
			'DIRECT_DISTRIBUTION' AS type,
			dd.is_auto_receipt
		FROM pangpang_brand_sku_location.direct_distribution AS dd
			JOIN pangpang_brand_sku_location.direct_distribution_item AS ddi
				ON dd.id = ddi.stock_distribute_id
			JOIN pangpang_brand_product.sku AS sku
				ON sku.id = ddi.sku_id
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = dd.receipt_location_id
		WHERE dd.tenant_code = ?
			AND dd.synced = false
			AND dd.is_auto_receipt = 1
			;
	`

	result, err := factory.GetP2BrandEngine().Query(sql, getTenantCode())
	if err != nil {
		return nil, err
	}

	return infra.ConvertByteResult(result), nil
}

// MarkWaybillSynced 标记某个运单已经同步过
func (DirectDistributionRepository) MarkWaybillSynced(receiptLocationCode, waybillNo string) error {
	sql := `
		UPDATE pangpang_brand_sku_location.direct_distribution dd
		JOIN pangpang_brand_place_management.store AS store
			ON store.id = dd.receipt_location_id
		SET dd.synced = true,
			dd.last_updated_at = now()
			WHERE dd.tenant_code = ?
			AND store.code = ?
			AND dd.waybill_no = ?
		;
	`
	_, err := factory.GetP2BrandEngine().Exec(sql, getTenantCode(), receiptLocationCode, waybillNo)
	if err != nil {
		log.Printf("MarkWaybillSynced error: %v", err)
		return err
	}

	return nil
}

// PutInStorage P2Brand 入库
func (DirectDistributionRepository) PutInStorage(order entities.DistributionOrder, isAutoReceipt bool) error {
	return StockDistributionRepository{}.putInStorage(order, isAutoReceipt)
}
