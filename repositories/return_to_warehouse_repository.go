package repositories

import (
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"time"
)

// ReturnToWarehouseRepository p2-brand 卖场退仓单仓库
type ReturnToWarehouseRepository struct{}

// GetReturnToWarehouseOrdersByCreateAt 获取某一时间段的退仓数据
func (ReturnToWarehouseRepository) GetReturnToWarehouseOrdersByCreateAt(start, end time.Time) ([]map[string]string, error) {
	sql := `
		SELECT
			rtwi.brand_code,
			store.code AS shipment_location_code,
			rtw.waybill_no,
			rtw.status AS status_code,
			sku.code AS sku_code,
			rtwi.quantity AS qty,
			emp.employee_no AS emp_id
		FROM pangpang_brand_sku_location.return_to_warehouse AS rtw
			JOIN pangpang_brand_sku_location.return_to_warehouse_item AS rtwi
				ON rtw.id = rtwi.return_to_warehouse_id
			JOIN pangpang_brand_product.sku AS sku
				ON sku.id = rtwi.sku_id
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = rtw.shipment_location_id
			JOIN pangpang_common_colleague_employee.employees AS emp
				ON emp.id = rtw.colleague_id
		WHERE rtw.tenant_code = 'pangpang'
			AND rtw.status = 'R'
			AND rtw.created_at BETWEEN ? AND ?
		;
	`

	result, err := factory.GetP2BrandEngine().Query(sql, start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}

	return infra.ConvertByteResult(result), nil
}
