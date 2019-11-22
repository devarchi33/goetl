package repositories

import (
	"clearance-adapter/config"
	"clearance-adapter/domain/entities"
	p2bConst "clearance-adapter/domain/p2brand-constants"
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"clearance-adapter/infra/proxy"
	"fmt"
	"log"
	"net/http"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
)

// SkuLocationDistributionType sku-location 使用的入库类型
// https://pangpang.dooray.com/project/posts/2565169056341902087
type SkuLocationDistributionType = string

const (
	// TypStockDistribute 物流分配
	TypStockDistribute = "Distribute"

	// TypDirectDistribute 工厂直送
	TypDirectDistribute = "DirectDistribute"
)

// StockDistributionRepository P2-brand 物流分配入库单仓库
type StockDistributionRepository struct{}

// GetUnsyncedDistributionOrders 获取未同步过的分配入库数据
func (StockDistributionRepository) GetUnsyncedDistributionOrders() ([]map[string]string, error) {
	sql := `
		SELECT
			sd.id AS distribution_id,
			sdi.brand_code,
			store.code AS receipt_location_code,
			sd.waybill_no,
			sd.box_no,
			sku.code AS sku_code,
			sdi.quantity AS qty,
			sd.created_at AS in_date,
			emp.emp_id AS in_emp_id,
			'STOCK_DISTRIBUTION' AS type,
			sd.is_auto_receipt
		FROM pangpang_brand_sku_location.stock_distribute AS sd
			JOIN pangpang_brand_sku_location.stock_distribute_item AS sdi
				ON sd.id = sdi.stock_distribute_id
			JOIN pangpang_brand_product.sku AS sku
				ON sku.id = sdi.sku_id
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = sd.receipt_location_id
			JOIN pangpang_common_colleague_auth.employees AS emp
				ON emp.colleague_id = sd.colleague_id
		WHERE sd.tenant_code = ?
			AND sd.synced = false
			AND sd.is_auto_receipt = 0
			;
	`

	result, err := factory.GetP2BrandEngine().Query(sql, getTenantCode())
	if err != nil {
		return nil, err
	}

	return infra.ConvertByteResult(result), nil
}

// GetUnsyncedAutoDistributionOrders 获取未同步过的分配入库数据
func (StockDistributionRepository) GetUnsyncedAutoDistributionOrders() ([]map[string]string, error) {
	sql := `
		SELECT
			sd.id AS distribution_id,
			sdi.brand_code,
			store.code AS receipt_location_code,
			sd.waybill_no,
			sd.box_no,
			sku.code AS sku_code,
			sdi.quantity AS qty,
			sd.created_at AS in_date,
			'auto' AS in_emp_id,
			'STOCK_DISTRIBUTION' AS type,
			sd.is_auto_receipt
		FROM pangpang_brand_sku_location.stock_distribute AS sd
			JOIN pangpang_brand_sku_location.stock_distribute_item AS sdi
				ON sd.id = sdi.stock_distribute_id
			JOIN pangpang_brand_product.sku AS sku
				ON sku.id = sdi.sku_id
			JOIN pangpang_brand_place_management.store AS store
				ON store.id = sd.receipt_location_id
		WHERE sd.tenant_code = ?
			AND sd.synced = false
			AND sd.is_auto_receipt = 1
			;
	`

	result, err := factory.GetP2BrandEngine().Query(sql, getTenantCode())
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
			WHERE sd.tenant_code = ?
			AND store.code = ?
			AND sd.waybill_no = ?
		;
	`
	_, err := factory.GetP2BrandEngine().Exec(sql, getTenantCode(), receiptLocationCode, waybillNo)
	if err != nil {
		log.Printf("MarkWaybillSynced error: %v", err)
		return err
	}

	return nil
}

func (StockDistributionRepository) putInStorage(order entities.DistributionOrder, isAutoReceipt bool) error {
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
	data["distributionType"] = TypStockDistribute
	if order.Type == p2bConst.TypFactoryToShop {
		data["distributionType"] = TypDirectDistribute
	}
	data["isAutoReceipt"] = isAutoReceipt

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
	if statusCode == 201 {
		return nil
	}

	return fmt.Errorf("Call PutInStorage API error, status code is %v, error message is: %v , data is: %v", statusCode, resp, data)
}

// PutInStorage P2Brand 入库
func (StockDistributionRepository) PutInStorage(order entities.DistributionOrder, isAutoReceipt bool) error {
	return StockDistributionRepository{}.putInStorage(order, isAutoReceipt)
}
