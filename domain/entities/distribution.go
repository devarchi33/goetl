package entities

import (
	"clearance-adapter/infra"
	"errors"
)

// Distribution 以StockDistribute为主的信息，会把id替换成code
type Distribution struct {
	BrandCode string
	ShopCode  string
	WaybillNo string
	BoxNo     string
	EmpID     string
	Items     []DistributionItem
}

// DistributionItem 以StockDistributeItem为主的信息，会把id替换成code
type DistributionItem struct {
	SkuCode string
	Qty     int
}

// Create 根据[]map类型的数据转换成Distribution
func (Distribution) Create(data []map[string]string) (Distribution, error) {
	distribution := Distribution{}
	if data == nil || len(data) == 0 {
		return distribution, errors.New("data is empty")
	}

	distData := data[0]
	err := Distribution{}.checkRequirement(distData, "brand_code", "shop_code", "waybill_no", "box_no", "emp_id")
	if err != nil {
		return distribution, err
	}

	distribution.BrandCode = distData["brand_code"]
	distribution.ShopCode = distData["shop_code"]
	distribution.WaybillNo = distData["waybill_no"]
	distribution.BoxNo = distData["box_no"]
	distribution.EmpID = distData["emp_id"]
	distribution.Items = make([]DistributionItem, 0)

	for _, item := range data {
		err := Distribution{}.checkRequirement(distData, "sku_code", "qty")
		if err != nil {
			return distribution, err
		}
		distribution.Items = append(distribution.Items, DistributionItem{
			SkuCode: item["sku_code"],
			Qty:     infra.ConvertStringToInt(item["qty"]),
		})
	}

	return distribution, nil
}

func (Distribution) checkRequirement(data map[string]string, requiredKeys ...string) error {
	for _, key := range requiredKeys {
		if _, ok := data[key]; !ok {
			return errors.New(key + " is required")
		}
	}
	return nil
}
