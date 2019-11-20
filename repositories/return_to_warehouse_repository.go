package repositories

import (
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"log"
)

// ReturnToWarehouseRepository p2-brand 卖场退仓单仓库
type ReturnToWarehouseRepository struct{}

// GetUnsyncedReturnToWarehouseAnytimeOrders 获取未同步过的随时退仓数据
func (ReturnToWarehouseRepository) GetUnsyncedReturnToWarehouseAnytimeOrders() ([]map[string]string, error) {
	sql := `
		SELECT
			rtwi.brand_code,
			store.code AS shipment_location_code,
			rtw.waybill_no,
			rtw.status AS status_code,
			sku.code AS sku_code,
			rtwi.quantity AS qty,
			emp.emp_id AS emp_id,
			'' AS defective_prod_reason_code,
			rtw.out_created_at AS out_date
		FROM pangpang_brand_sku_location.return_to_warehouse_anytime AS rtw
			JOIN pangpang_brand_sku_location.return_to_warehouse_anytime_item AS rtwi
				ON rtw.id = rtwi.anytime_id
			JOIN pangpang_brand_product.sku AS sku
				ON sku.id = rtwi.sku_id
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = rtw.shipment_location_id
			JOIN pangpang_common_colleague_auth.employees AS emp
				ON emp.colleague_id = rtw.colleague_id
		WHERE rtw.tenant_code = ?
			AND rtw.status = 'R'
			AND rtw.synced = false
		;
	`

	result, err := factory.GetP2BrandEngine().Query(sql, getTenantCode())
	if err != nil {
		return nil, err
	}

	return infra.ConvertByteResult(result), nil
}

// MarkAnytimeWaybillSynced 标记某个随时退仓运单已经同步过
func (ReturnToWarehouseRepository) MarkAnytimeWaybillSynced(shipmentLocationCode, waybillNo string) error {
	sql := `
		UPDATE pangpang_brand_sku_location.return_to_warehouse_anytime rtw
		JOIN pangpang_brand_place_management.store AS store
			ON store.id = rtw.shipment_location_id
		SET rtw.synced = true,
			rtw.last_updated_at = now()
			WHERE rtw.tenant_code = ?
			AND store.code = ?
			AND rtw.waybill_no = ?
		;
	`
	_, err := factory.GetP2BrandEngine().Exec(sql, getTenantCode(), shipmentLocationCode, waybillNo)
	if err != nil {
		log.Printf("MarkAnytimeWaybillSynced error: %v", err)
		return err
	}

	return nil
}

// GetUnsyncedReturnToWarehouseSeasonalOrders 获取未同步过的季节退仓数据
func (ReturnToWarehouseRepository) GetUnsyncedReturnToWarehouseSeasonalOrders() ([]map[string]string, error) {
	sql := `
		SELECT
			rtwi.brand_code,
			store.code AS shipment_location_code,
			rtw.waybill_no,
			rtw.status AS status_code,
			sku.code AS sku_code,
			rtwi.quantity AS qty,
			emp.emp_id AS emp_id,
			'' AS defective_prod_reason_code,
			rtw.out_created_at AS out_date
		FROM pangpang_brand_sku_location.return_to_warehouse_seasonal AS rtw
			JOIN pangpang_brand_sku_location.return_to_warehouse_seasonal_item AS rtwi
				ON rtw.id = rtwi.seasonal_id
			JOIN pangpang_brand_product.sku AS sku
				ON sku.id = rtwi.sku_id
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = rtw.shipment_location_id
			JOIN pangpang_common_colleague_auth.employees AS emp
				ON emp.colleague_id = rtw.colleague_id
		WHERE rtw.tenant_code = ?
			AND rtw.status = 'R'
			AND rtw.synced = false
		;
	`

	result, err := factory.GetP2BrandEngine().Query(sql, getTenantCode())
	if err != nil {
		return nil, err
	}

	return infra.ConvertByteResult(result), nil
}

// MarkSeasonalWaybillSynced 标记某个季节退仓运单已经同步过
func (ReturnToWarehouseRepository) MarkSeasonalWaybillSynced(shipmentLocationCode, waybillNo string) error {
	sql := `
		UPDATE pangpang_brand_sku_location.return_to_warehouse_seasonal rtw
		JOIN pangpang_brand_place_management.store AS store
			ON store.id = rtw.shipment_location_id
		SET rtw.synced = true,
			rtw.last_updated_at = now()
			WHERE rtw.tenant_code = ?
			AND store.code = ?
			AND rtw.waybill_no = ?
		;
	`
	_, err := factory.GetP2BrandEngine().Exec(sql, getTenantCode(), shipmentLocationCode, waybillNo)
	if err != nil {
		log.Printf("MarkAnytimeWaybillSynced error: %v", err)
		return err
	}

	return nil
}

// GetUnsyncedReturnToWarehouseDefectiveOrders 获取未同步过的次品退仓数据
func (ReturnToWarehouseRepository) GetUnsyncedReturnToWarehouseDefectiveOrders() ([]map[string]string, error) {
	sql := `
		SELECT
			rtwi.brand_code,
			store.code AS shipment_location_code,
			rtw.waybill_no,
			rtw.status AS status_code,
			sku.code AS sku_code,
			rtwi.quantity AS qty,
			emp.emp_id AS emp_id,
			rtwi.defective_prod_reason_code,
			rtw.out_created_at AS out_date
		FROM pangpang_brand_sku_location.return_to_warehouse_defective AS rtw
			JOIN pangpang_brand_sku_location.return_to_warehouse_defective_item AS rtwi
				ON rtw.id = rtwi.defective_id
			JOIN pangpang_brand_product.sku AS sku
				ON sku.id = rtwi.sku_id
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = rtw.shipment_location_id
			JOIN pangpang_common_colleague_auth.employees AS emp
				ON emp.colleague_id = rtw.colleague_id
		WHERE rtw.tenant_code = ?
			AND rtw.status = 'R'
			AND rtw.synced = false
		;
	`

	result, err := factory.GetP2BrandEngine().Query(sql, getTenantCode())
	if err != nil {
		return nil, err
	}

	return infra.ConvertByteResult(result), nil
}

// MarkDefectiveWaybillSynced 标记某个次品退仓运单已经同步过
func (ReturnToWarehouseRepository) MarkDefectiveWaybillSynced(shipmentLocationCode, waybillNo string) error {
	sql := `
		UPDATE pangpang_brand_sku_location.return_to_warehouse_defective rtw
		JOIN pangpang_brand_place_management.store AS store
			ON store.id = rtw.shipment_location_id
		SET rtw.synced = true,
			rtw.last_updated_at = now()
			WHERE rtw.tenant_code = ?
			AND store.code = ?
			AND rtw.waybill_no = ?
		;
	`
	_, err := factory.GetP2BrandEngine().Exec(sql, getTenantCode(), shipmentLocationCode, waybillNo)
	if err != nil {
		log.Printf("MarkAnytimeWaybillSynced error: %v", err)
		return err
	}

	return nil
}
