package entities

import (
	"clearance-adapter/infra"
	"errors"
)

// ReturnToWarehouseOrder 退仓出库单 以ReturnToWarehouse为主的信息，会把id替换成code
type ReturnToWarehouseOrder struct {
	BrandCode            string
	ShipmentLocationCode string
	WaybillNo            string
	StatusCode           string
	EmpID                string
	OutDate              string
	DeliveryOrderNo      string
	Items                []ReturnToWarehouseOrderItem
}

// ReturnToWarehouseOrderItem 以ReturnToWarehouseItem为主的信息，会把id替换成code
type ReturnToWarehouseOrderItem struct {
	SkuCode string
	Qty     int
}

// Create 根据[]map类型的数据转换成 ReturnToWarehouseOrder
func (ReturnToWarehouseOrder) Create(data []map[string]string) (ReturnToWarehouseOrder, error) {
	order := ReturnToWarehouseOrder{}
	if data == nil || len(data) == 0 {
		return order, errors.New("data is empty")
	}

	orderData := data[0]
	err := checkRequirement(orderData, "brand_code", "shipment_location_code", "waybill_no", "status_code", "emp_id", "out_date")
	if err != nil {
		return order, err
	}

	order.BrandCode = orderData["brand_code"]
	order.ShipmentLocationCode = orderData["shipment_location_code"]
	order.WaybillNo = orderData["waybill_no"]
	order.StatusCode = orderData["status_code"]
	order.EmpID = orderData["emp_id"]
	order.OutDate = infra.Parse8BitsDate(orderData["out_date"], nil)
	order.DeliveryOrderNo = orderData["delivery_order_no"]
	order.Items = make([]ReturnToWarehouseOrderItem, 0)

	for _, item := range data {
		err := checkRequirement(orderData, "sku_code", "qty")
		if err != nil {
			return order, err
		}
		order.Items = append(order.Items, ReturnToWarehouseOrderItem{
			SkuCode: item["sku_code"],
			Qty:     infra.ConvertStringToInt(item["qty"]),
		})
	}

	return order, nil
}
