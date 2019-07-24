package goetl

import (
	"context"
	"errors"

	"clearance-adapter/factory"
	"clearance-adapter/models"

	"github.com/pangpanglabs/goetl"
)

// ConvertByteResult ...
func ConvertByteResult(source []map[string][]byte) []map[string]string {
	result := make([]map[string]string, 0)
	if source == nil || len(source) == 0 {
		return result
	}
	for _, sourceItem := range source {
		item := make(map[string]string, 0)
		if sourceItem == nil || len(sourceItem) == 0 {
			continue
		}
		for key, value := range sourceItem {
			item[key] = string(value)
		}
		result = append(result, item)
	}
	return result
}

func buildETL() *goetl.ETL {
	etl := goetl.New(AssignETL{})
	return etl
}

// AssignETL 分配
type AssignETL struct{}

// Extract ...
func (etl AssignETL) Extract(ctx context.Context) (interface{}, error) {
	engine := factory.GetCSLEngine()
	masters := make([]models.RecvSuppMst, 0)
	if err := engine.Where("BrandCode = ? AND ShopCode = ?", "EE", "CDVQ").Find(&masters); err != nil {
		return nil, err
	}

	return masters, nil
}

// Transform ...
func (etl AssignETL) Transform(ctx context.Context, source interface{}) (interface{}, error) {
	cslMasters, ok := source.([]models.RecvSuppMst)
	if !ok {
		return nil, errors.New("Convert Failed")
	}
	mslMasters := make([]models.TransactionMaster, 0)
	for _, mst := range cslMasters {
		mslMasters = append(mslMasters, models.TransactionMaster{
			Date:            mst.Dates,
			PlantCode:       mst.BrandCode + "-" + mst.ShopCode,
			WaybillNo:       mst.WayBillNo,
			OrderNo:         mst.RecvSuppNo,
			TransactionCode: "OS100",
			Channel:         "CLEARANCE",
		})
	}

	return mslMasters, nil
}

// ReadyToLoad ...
func (etl AssignETL) ReadyToLoad(ctx context.Context, source interface{}) error {
	masters, ok := source.([]models.TransactionMaster)
	if !ok {
		return errors.New("Convert Failed")
	}
	savedMasters := make([]models.TransactionMaster, 0)
	for _, mst := range masters {
		sql := `SELECT id
		FROM transaction_masters
		WHERE order_no = ?
		`

		engine := factory.GetINVEngine()
		result, err := engine.Query(sql, mst.OrderNo)
		if err != nil {
			return err
		}

		if len(result) == 0 {
			savedMasters = append(savedMasters, mst)
		}
	}

	return nil
}

// Load ...
func (etl AssignETL) Load(ctx context.Context, source interface{}) error {
	if source == nil {
		return errors.New("source is nil")
	}
	mslMasters, ok := source.([]models.TransactionMaster)
	if !ok {
		return errors.New("Convert Failed")
	}
	engine := factory.GetINVEngine()

	if _, err := engine.Insert(&mslMasters); err != nil {
		return err
	}
	return nil
}
