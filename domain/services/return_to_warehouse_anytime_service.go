package services

import (
	"clearance-adapter/domain/entities"
	"clearance-adapter/repositories"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/pangpanglabs/goetl"
)

// ReturnToWarehouseAnytimeETL 随时退仓 p2-brand -> CSL
type ReturnToWarehouseAnytimeETL struct {
	ErrorRepository repositories.ReturnToWarehouseErrorRepository
}

// New 创建 ReturnToWarehouseAnytimeETL 对象，从Clearance到CSL
func (ReturnToWarehouseAnytimeETL) New() *goetl.ETL {
	returnToWarehouseETL := ReturnToWarehouseAnytimeETL{
		ErrorRepository: repositories.ReturnToWarehouseErrorRepository{}}
	etl := goetl.New(returnToWarehouseETL)

	return etl
}

func (etl ReturnToWarehouseAnytimeETL) saveError(order entities.ReturnToWarehouseOrder, errMsg string) {
	log.Printf(errMsg)
	go etl.ErrorRepository.Save(order.BrandCode, order.ShipmentLocationCode, order.WaybillNo, errMsg)
}

// Extract ...
func (etl ReturnToWarehouseAnytimeETL) Extract(ctx context.Context) (interface{}, error) {
	result, err := repositories.ReturnToWarehouseRepository{}.GetUnsyncedReturnToWarehouseAnytimeOrders()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Transform ...
func (etl ReturnToWarehouseAnytimeETL) Transform(ctx context.Context, source interface{}) (interface{}, error) {
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
	for k, v := range items {
		order, err := entities.ReturnToWarehouseOrder{}.Create(v)
		if err != nil {
			etl.saveError(entities.ReturnToWarehouseOrder{
				BrandCode:            strings.Split(k, "-")[0],
				ShipmentLocationCode: strings.Split(k, "-")[1],
				WaybillNo:            strings.Split(k, "-")[2],
			}, "ReturnToWarehouseAnytimeETL.Transform.orders | "+err.Error())
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Load ...
func (etl ReturnToWarehouseAnytimeETL) Load(ctx context.Context, source interface{}) error {
	orders, ok := source.([]entities.ReturnToWarehouseOrder)
	if !ok {
		return errors.New("Convert Failed")
	}

	for _, order := range orders {
		if len(order.Items) == 0 {
			err := fmt.Errorf("运单号为: %v 的出库单没有商品", order.WaybillNo)
			etl.saveError(order, "ReturnToWarehouseAnytimeETL.Load.orders | "+err.Error())
			continue
		}

		shopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(order.ShipmentLocationCode, order.BrandCode)
		if err != nil {
			etl.saveError(order, "ReturnToWarehouseAnytimeETL.Load.GetShopCodeByChiefShopCodeAndBrandCode | "+err.Error())
			continue
		}

		recvSuppNo, err := repositories.RecvSuppRepository{}.CreateReturnToWarehouseAnytimeOrder(order.BrandCode, shopCode, order.WaybillNo, order.OutDate, order.EmpID, order.DeliveryOrderNo)
		if err != nil {
			etl.saveError(order, "ReturnToWarehouseAnytimeETL.Load.CreateReturnToWarehouseAnytimeOrder | "+err.Error())
			continue
		}

		var wg sync.WaitGroup
		wg.Add(len(order.Items))
		for _, v := range order.Items {
			go func(item entities.ReturnToWarehouseOrderItem, wg *sync.WaitGroup) {
				err := repositories.RecvSuppRepository{}.AddReturnToWarehouseOrderItem(order.BrandCode, shopCode, order.OutDate, recvSuppNo, item.SkuCode, item.Qty, order.EmpID)
				if err != nil {
					etl.saveError(order, "ReturnToWarehouseAnytimeETL.Load.AddReturnToWarehouseOrderItem | "+err.Error())
				}
				wg.Done()
			}(v, &wg)
		}
		wg.Wait()

		// 更新状态的时候需要使用主卖场的Code
		err = repositories.ReturnToWarehouseRepository{}.MarkAnytimeWaybillSynced(order.ShipmentLocationCode, order.WaybillNo)
		if err != nil {
			etl.saveError(order, "ReturnToWarehouseAnytimeETL.Load.MarkAnytimeWaybillSynced | "+err.Error())
			continue
		}
		log.Printf("运单号为：%v 的退仓单（卖场：%v，品牌：%v）已经同步完成。", order.WaybillNo, order.ShipmentLocationCode, order.BrandCode)
	}

	return nil
}
