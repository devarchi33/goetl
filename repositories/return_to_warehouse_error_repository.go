package repositories

import (
	clrConst "clearance-adapter/domain/clr-constants"
	"clearance-adapter/factory"
	"clearance-adapter/models"
	"log"
)

// ReturnToWarehouseErrorRepository Clearance 记录退仓时错误信息的仓库
type ReturnToWarehouseErrorRepository struct{}

// Save ...
func (ReturnToWarehouseErrorRepository) Save(brandCode, shipmentLocationCode, waybillNo string, errMsg string) error {
	data := make(map[string]interface{})
	param := make(map[string]string)
	param["brand_code"] = brandCode
	param["shipment_location_code"] = shipmentLocationCode
	param["waybill_no"] = waybillNo
	param["type"] = clrConst.TypReturnToWarehouseError
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

// GetByWaybillNo 根据品牌，出库卖场和运单号获取错误信息
func (ReturnToWarehouseErrorRepository) GetByWaybillNo(brandCode, shipmentLocationCode, waybillNo string) (bool, models.ReturnToWarehouseError, error) {
	rtwError := new(models.ReturnToWarehouseError)
	has, err := factory.GetClrEngine().
		Where("brand_code=? AND shipment_location_code=? AND waybill_no=?", brandCode, shipmentLocationCode, waybillNo).
		Get(rtwError)

	if err != nil {
		log.Printf(err.Error())
		return has, *rtwError, err
	}

	return has, *rtwError, err
}
