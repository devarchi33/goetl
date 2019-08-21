package repositories

import (
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"log"
)

// StockRoundRepository P2-brand 调货单仓库
type StockRoundRepository struct{}

// GetUnsyncedTransferInOrders 获取未同步过的已调货入库数据
func (StockRoundRepository) GetUnsyncedTransferInOrders() ([]map[string]string, error) {
	sql := `
		SELECT
			sri.brand_code,
			shipmentStore.code AS shipment_location_code,
			receitpStore.code AS receipt_location_code,
			sr.waybill_no,
			sr.box_no,
			sku.code AS sku_code,
			sri.quantity AS qty,
			sr.out_created_at AS out_date,
			outEmp.employee_no AS out_emp_id,
			sr.in_created_at AS in_date,
			inEmp.employee_no AS in_emp_id,
			sr.shipping_company_code
		FROM pangpang_brand_sku_location.stock_round AS sr
			JOIN pangpang_brand_sku_location.stock_round_item AS sri
				ON sr.id = sri.stock_round_id
			JOIN pangpang_brand_product.sku AS sku
				ON sku.id = sri.sku_id
			JOIN pangpang_brand_place_management.store AS shipmentStore
				ON shipmentStore.id = sr.shipment_location_id
			JOIN pangpang_brand_place_management.store AS receitpStore
				ON receitpStore.id = sr.receipt_location_id
			JOIN pangpang_common_colleague_employee.employees AS outEmp
				ON outEmp.id = sr.out_colleague_id
			JOIN pangpang_common_colleague_employee.employees AS inEmp
				ON inEmp.id = sr.in_colleague_id
		WHERE sr.tenant_code = ?
			AND sr.synced = false
			AND sr.status = 'F'
			;
	`

	result, err := factory.GetP2BrandEngine().Query(sql, getTenantCode())
	if err != nil {
		return nil, err
	}

	return infra.ConvertByteResult(result), nil
}

// MarkWaybillSynced 标记某个运单已经同步过
func (StockRoundRepository) MarkWaybillSynced(shipmentLocationCode, receiptLocationCode, waybillNo string) error {
	sql := `
		UPDATE pangpang_brand_sku_location.stock_round sr
		JOIN pangpang_brand_place_management.store AS shipmentStore
			ON shipmentStore.id = sr.shipment_location_id
		JOIN pangpang_brand_place_management.store AS receitpStore
			ON receitpStore.id = sr.receipt_location_id
		SET sr.synced = true,
			sr.last_updated_at = now()
			WHERE sr.tenant_code = ?
			AND shipmentStore.code = ?
			AND receitpStore.code = ?
			AND sr.waybill_no = ?
		;
	`
	_, err := factory.GetP2BrandEngine().Exec(sql, getTenantCode(), shipmentLocationCode, receiptLocationCode, waybillNo)
	if err != nil {
		log.Printf("MarkWaybillSynced error: %v", err)
		return err
	}

	return nil
}
