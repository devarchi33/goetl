package entities

import (
	"clearance-adapter/infra"
	"errors"
)

// DistributionOrder 以StockDistribute为主的信息，会把id替换成code
type DistributionOrder struct {
	BrandCode           string
	ReceiptLocationCode string
	WaybillNo           string
	BoxNo               string
	InDate              string
	InEmpID             string
	Items               []DistributionOrderItem
}

// RequiredKeys 创建 DistributionOrder 的必须项
func (DistributionOrder) RequiredKeys() []string {
	return []string{"brand_code", "receipt_location_code", "waybill_no", "box_no", "in_date", "in_emp_id"}
}

// Create 根据[]map类型的数据转换成 DistributionOrder
func (DistributionOrder) Create(data []map[string]string) (DistributionOrder, error) {
	order := DistributionOrder{}
	if data == nil || len(data) == 0 {
		return order, errors.New("data is empty")
	}

	orderData := data[0]
	err := checkRequirement(orderData, DistributionOrder{}.RequiredKeys()...)
	if err != nil {
		return order, err
	}

	order.BrandCode = orderData["brand_code"]
	order.ReceiptLocationCode = orderData["receipt_location_code"]
	order.WaybillNo = orderData["waybill_no"]
	order.BoxNo = orderData["box_no"]
	order.InDate = infra.Parse8BitsDate(orderData["in_date"], nil)
	order.InEmpID = orderData["in_emp_id"]
	order.Items = make([]DistributionOrderItem, 0)

	for _, item := range data {
		err := checkRequirement(orderData, DistributionOrderItem{}.RequiredKeys()...)
		if err != nil {
			return order, err
		}
		order.Items = append(order.Items, DistributionOrderItem{
			SkuCode: item["sku_code"],
			Qty:     infra.ConvertStringToInt(item["qty"]),
		})
	}

	return order, nil
}

// DistributionOrderItem 以 StockDistributeItem 为主的信息，会把id替换成code
type DistributionOrderItem struct {
	SkuCode string
	Qty     int
}

// RequiredKeys 创建 DistributionOrderItem 的必须项
func (DistributionOrderItem) RequiredKeys() []string {
	return []string{"sku_code", "qty"}
}
