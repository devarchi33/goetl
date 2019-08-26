package entities

import (
	"clearance-adapter/infra"
	"errors"
)

// TransferOrder 以StockRound为主的信息，会把id替换成code
type TransferOrder struct {
	BrandCode            string
	ShipmentLocationCode string
	ReceiptLocationCode  string
	WaybillNo            string
	BoxNo                string
	OutDate              string
	OutEmpID             string
	InDate               string
	InEmpID              string
	ShippingCompanyCode  string
	Items                []TransferOrderItem
}

// TransferOrderItem 以StockRoundItem为主的信息，会把id替换成code
type TransferOrderItem struct {
	SkuCode string
	Qty     int
}

// Create 根据[]map类型的数据转换成 TransferOrder
func (TransferOrder) Create(data []map[string]string) (TransferOrder, error) {
	order := TransferOrder{}
	if data == nil || len(data) == 0 {
		return order, errors.New("data is empty")
	}

	orderData := data[0]
	err := checkRequirement(orderData, "brand_code", "shipment_location_code", "receipt_location_code", "waybill_no", "box_no", "out_date", "out_emp_id", "in_date", "in_emp_id", "shipping_company_code", "sku_code", "qty")
	if err != nil {
		return order, err
	}

	order.BrandCode = orderData["brand_code"]
	order.ShipmentLocationCode = orderData["shipment_location_code"]
	order.ReceiptLocationCode = orderData["receipt_location_code"]
	order.WaybillNo = orderData["waybill_no"]
	order.BoxNo = orderData["box_no"]
	order.OutEmpID = orderData["out_emp_id"]
	order.OutDate = infra.Parse8BitsDate(orderData["out_date"], nil)
	order.InEmpID = orderData["in_emp_id"]
	order.InDate = infra.Parse8BitsDate(orderData["in_date"], nil)
	order.ShippingCompanyCode = orderData["shipping_company_code"]
	order.Items = make([]TransferOrderItem, 0)

	for _, item := range data {
		err := checkRequirement(orderData, "sku_code", "qty")
		if err != nil {
			return order, err
		}
		order.Items = append(order.Items, TransferOrderItem{
			SkuCode: item["sku_code"],
			Qty:     infra.ConvertStringToInt(item["qty"]),
		})
	}

	return order, nil
}
