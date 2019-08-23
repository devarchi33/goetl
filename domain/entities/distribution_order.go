package entities

import (
	"clearance-adapter/infra"
	"errors"
)

// DistributionOrder 以StockDistribute为主的信息，会把id替换成code
type DistributionOrder struct {
	BrandCode string
	ShopCode  string
	WaybillNo string
	BoxNo     string
	InDate    string
	EmpID     string
	Items     []DistributionOrderItem
}

// DistributionOrderItem 以StockDistributeItem为主的信息，会把id替换成code
type DistributionOrderItem struct {
	SkuCode string
	Qty     int
}

// Create 根据[]map类型的数据转换成 DistributionOrder
func (DistributionOrder) Create(data []map[string]string) (DistributionOrder, error) {
	order := DistributionOrder{}
	if data == nil || len(data) == 0 {
		return order, errors.New("data is empty")
	}

	orderData := data[0]
	err := checkRequirement(orderData, "brand_code", "shop_code", "waybill_no", "box_no", "in_date", "emp_id")
	if err != nil {
		return order, err
	}

	order.BrandCode = orderData["brand_code"]
	order.ShopCode = orderData["shop_code"]
	order.WaybillNo = orderData["waybill_no"]
	order.BoxNo = orderData["box_no"]
	order.InDate = infra.Parse8BitsDate(orderData["in_date"], nil)
	order.EmpID = orderData["emp_id"]
	order.Items = make([]DistributionOrderItem, 0)

	for _, item := range data {
		err := checkRequirement(orderData, "sku_code", "qty")
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
