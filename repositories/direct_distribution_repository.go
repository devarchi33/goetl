package repositories

import (
	"clearance-adapter/config"
	"clearance-adapter/domain/entities"
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"clearance-adapter/infra/proxy"
	"fmt"
	"log"
	"net/http"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
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
			emp.employee_no AS in_emp_id,
			'DIRECT_DISTRIBUTION' AS type
		FROM pangpang_brand_sku_location.direct_distribution AS dd
			JOIN pangpang_brand_sku_location.direct_distribution_item AS ddi
				ON dd.id = ddi.stock_distribute_id
			JOIN pangpang_brand_product.sku AS sku
				ON sku.id = ddi.sku_id
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = dd.receipt_location_id
			JOIN pangpang_common_colleague_employee.employees AS emp
				ON emp.id = dd.colleague_id
		WHERE dd.tenant_code = ?
			AND dd.synced = false
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
func (DirectDistributionRepository) PutInStorage(order entities.DistributionOrder) error {
	store, has, err := PlaceRepository{}.GetStoreByCode(order.ReceiptLocationCode)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("PutInStorage error: 没有找到%v卖场", order.ReceiptLocationCode)
	}

	data := make(map[string]interface{}, 0)
	data["receiptLocationId"] = store.ID
	data["waybillNo"] = order.WaybillNo
	data["boxNo"] = order.BoxNo
	data["brandCode"] = order.BrandCode
	data["version"] = order.Version
	data["distributionType"] = "16"

	items := make([]map[string]interface{}, 0)
	for _, v := range order.Items {
		sku, _, err := ProductRepository{}.GetSkuByCode(v.SkuCode)
		if err != nil {
			return err
		}
		item := make(map[string]interface{}, 0)
		item["productId"] = sku.ProductID
		item["skuId"] = sku.ID
		item["quantity"] = v.Qty
		item["outQuantity"] = v.Qty
		item["barcode"] = v.SkuCode
		item["brandCode"] = order.BrandCode
		items = append(items, item)
	}
	data["items"] = items

	url := config.GetP2BrandSkuLocationAPIRoot() + "/stock-distribute"
	fmt.Println(url)

	fmt.Println("--------")
	fmt.Println(data)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json;charset=utf-8"
	var resp struct {
		Result  interface{} `json:"result"`
		Success bool        `json:"success"`
		Error   struct {
			Code    int         `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details"`
		} `json:"error"`
	}
	tokenProxy := proxy.TokenProxy{}.GetInstance()
	token, err := tokenProxy.GetToken()
	if err != nil {
		return err
	}
	statusCode, err := httpreq.New(http.MethodPost, url, []map[string]interface{}{data}).
		WithToken(token).
		WithBehaviorLogContext(behaviorlog.FromCtx(nil)).
		Call(&resp)

	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	if statusCode == 201 {
		return nil
	}

	return fmt.Errorf("Call PutInStorage API error, status code is %v, error message is: %v ", statusCode, resp)
}
