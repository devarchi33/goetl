package services

import (
	"clearance-adapter/domain/entities"
	"clearance-adapter/repositories"
	"context"
	"errors"
	"log"

	"github.com/pangpanglabs/goetl"
)

// TransferETL 退仓 p2-brand -> CSL
type TransferETL struct{}

// New 创建 TransferETL 对象，从Clearance到CSL
func (TransferETL) New() *goetl.ETL {
	transferETL := TransferETL{}
	etl := goetl.New(transferETL)

	return etl
}

// Extract ...
func (etl TransferETL) Extract(ctx context.Context) (interface{}, error) {
	result, err := repositories.StockRoundRepository{}.GetUnsyncedTransferInOrders()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Transform ...
func (etl TransferETL) Transform(ctx context.Context, source interface{}) (interface{}, error) {
	data, ok := source.([]map[string]string)
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
	for _, v := range items {
		order, err := entities.TransferOrder{}.Create(v)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Load ...
func (etl TransferETL) Load(ctx context.Context, source interface{}) error {
	orders, ok := source.([]entities.TransferOrder)
	if !ok {
		return errors.New("Convert Failed")
	}

	for _, order := range orders {
		if len(order.Items) == 0 {
			log.Printf("运单号为: %v 的调货单没有商品", order.WaybillNo)
			continue
		}

		shipmentShopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(order.ShipmentLocationCode, order.BrandCode)
		if err != nil {
			log.Printf(err.Error())
		}

		receiptShopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(order.ReceiptLocationCode, order.BrandCode)
		if err != nil {
			log.Printf(err.Error())
		}

		deliveryOrderNo := ""
		recvSuppNo, err := repositories.RecvSuppRepository{}.CreateTransferOrder(order.BrandCode, shipmentShopCode, receiptShopCode, order.WaybillNo, order.BoxNo, order.ShippingCompanyCode, deliveryOrderNo, order.OutEmpID)
		if err != nil {
			log.Printf(err.Error())
		}
		log.Printf("运单号：%v 的出库单RecvSuppNo：%v", order.WaybillNo, recvSuppNo)

		for _, item := range order.Items {
			err := repositories.RecvSuppRepository{}.AddTransferOrderItem(order.BrandCode, shipmentShopCode, recvSuppNo, item.SkuCode, item.Qty, order.OutEmpID)
			if err != nil {
				log.Printf(err.Error())
			}
		}

		confirmRecvSuppNo, err := repositories.RecvSuppRepository{}.ConfirmTransferOrder(order.BrandCode, receiptShopCode, shipmentShopCode, order.WaybillNo, order.BoxNo, recvSuppNo, order.InEmpID)
		if err != nil {
			log.Printf(err.Error())
		}
		log.Printf("运单号：%v 的入库确认RecvSuppNo：%v", order.WaybillNo, confirmRecvSuppNo)

		// 更新状态的时候需要使用主卖场的Code
		err = repositories.StockRoundRepository{}.MarkWaybillSynced(order.ShipmentLocationCode, order.ReceiptLocationCode, order.WaybillNo)
		if err != nil {
			log.Printf(err.Error())
		}
		log.Printf("运单号为：%v 的调货单（出库卖场：%v，入库卖场：%v，品牌：%v）已经同步完成。", order.WaybillNo, order.ShipmentLocationCode, order.ReceiptLocationCode, order.BrandCode)
	}

	return nil
}
