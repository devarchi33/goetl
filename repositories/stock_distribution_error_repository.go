package repositories

import (
	clrConst "clearance-adapter/domain/clr-constants"
	"clearance-adapter/errorlog"
	"clearance-adapter/factory"
	"clearance-adapter/models"
	"log"
)

// StockDistributionErrorRepository Clearance 记录物流分配入库时错误信息的仓库
type StockDistributionErrorRepository struct{}

// Save ...
func (StockDistributionErrorRepository) Save(logID int64, brandCode, receiptLocationCode, waybillNo, errMsg string, errType clrConst.ClrErrorType) error {
	param := make(map[string]string)
	param["brand_code"] = brandCode
	param["receipt_location_code"] = receiptLocationCode
	param["waybill_no"] = waybillNo
	param["type"] = errType
	param["error_message"] = errMsg
	errorlog.ErrorLog{}.AppendLogs(logID, param)
	return nil
}

// GetByWaybillNo 根据品牌，卖场和运单号获取错误信息
func (StockDistributionErrorRepository) GetByWaybillNo(brandCode, receiptLocationCode, waybillNo string) (bool, models.StockDistributionError, error) {
	distError := new(models.StockDistributionError)
	has, err := factory.GetClrEngine().
		Where("brand_code=? AND receipt_location_code=? AND waybill_no=?", brandCode, receiptLocationCode, waybillNo).
		Get(distError)

	if err != nil {
		log.Printf(err.Error())
		return has, *distError, err
	}

	return has, *distError, err
}
