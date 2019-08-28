package services

import (
	"clearance-adapter/domain/entities"
	"clearance-adapter/repositories"
	"context"
	"errors"
	"log"
	"sync"

	"github.com/pangpanglabs/goetl"
)

// ReturnToWarehouseETL 退仓 p2-brand -> CSL
type ReturnToWarehouseETL struct{}

// New 创建 ReturnToWarehouseETL 对象，从Clearance到CSL
func (ReturnToWarehouseETL) New() *goetl.ETL {
	returnToWarehouseETL := ReturnToWarehouseETL{}
	etl := goetl.New(returnToWarehouseETL)

	return etl
}

// Extract ...
func (etl ReturnToWarehouseETL) Extract(ctx context.Context) (interface{}, error) {
	result, err := repositories.ReturnToWarehouseRepository{}.GetUnsyncedReturnToWarehouseOrders()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Transform ...
func (etl ReturnToWarehouseETL) Transform(ctx context.Context, source interface{}) (interface{}, error) {
	data, ok := source.([]map[string]string)
	if !ok {
		return nil, errors.New("Convert Failed")
	}

	items := make(map[string][]map[string]string, 0)
	for _, item := range data {
		key := item["brand_code"] + "-" + item["shipment_location_code"] + "-" + item["waybill_no"]
		if _, ok := items[key]; ok {
			items[key] = append(items[key], item)
		} else {
			items[key] = make([]map[string]string, 0)
			items[key] = append(items[key], item)
		}
	}

	orders := make([]entities.ReturnToWarehouseOrder, 0)
	for _, v := range items {
		order, err := entities.ReturnToWarehouseOrder{}.Create(v)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Load ...
func (etl ReturnToWarehouseETL) Load(ctx context.Context, source interface{}) error {
	orders, ok := source.([]entities.ReturnToWarehouseOrder)
	if !ok {
		return errors.New("Convert Failed")
	}

	for _, order := range orders {
		if len(order.Items) == 0 {
			log.Printf("运单号为: %v 的出库单没有商品", order.WaybillNo)
			continue
		}

		shopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(order.ShipmentLocationCode, order.BrandCode)
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		recvSuppNo, err := repositories.RecvSuppRepository{}.CreateReturnToWarehouseOrder(order.BrandCode, shopCode, order.WaybillNo, order.OutDate, order.EmpID, order.DeliveryOrderNo)
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		var wg sync.WaitGroup
		wg.Add(len(order.Items))
		for _, v := range order.Items {
			go func(item entities.ReturnToWarehouseOrderItem, wg *sync.WaitGroup) {
				err := repositories.RecvSuppRepository{}.AddReturnToWarehouseOrderItem(order.BrandCode, shopCode, order.OutDate, recvSuppNo, item.SkuCode, item.Qty, order.EmpID)
				if err != nil {
					log.Printf(err.Error())
				}
				wg.Done()
			}(v, &wg)
		}
		wg.Wait()

		// 更新状态的时候需要使用主卖场的Code
		err = repositories.ReturnToWarehouseRepository{}.MarkWaybillSynced(order.ShipmentLocationCode, order.WaybillNo)
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		log.Printf("运单号为：%v 的退仓单（卖场：%v，品牌：%v）已经同步完成。", order.WaybillNo, order.ShipmentLocationCode, order.BrandCode)
	}

	return nil
}
