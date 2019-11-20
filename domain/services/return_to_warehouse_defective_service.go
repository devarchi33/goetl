package services

import (
	clrConst "clearance-adapter/domain/clr-constants"
	"clearance-adapter/domain/entities"
	"clearance-adapter/errorlog"
	"clearance-adapter/repositories"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/pangpanglabs/goetl"
)

// ReturnToWarehouseDefectiveETL 次品退仓 p2-brand -> CSL
type ReturnToWarehouseDefectiveETL struct {
	ErrorRepository repositories.ReturnToWarehouseErrorRepository
	ErrLogID        int64
}

// New 创建 ReturnToWarehouseDefectiveETL 对象，从Clearance到CSL
func (ReturnToWarehouseDefectiveETL) New() *goetl.ETL {
	logID, _ := errorlog.ErrorLog{}.CreateLog(clrConst.TypReturnToWarehouseError)
	returnToWarehouseETL := ReturnToWarehouseDefectiveETL{
		ErrorRepository: repositories.ReturnToWarehouseErrorRepository{},
		ErrLogID:        logID}
	etl := goetl.New(returnToWarehouseETL)
	return etl
}
func (etl ReturnToWarehouseDefectiveETL) saveError(order entities.ReturnToWarehouseOrder, errMsg string) {
	log.Printf(errMsg)
	etl.ErrorRepository.Save(etl.ErrLogID, order.BrandCode, order.ShipmentLocationCode, order.WaybillNo, errMsg)
}

// Extract ...
func (etl ReturnToWarehouseDefectiveETL) Extract(ctx context.Context) (interface{}, error) {
	result, err := repositories.ReturnToWarehouseRepository{}.GetUnsyncedReturnToWarehouseDefectiveOrders()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Transform ...
func (etl ReturnToWarehouseDefectiveETL) Transform(ctx context.Context, source interface{}) (interface{}, error) {
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
			}, "ReturnToWarehouseDefectiveETL.Transform.orders | "+err.Error())
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Load ...
func (etl ReturnToWarehouseDefectiveETL) Load(ctx context.Context, source interface{}) error {
	orders, ok := source.([]entities.ReturnToWarehouseOrder)
	if !ok {
		return errors.New("Convert Failed")
	}

	for _, order := range orders {
		if len(order.Items) == 0 {
			err := fmt.Errorf("运单号为: %v 的出库单没有商品", order.WaybillNo)
			etl.saveError(order, "ReturnToWarehouseDefectiveETL.Load.orders | "+err.Error())
			continue
		}

		shopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(order.ShipmentLocationCode, order.BrandCode)
		if err != nil {
			etl.saveError(order, "ReturnToWarehouseDefectiveETL.Load.GetShopCodeByChiefShopCodeAndBrandCode | "+err.Error())
			continue
		}

		recvSuppNo, err := repositories.RecvSuppRepository{}.CreateReturnToWarehouseDefectiveOrder(order.BrandCode, shopCode, order.WaybillNo, order.OutDate, order.EmpID, order.DeliveryOrderNo)
		if err != nil {
			etl.saveError(order, "ReturnToWarehouseDefectiveETL.Load.CreateReturnToWarehouseDefectiveOrder | "+err.Error())
			continue
		}

		var wg sync.WaitGroup
		wg.Add(len(order.Items))
		for _, v := range order.Items {
			go func(item entities.ReturnToWarehouseOrderItem, wg *sync.WaitGroup) {
				err := repositories.RecvSuppRepository{}.AddReturnToWarehouseDefectiveOrderItem(order.BrandCode, shopCode, order.OutDate, recvSuppNo, item.SkuCode, item.DefectiveProdReasonCode, item.Qty, order.EmpID)
				if err != nil {
					etl.saveError(order, "ReturnToWarehouseDefectiveETL.Load.AddReturnToWarehouseDefectiveOrderItem | "+err.Error())
				}
				wg.Done()
			}(v, &wg)
		}
		wg.Wait()

		// 更新状态的时候需要使用主卖场的Code
		err = repositories.ReturnToWarehouseRepository{}.MarkDefectiveWaybillSynced(order.ShipmentLocationCode, order.WaybillNo)
		if err != nil {
			etl.saveError(order, "ReturnToWarehouseDefectiveETL.Load.MarkDefectiveWaybillSynced | "+err.Error())
			continue
		}
		log.Printf("运单号为：%v 的退仓单（卖场：%v，品牌：%v）已经同步完成。", order.WaybillNo, order.ShipmentLocationCode, order.BrandCode)
	}
	errorlog.ErrorLog{}.Finish(etl.ErrLogID)

	return nil
}
