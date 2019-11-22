package services

import (
	clrConst "clearance-adapter/domain/clr-constants"
	cslConst "clearance-adapter/domain/csl-constants"
	"clearance-adapter/domain/entities"
	p2bConst "clearance-adapter/domain/p2brand-constants"
	"clearance-adapter/errorlog"
	"clearance-adapter/repositories"
	"context"
	"errors"
	"log"
	"strings"

	"github.com/pangpanglabs/goetl"
)

// DistributionETL 入库 p2-brand -> CSL
type DistributionETL struct {
	ErrorRepository repositories.StockDistributionErrorRepository
	ErrLogID        int64
}

// New 创建DistributionETL对象，从Clearance到CSL
func (DistributionETL) New() *goetl.ETL {
	logID, _ := errorlog.ErrorLog{}.CreateLog(clrConst.TypStockDistributionError)
	distributionETL := DistributionETL{
		ErrorRepository: repositories.StockDistributionErrorRepository{},
		ErrLogID:        logID,
	}
	etl := goetl.New(distributionETL)
	etl.Before(DistributionETL{}.buildDistributionOrders)
	etl.Before(DistributionETL{}.filterStorableDistributions)
	return etl
}
func (etl DistributionETL) saveError(order entities.DistributionOrder, errMsg string) {
	log.Printf(errMsg)
	etl.ErrorRepository.Save(etl.ErrLogID, order.BrandCode, order.ReceiptLocationCode, order.WaybillNo, errMsg, clrConst.TypStockDistributionError)
}

// Extract ...
func (etl DistributionETL) Extract(ctx context.Context) (interface{}, error) {
	result, err := repositories.StockDistributionRepository{}.GetUnsyncedDistributionOrders()
	if err != nil {
		return nil, err
	}
	resultAuto, err := repositories.StockDistributionRepository{}.GetUnsyncedAutoDistributionOrders()
	if err != nil {
		return nil, err
	}

	directDistResult, err := repositories.DirectDistributionRepository{}.GetUnsyncedDistributionOrders()
	if err != nil {
		return nil, err
	}

	directDistResultAuto, err := repositories.DirectDistributionRepository{}.GetUnsyncedAutoDistributionOrders()
	if err != nil {
		return nil, err
	}

	result = append(result, resultAuto...)
	result = append(result, directDistResult...)
	result = append(result, directDistResultAuto...)
	return result, nil
}

// 相同运单号合并为一个 DistributionOrder
func (etl DistributionETL) buildDistributionOrders(ctx context.Context, source interface{}) (interface{}, error) {
	data, ok := source.([]map[string]string)
	if !ok {
		return nil, errors.New("Convert Failed")
	}

	items := make(map[string][]map[string]string, 0)
	for _, item := range data {
		key := item["brand_code"] + "-" + item["receipt_location_code"] + "-" + item["waybill_no"]
		if _, ok := items[key]; ok {
			items[key] = append(items[key], item)
		} else {
			items[key] = make([]map[string]string, 0)
			items[key] = append(items[key], item)
		}
	}

	orders := make([]entities.DistributionOrder, 0)
	for k, v := range items {
		order, err := entities.DistributionOrder{}.Create(v)
		if err != nil {

			etl.saveError(entities.DistributionOrder{
				BrandCode:           strings.Split(k, "-")[0],
				ReceiptLocationCode: strings.Split(k, "-")[1],
				WaybillNo:           strings.Split(k, "-")[2],
			}, err.Error())

			continue
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// validateDistributions 验证distribution是否可以入库
func (etl DistributionETL) validateDistribution(distribution entities.DistributionOrder) (bool, error) {
	shopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(distribution.ReceiptLocationCode, distribution.BrandCode)
	if err != nil {
		return false, err
	}
	recvSupp, err := repositories.RecvSuppRepository{}.GetByWaybillNo(distribution.BrandCode, shopCode, distribution.WaybillNo)
	if err != nil {
		return false, err
	}

	if len(recvSupp) == 0 {
		err = errors.New("there is no outbound order which waybill no is: " + distribution.WaybillNo + ", shop is: " + distribution.BrandCode + "-" + distribution.ReceiptLocationCode)
		return false, err
	}

	ok := false
	for _, v := range recvSupp {
		ok = ok || (v.RecvSuppStatusCode == cslConst.StsSentOut) // 同一个运单号中只有有一条数据的v.RecvSuppStatusCode == "R" 就是ok的
	}

	for _, v := range recvSupp {
		ok = ok || (v.RecvChk == false) // 同一个运单号中只有有一条数据的 v.RecvChk == false 就是ok的
	}

	if !ok {
		err = errors.New("outbound order that waybill no is " + distribution.WaybillNo + " has been put in storage")
		return false, err
	}
	return true, nil
}

// filterStorableDistributions 过滤出可以入库的运单
func (etl DistributionETL) filterStorableDistributions(ctx context.Context, source interface{}) (interface{}, error) {
	orders, ok := source.([]entities.DistributionOrder)
	if !ok {
		return nil, errors.New("Convert Failed")
	}

	storableDistributions := make([]entities.DistributionOrder, 0)
	for _, order := range orders {
		ok, err := DistributionETL{}.validateDistribution(order)
		if err != nil {
			etl.saveError(order, "DistributionETL.filterStorableDistributions.orders | "+err.Error())
			continue
		}
		if ok {
			storableDistributions = append(storableDistributions, order)
		}
	}
	return storableDistributions, nil
}

// Transform ...
func (etl DistributionETL) Transform(ctx context.Context, source interface{}) (interface{}, error) {
	orders, ok := source.([]entities.DistributionOrder)
	if !ok {
		return nil, errors.New("Convert Failed")
	}

	return orders, nil
}

// Load ...
func (etl DistributionETL) Load(ctx context.Context, source interface{}) error {
	orders, ok := source.([]entities.DistributionOrder)
	if !ok {
		return errors.New("Convert Failed")
	}
	for _, order := range orders {
		shopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(order.ReceiptLocationCode, order.BrandCode)
		if err != nil {
			etl.saveError(order, "DistributionETL.Load.GetShopCodeByChiefShopCodeAndBrandCode | "+err.Error())
			continue
		}
		if !order.IsAutoReceipt {
			err = repositories.RecvSuppRepository{}.PutInStorage(order.BrandCode, shopCode, order.WaybillNo, order.InDate, order.InEmpID)
		} else {
			err = repositories.RecvSuppRepository{}.AutoPutInStorage(order.BrandCode, shopCode, order.WaybillNo, order.InDate)
		}
		if err != nil {
			etl.saveError(order, "DistributionETL.Load.PutInStorage | "+err.Error())
			continue
		}

		if order.Type == p2bConst.TypFactoryToShop {

			etl.writeDownStockMiss(order)

			// 更新状态的时候需要使用主卖场的Code
			err = repositories.DirectDistributionRepository{}.MarkWaybillSynced(order.ReceiptLocationCode, order.WaybillNo)
			if err != nil {
				etl.saveError(order, "DistributionETL.Load。DirectDistributionRepository.MarkWaybillSynced | "+err.Error())
				continue
			}
			log.Printf("【工厂直送】运单号为：%v 的运单（卖场：%v，品牌：%v）已经同步完成。", order.WaybillNo, order.ReceiptLocationCode, order.BrandCode)
		} else {

			etl.writeDownStockMiss(order)

			// 更新状态的时候需要使用主卖场的Code
			err = repositories.StockDistributionRepository{}.MarkWaybillSynced(order.ReceiptLocationCode, order.WaybillNo)
			if err != nil {
				etl.saveError(order, "DistributionETL.Load.StockDistributionRepository。MarkWaybillSynced | "+err.Error())
				continue
			}
			log.Printf("【物流分配】运单号为：%v 的运单（卖场：%v，品牌：%v）已经同步完成。", order.WaybillNo, order.ReceiptLocationCode, order.BrandCode)
		}
	}
	errorlog.ErrorLog{}.Finish(etl.ErrLogID)

	return nil
}

// 记录误差
func (etl DistributionETL) writeDownStockMiss(distribution entities.DistributionOrder) error {
	shopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(distribution.ReceiptLocationCode, distribution.BrandCode)
	if err != nil {
		etl.saveError(distribution, "DistributionETL.writeDownStockMiss.GetShopCodeByChiefShopCodeAndBrandCode | "+err.Error())
		return err
	}
	recvSupp, err := repositories.RecvSuppRepository{}.GetByWaybillNo(distribution.BrandCode, shopCode, distribution.WaybillNo)
	if err != nil {
		etl.saveError(distribution, err.Error())
		return err
	}

	type StockMiss struct {
		BrandCode string
		ShopCode  string
		InDate    string
		WaybillNo string
		EmpID     string
		SkuCode   string
		OutQty    int
		InQty     int
	}

	stockMissMap := make(map[string]StockMiss, 0)
	for _, v := range recvSupp {
		key := distribution.BrandCode + "-" + distribution.ReceiptLocationCode + "-" + recvSupp[0].ShopSuppRecvDate + "-" + distribution.WaybillNo + "-" + v.ProdCode
		_, ok := stockMissMap[key]
		if ok {
			stockMiss := stockMissMap[key]
			stockMiss.OutQty += v.RecvSuppQty
			stockMissMap[key] = stockMiss

		} else {
			stockMiss := StockMiss{
				BrandCode: v.RecvSuppMst.BrandCode,
				ShopCode:  v.RecvSuppMst.ShopCode,
				InDate:    v.ShopSuppRecvDate,
				WaybillNo: v.WayBillNo,
				EmpID:     v.RecvEmpID,
				SkuCode:   v.ProdCode,
				OutQty:    v.RecvSuppQty,
				InQty:     0,
			}
			stockMissMap[key] = stockMiss
		}
	}

	for _, v := range distribution.Items {
		key := distribution.BrandCode + "-" + distribution.ReceiptLocationCode + "-" + recvSupp[0].ShopSuppRecvDate + "-" + distribution.WaybillNo + "-" + v.SkuCode
		_, ok := stockMissMap[key]
		if ok {
			stockMiss := stockMissMap[key]
			stockMiss.InQty += v.Qty
			stockMiss.EmpID = distribution.InEmpID
			stockMissMap[key] = stockMiss
		} else {
			shopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(distribution.ReceiptLocationCode, distribution.BrandCode)
			if err != nil {
				etl.saveError(distribution, "DistributionETL.writeDownStockMiss.Items | "+err.Error())
				continue
			}
			stockMiss := StockMiss{
				BrandCode: distribution.BrandCode,
				ShopCode:  shopCode,
				InDate:    distribution.InDate,
				WaybillNo: distribution.WaybillNo,
				EmpID:     distribution.InEmpID,
				SkuCode:   v.SkuCode,
				OutQty:    0,
				InQty:     v.Qty,
			}
			stockMissMap[key] = stockMiss
		}
	}

	if len(stockMissMap) > 0 {
		for _, v := range stockMissMap {
			if v.OutQty != v.InQty {
				log.Printf("运单号：%v, 有误差，参数：%v", v.WaybillNo, v)
				err := repositories.RecvSuppRepository{}.WriteDownStockMiss(v.BrandCode, v.ShopCode, v.InDate, v.WaybillNo, v.SkuCode, v.EmpID, v.OutQty, v.InQty)
				if err != nil {
					etl.saveError(distribution, "DistributionETL.writeDownStockMiss.WriteDownStockMiss | "+err.Error())
				}
			}
		}
	}

	return nil
}
