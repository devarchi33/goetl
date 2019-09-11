package services

import (
	clrConst "clearance-adapter/domain/clr-constants"
	"clearance-adapter/domain/entities"
	"clearance-adapter/repositories"
	repoEntities "clearance-adapter/repositories/entities"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/pangpanglabs/goetl"
)

// TransferETL 退仓 p2-brand -> CSL
// 调入之后才会同步到CSL（调出&调入），如果只有调出单，那么是不会同步到CSL的
type TransferETL struct {
	ErrorRepository repositories.StockRoundErrorRepository
}

// New 创建 TransferETL 对象，从Clearance到CSL
func (TransferETL) New() *goetl.ETL {
	transferETL := TransferETL{
		ErrorRepository: repositories.StockRoundErrorRepository{},
	}
	etl := goetl.New(transferETL)

	return etl
}

func (etl TransferETL) saveError(order entities.TransferOrder, errMsg string) {
	log.Printf(errMsg)
	go etl.ErrorRepository.Save(order.BrandCode, order.ShipmentLocationCode, order.ReceiptLocationCode, order.WaybillNo, errMsg, clrConst.TypTransferError)
}

func (etl TransferETL) saveErrorByTransferOrderSet(orderSet repoEntities.TransferOrderSet, errMsg string) {
	etl.saveError(entities.TransferOrder{
		BrandCode:            orderSet.BrandCode,
		ShipmentLocationCode: orderSet.ShipmentLocationCode,
		ReceiptLocationCode:  orderSet.ReceiptLocationCode,
		WaybillNo:            orderSet.WaybillNo,
	}, errMsg)
}

// Extract ...
func (etl TransferETL) Extract(ctx context.Context) (interface{}, error) {
	result, err := repositories.StockRoundRepository{}.GetUnsyncedTransferInOrders()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (etl TransferETL) buildTransferOrders(ctx context.Context, source interface{}) (interface{}, error) {
	data, ok := source.([]map[string]string)

	// fmt.Println(len(data))
	if !ok {
		return nil, errors.New("Convert Failed")
	}

	items := make(map[string][]map[string]string, 0)
	for _, item := range data {
		key := item["brand_code"] + "-" + item["shipment_location_code"] + "-" + item["receipt_location_code"] + "-" + item["waybill_no"]
		if _, ok := items[key]; ok {
			items[key] = append(items[key], item)
		} else {
			items[key] = make([]map[string]string, 0)
			items[key] = append(items[key], item)
		}
	}

	orders := make([]entities.TransferOrder, 0)
	for k, v := range items {
		order, err := entities.TransferOrder{}.Create(v)
		if err != nil {

			etl.saveError(entities.TransferOrder{
				BrandCode:            strings.Split(k, "-")[0],
				ShipmentLocationCode: strings.Split(k, "-")[1],
				ReceiptLocationCode:  strings.Split(k, "-")[2],
				WaybillNo:            strings.Split(k, "-")[3],
			}, err.Error())

			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (etl TransferETL) buildTransferOrderSets(ctx context.Context, source interface{}) (interface{}, error) {
	orders, ok := source.([]entities.TransferOrder)
	if !ok {
		return nil, errors.New("Convert to TransferOrders failed")
	}

	orderSets := make([]repoEntities.TransferOrderSet, 0)
	for _, order := range orders {
		shipmentShopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(order.ShipmentLocationCode, order.BrandCode)
		if err != nil {
			etl.saveError(order, "TransferETL.buildTransferOrderSets.GetShopCodeByChiefShopCodeAndBrandCode(shipmentShopCode) | "+err.Error())
			continue
		}

		receiptShopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(order.ReceiptLocationCode, order.BrandCode)
		if err != nil {
			etl.saveError(order, "TransferETL.buildTransferOrderSets.GetShopCodeByChiefShopCodeAndBrandCode(receiptShopCode) | "+err.Error())
			continue
		}

		orderSet := repoEntities.TransferOrderSet{
			BrandCode:            order.BrandCode,
			ShipmentLocationCode: order.ShipmentLocationCode,
			ShipmentShopCode:     shipmentShopCode,
			ReceiptLocationCode:  order.ReceiptLocationCode,
			ReceiptShopCode:      receiptShopCode,
			WaybillNo:            order.WaybillNo,
			BoxNo:                order.BoxNo,
			ShippingCompanyCode:  order.ShippingCompanyCode,
			DeliveryOrderNo:      "",
			OutDate:              order.OutDate,
			InDate:               order.InDate,
			OutEmpID:             order.OutEmpID,
			InEmpID:              order.InEmpID,
			Items:                make([]repoEntities.TransferOrderSetItem, 0),
		}
		for _, item := range order.Items {
			orderSet.Items = append(orderSet.Items, repoEntities.TransferOrderSetItem{
				SkuCode: item.SkuCode,
				Qty:     item.Qty,
			})
		}
		orderSets = append(orderSets, orderSet)
	}

	return orderSets, nil
}

// Transform ...
func (etl TransferETL) Transform(ctx context.Context, source interface{}) (interface{}, error) {
	orders, err := TransferETL{}.buildTransferOrders(ctx, source)
	if err != nil {
		return nil, errors.New("Build TransferOrders Failed " + err.Error())
	}

	orderSets, err := TransferETL{}.buildTransferOrderSets(ctx, orders)
	if err != nil {
		return nil, errors.New("Build TransferOrderSets Failed " + err.Error())
	}

	return orderSets, nil
}

// Load ...
func (etl TransferETL) Load(ctx context.Context, source interface{}) error {
	orderSets, ok := source.([]repoEntities.TransferOrderSet)
	if !ok {
		return errors.New("Convert to TransferOrderSets failed")
	}

	for _, orderSet := range orderSets {
		if len(orderSet.Items) == 0 {
			err := fmt.Errorf("运单号为: %v 的调货单没有商品", orderSet.WaybillNo)

			etl.saveErrorByTransferOrderSet(orderSet, "TransferETL.Load.orderSets(len(orderSet.Items) == 0) | "+err.Error())

			continue
		}

		err := repositories.RecvSuppRepository{}.CreateTransferOrderSet(orderSet)
		if err != nil {
			etl.saveErrorByTransferOrderSet(orderSet, "TransferETL.Load.orderSets.CreateTransferOrderSet | "+err.Error())

			return err
		}

		// 更新状态的时候需要使用主卖场的Code
		err = repositories.StockRoundRepository{}.MarkWaybillSynced(orderSet.ShipmentLocationCode, orderSet.ReceiptLocationCode, orderSet.WaybillNo)
		if err != nil {
			etl.saveErrorByTransferOrderSet(orderSet, "TransferETL.Load.orderSets.MarkWaybillSynced | "+err.Error())
			return err
		}

		log.Printf("运单号为：%v 的调货单（出库卖场：%v，入库卖场：%v，品牌：%v）已经同步完成。", orderSet.WaybillNo, orderSet.ShipmentLocationCode, orderSet.ReceiptLocationCode, orderSet.BrandCode)
	}

	return nil
}
