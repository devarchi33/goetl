package repositories

import (
	clrConst "clearance-adapter/domain/clr-constants"
	"clearance-adapter/factory"
	"clearance-adapter/models"
	"fmt"
	"log"
)

// ReturnToWarehouseErrorRepository Clearance 记录退仓时错误信息的仓库
type ReturnToWarehouseErrorRepository struct{}

// Save ...
func (ReturnToWarehouseErrorRepository) Save(brandCode, shipmentLocationCode, waybillNo string, errMsg string) error {
	has, rtwError, err := ReturnToWarehouseErrorRepository{}.GetByWaybillNo(brandCode, shipmentLocationCode, waybillNo)
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	rtwError.BrandCode = brandCode
	rtwError.ShipmentLocationCode = shipmentLocationCode
	rtwError.WaybillNo = waybillNo
	rtwError.Type = clrConst.TypReturnToWarehouseError
	rtwError.ErrorMessage = errMsg
	rtwError.IsProcessed = false
	rtwError.CreatedBy = "Clearance"

	affected, err := factory.GetClrEngine().Insert(&rtwError)
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	if affected == 0 {
		err = fmt.Errorf("ReturnToWarehouseErrorRepository.Save failed affected is 0")
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
