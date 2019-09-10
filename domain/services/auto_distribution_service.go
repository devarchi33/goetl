package services

import (
	"clearance-adapter/config"
	"clearance-adapter/domain/entities"
	p2bConst "clearance-adapter/domain/p2brand-constants"
	"clearance-adapter/infra"
	"clearance-adapter/models"
	"clearance-adapter/repositories"
	"context"
	"errors"
	"log"
	"time"

	clrConst "clearance-adapter/domain/clr-constants"

	"github.com/pangpanglabs/goetl"
)

// AutoDistributionETL 自动入库入库 CSL ->  p2-brand
// p2-brand -> CSL 部分由DistributionETL完成
type AutoDistributionETL struct {
	ErrorRepository repositories.StockDistributionErrorRepository
}

// New 创建 AutoDistributionETL 对象，从Clearance到CSL
func (AutoDistributionETL) New() *goetl.ETL {
	distETL := AutoDistributionETL{
		ErrorRepository: repositories.StockDistributionErrorRepository{}}

	etl := goetl.New(distETL)
	etl.Before(AutoDistributionETL{}.buildDistributionOrders)

	return etl
}

func (etl AutoDistributionETL) saveError(order entities.DistributionOrder, errMsg string) {
	log.Printf(errMsg)
	go etl.ErrorRepository.Save(order.BrandCode, order.ReceiptLocationCode, order.WaybillNo, errMsg, clrConst.TypAutoStockDistributionError)
}

// Extract 获取14天未入库的出库单
func (etl AutoDistributionETL) Extract(ctx context.Context) (interface{}, error) {
	today := time.Now().UTC()
	day, _ := time.ParseDuration("-24h")
	autoDistributeDeadlineDays := time.Duration(config.GetAutoDistributeDeadlineDays())
	twoWeeksAgo := today.Add(autoDistributeDeadlineDays * day).Format("2006-01-02T15:04:05Z")
	deadline := infra.Parse8BitsDate(twoWeeksAgo, nil)

	result, err := repositories.RecvSuppRepository{}.GetUnconfirmedDistributionOrdersByDeadline(deadline)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 相同运单号合并为一个 DistributionOrder
func (etl AutoDistributionETL) buildDistributionOrders(ctx context.Context, source interface{}) (interface{}, error) {
	data, ok := source.([]models.RecvSupp)
	if !ok {
		return nil, errors.New("AutoDistributionETL.buildDistributionOrders error: convert to distribution order failed")
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

	orders := make([]entities.DistributionOrder, 0)
	for _, v := range items {
		order, err := entities.DistributionOrder{}.CreateByRecvSupp(v)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Transform ...
func (etl AutoDistributionETL) Transform(ctx context.Context, source interface{}) (interface{}, error) {
	orders, ok := source.([]entities.DistributionOrder)
	if !ok {
		return nil, errors.New("AutoDistributionETL.Transform error: convert to distribution order failed")
	}

	return orders, nil
}

// Load ...
func (etl AutoDistributionETL) Load(ctx context.Context, source interface{}) error {
	orders, ok := source.([]entities.DistributionOrder)
	if !ok {
		return errors.New("AutoDistributionETL.Load error: convert to distribution order failed")
	}

	for _, order := range orders {
		shopCode, err := repositories.RecvSuppRepository{}.GetChiefShopCodeByShopCodeAndBrandCode(order.ReceiptLocationCode, order.BrandCode)
		if err != nil {
			etl.saveError(order, "AutoDistributionETL.Load.orders | "+err.Error())
			continue
		}
		order.ReceiptLocationCode = shopCode

		if order.Type == p2bConst.TypFactoryToShop {
			err = repositories.DirectDistributionRepository{}.PutInStorage(order)
			if err != nil {
				etl.saveError(order, "AutoDistributionETL.Load.【工厂直送】PutInStorage | "+err.Error())
				continue
			}

			log.Printf("Clearance将【工厂直送】运单号为：%v 的运单（卖场：%v，品牌：%v）自动入库到P2Brand，需要继续等待Clearance将其同步到CSL。", order.WaybillNo, order.ReceiptLocationCode, order.BrandCode)
		} else {
			err = repositories.StockDistributionRepository{}.PutInStorage(order)
			if err != nil {
				etl.saveError(order, "AutoDistributionETL.Load.【物流分配】PutInStorage | "+err.Error())
				continue
			}

			log.Printf("Clearance将【物流分配】运单号为：%v 的运单（卖场：%v，品牌：%v）自动入库到P2Brand，需要继续等待Clearance将其同步到CSL。", order.WaybillNo, order.ReceiptLocationCode, order.BrandCode)
		}
	}

	return nil
}
