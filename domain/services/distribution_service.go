package services

import (
	"clearance-adapter/domain/entities"
	"clearance-adapter/repositories"
	"context"
	"errors"
	"log"
	"time"

	"github.com/pangpanglabs/goetl"
)

// DistributionETL 入库 p2-brand -> CSL
type DistributionETL struct {
	StartDateTime time.Time
	EndDateTime   time.Time
}

// New 创建DistributionETL对象，从Clearance到CSL
func (DistributionETL) New(startDatetime, endDateTime string) *goetl.ETL {
	local, _ := time.LoadLocation("Local")
	start, _ := time.ParseInLocation("2006-01-02 15:04:05", startDatetime, local)
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", endDateTime, local)
	distributionETL := DistributionETL{
		StartDateTime: start,
		EndDateTime:   end,
	}

	etl := goetl.New(distributionETL)
	etl.Before(DistributionETL{}.buildDistributionOrders)
	etl.Before(DistributionETL{}.filterStorableDistributions)

	return etl
}

// Extract ...
func (etl DistributionETL) Extract(ctx context.Context) (interface{}, error) {
	result, err := repositories.StockDistributionRepository{}.GetDistributionOrdersByCreateAt(etl.StartDateTime, etl.EndDateTime)
	if err != nil {
		return nil, err
	}

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
		key := item["brand_code"] + "-" + item["shop_code"] + "-" + item["waybill_no"]
		if _, ok := items[key]; ok {
			items[key] = append(items[key], item)
		} else {
			items[key] = make([]map[string]string, 0)
			items[key] = append(items[key], item)
		}
	}

	orders := make([]entities.DistributionOrder, 0)
	for _, v := range items {
		order, err := entities.DistributionOrder{}.Create(v)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// validateDistributions 验证distribution是否可以入库
func (etl DistributionETL) validateDistribution(distribution entities.DistributionOrder) (bool, error) {
	shopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(distribution.ShopCode, distribution.BrandCode)
	if err != nil {
		log.Printf(err.Error())
	}
	recvSupp, err := repositories.RecvSuppRepository{}.GetByWaybillNo(distribution.BrandCode, shopCode, distribution.WaybillNo)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}

	if len(recvSupp) == 0 {
		return false, errors.New("there is no outbound order which waybill no is: " + distribution.WaybillNo + ", shop is: " + distribution.BrandCode + "-" + distribution.ShopCode)
	}

	ok := true
	for _, v := range recvSupp {
		ok = ok && (v.RecvSuppStatusCode == "R") // v.RecvSuppStatusCode == "R"的是ok的
	}

	for _, v := range recvSupp {
		ok = ok && (v.RecvChk == false) // v.RecvChk == false 的都是ok的
	}

	if !ok {
		return false, errors.New("outbound order that waybill no is " + distribution.WaybillNo + " has been put in storage")
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
			log.Println(err.Error())
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
		shopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(order.ShopCode, order.BrandCode)
		if err != nil {
			log.Printf(err.Error())
		}
		err = repositories.RecvSuppRepository{}.PutInStorage(order.BrandCode, shopCode, order.WaybillNo, order.EmpID)
		if err != nil {
			log.Printf(err.Error())
		}
		etl.writeDownStockMiss(order)
	}

	return nil
}

// 记录误差
func (etl DistributionETL) writeDownStockMiss(distribution entities.DistributionOrder) error {
	shopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(distribution.ShopCode, distribution.BrandCode)
	if err != nil {
		log.Printf(err.Error())
	}
	recvSupp, err := repositories.RecvSuppRepository{}.GetByWaybillNo(distribution.BrandCode, shopCode, distribution.WaybillNo)
	if err != nil {
		log.Println(err.Error())
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
		key := distribution.BrandCode + "-" + distribution.ShopCode + "-" + recvSupp[0].ShopSuppRecvDate + "-" + distribution.WaybillNo + "-" + v.ProdCode
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
		key := distribution.BrandCode + "-" + distribution.ShopCode + "-" + recvSupp[0].ShopSuppRecvDate + "-" + distribution.WaybillNo + "-" + v.SkuCode
		_, ok := stockMissMap[key]
		if ok {
			stockMiss := stockMissMap[key]
			stockMiss.InQty += v.Qty
			stockMiss.EmpID = distribution.EmpID
			stockMissMap[key] = stockMiss
		} else {
			shopCode, err := repositories.RecvSuppRepository{}.GetShopCodeByChiefShopCodeAndBrandCode(distribution.ShopCode, distribution.BrandCode)
			if err != nil {
				log.Printf(err.Error())
			}
			stockMiss := StockMiss{
				BrandCode: distribution.BrandCode,
				ShopCode:  shopCode,
				InDate:    time.Now().Format("20061012"),
				WaybillNo: distribution.WaybillNo,
				EmpID:     distribution.EmpID,
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
				err := repositories.RecvSuppRepository{}.WriteDownStockMiss(v.BrandCode, v.ShopCode, v.InDate, v.WaybillNo, v.SkuCode, v.EmpID, v.OutQty, v.InQty)
				if err != nil {
					log.Printf(err.Error())
				}
			}
		}
	}

	return nil
}