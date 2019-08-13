package services

import (
	"clearance-adapter/domain/entities"
	"clearance-adapter/repositories"
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/pangpanglabs/goetl"
)

// ReturnToWarehouseETL 退仓 p2-brand -> CSL
type ReturnToWarehouseETL struct {
	StartDateTime time.Time
	EndDateTime   time.Time
}

// New 创建 ReturnToWarehouseETL 对象，从Clearance到CSL
func (ReturnToWarehouseETL) New(startDatetime, endDateTime string) *goetl.ETL {
	local, _ := time.LoadLocation("Local")
	start, _ := time.ParseInLocation("2006-01-02 15:04:05", startDatetime, local)
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", endDateTime, local)
	returnToWarehouseETL := ReturnToWarehouseETL{
		StartDateTime: start,
		EndDateTime:   end,
	}

	etl := goetl.New(returnToWarehouseETL)
	// etl.Before(ReturnToWarehouseETL{}.buildReturnToWarehouseOrders)
	// etl.Before(ReturnToWarehouseETL{}.filterStorableDistributions)

	return etl
}

// Extract ...
func (etl ReturnToWarehouseETL) Extract(ctx context.Context) (interface{}, error) {
	result, err := repositories.ReturnToWarehouseRepository{}.GetReturnToWarehouseOrdersByCreateAt(etl.StartDateTime, etl.EndDateTime)
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

		recvSuppNo, err := repositories.RecvSuppRepository{}.CreateReturnToWarehouseOrder(order.BrandCode, order.ShipmentLocationCode, order.WaybillNo, order.EmpID, order.DeliveryOrderNo)
		if err != nil {
			log.Printf(err.Error())
		}

		var wg sync.WaitGroup
		wg.Add(len(order.Items))
		for _, v := range order.Items {
			go func(item entities.ReturnToWarehouseOrderItem, wg *sync.WaitGroup) {
				err := repositories.RecvSuppRepository{}.AddReturnToWarehouseOrderItem(order.BrandCode, order.ShipmentLocationCode, recvSuppNo, item.SkuCode, item.Qty, order.EmpID)
				if err != nil {
					log.Printf(err.Error())
				}
				wg.Done()
			}(v, &wg)
		}
		wg.Wait()
	}

	return nil
}
