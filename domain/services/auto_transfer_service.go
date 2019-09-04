package services

import (
	"clearance-adapter/domain/entities"
	"clearance-adapter/infra"
	"clearance-adapter/models"
	"clearance-adapter/repositories"
	"context"
	"errors"
	"log"
	"time"

	"github.com/pangpanglabs/goetl"
)

// AutoTransferETL 自动调货入库 CSL ->  p2-brand
// p2-brand -> CSL 部分由TransferETL完成
type AutoTransferETL struct{}

// New 创建 AutoTransferETL 对象，从Clearance到CSL
func (AutoTransferETL) New() *goetl.ETL {
	transferETL := AutoTransferETL{}

	etl := goetl.New(transferETL)
	etl.Before(AutoTransferETL{}.buildTransferOrders)

	return etl
}

// Extract 获取14天未入库的出库单
func (etl AutoTransferETL) Extract(ctx context.Context) (interface{}, error) {
	today := time.Now().UTC()
	day, _ := time.ParseDuration("-24h")
	twoWeeksAgo := today.Add(14 * day).Format("2006-01-02T15:04:05Z")
	deadline := infra.Parse8BitsDate(twoWeeksAgo, nil)

	result, err := repositories.RecvSuppRepository{}.GetUnconfirmedTransferOrdersByDeadline(deadline)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 相同运单号合并为一个 TransferOrder
func (etl AutoTransferETL) buildTransferOrders(ctx context.Context, source interface{}) (interface{}, error) {
	data, ok := source.([]models.RecvSupp)
	if !ok {
		return nil, errors.New("AutoTransferETL.buildTransferOrders error: convert to distribution order failed")
	}

	items := make(map[string][]models.RecvSupp, 0)
	for _, item := range data {
		key := item.RecvSuppMst.BrandCode + "-" + item.RecvSuppMst.ShopCode + "-" + item.WayBillNo
		if _, ok := items[key]; ok {
			items[key] = append(items[key], item)
		} else {
			items[key] = make([]models.RecvSupp, 0)
			items[key] = append(items[key], item)
		}
	}

	orders := make([]entities.TransferOrder, 0)
	for _, v := range items {
		order, err := entities.TransferOrder{}.CreateByRecvSupp(v)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Transform ...
func (etl AutoTransferETL) Transform(ctx context.Context, source interface{}) (interface{}, error) {
	orders, ok := source.([]entities.TransferOrder)
	if !ok {
		return nil, errors.New("AutoTransferETL.Transform error: convert to distribution order failed")
	}

	return orders, nil
}

// Load ...
func (etl AutoTransferETL) Load(ctx context.Context, source interface{}) error {
	orders, ok := source.([]entities.TransferOrder)
	if !ok {
		return errors.New("AutoTransferETL.Load error: convert to distribution order failed")
	}

	for _, order := range orders {
		chiefShipmentLocationCode, err := repositories.RecvSuppRepository{}.GetChiefShopCodeByShopCodeAndBrandCode(order.ShipmentLocationCode, order.BrandCode)
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		chiefReceiptLocationCode, err := repositories.RecvSuppRepository{}.GetChiefShopCodeByShopCodeAndBrandCode(order.ReceiptLocationCode, order.BrandCode)
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		err = repositories.StockRoundRepository{}.TransferIn(order.WaybillNo, chiefShipmentLocationCode, chiefReceiptLocationCode)
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		log.Printf("Clearance将运单号为：%v 的运单（从卖场：%v到卖场%v，品牌：%v）自动调货入库到P2Brand，需要继续等待Clearance将其同步到CSL。", order.WaybillNo, order.ShipmentLocationCode, order.ReceiptLocationCode, order.BrandCode)
	}

	return nil
}
