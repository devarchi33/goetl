package repositories

import (
	clrConst "clearance-adapter/domain/clr-constants"
	"clearance-adapter/factory"
	"clearance-adapter/models"
	"log"
)

// StockRoundErrorRepository Clearance 记录调货时错误信息的仓库
type StockRoundErrorRepository struct{}

// Save ...
func (StockRoundErrorRepository) Save(brandCode, shipmentLocationCode, receiptLocationCode, waybillNo, errMsg string, errType clrConst.ClrErrorType) error {
	data := make(map[string]interface{})
	param := make(map[string]string)
	param["brand_code"] = brandCode
	param["shipment_location_code"] = shipmentLocationCode
	param["receipt_location_code"] = receiptLocationCode
	param["waybill_no"] = waybillNo
	param["type"] = errType
	param["error_message"] = errMsg
	data["service"] = "clearance-adapter"
	data["param"] = param
	_, err := CreateErrData(data)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	return nil
}

// GetByWaybillNo 根据品牌，卖场和运单号获取错误信息
func (StockRoundErrorRepository) GetByWaybillNo(brandCode, shipmentLocationCode, receiptLocationCode, waybillNo string) (bool, models.StockRoundError, error) {
	tranError := new(models.StockRoundError)
	has, err := factory.GetClrEngine().
		Where("brand_code=? AND shipment_location_code=? AND receipt_location_code=? AND waybill_no=?", brandCode, shipmentLocationCode, receiptLocationCode, waybillNo).
		Get(tranError)

	if err != nil {
		log.Printf(err.Error())
		return has, *tranError, err
	}

	return has, *tranError, err
}
