package repositories

import (
	clrConst "clearance-adapter/domain/clr-constants"
	"clearance-adapter/factory"
	"clearance-adapter/models"
	"fmt"
	"log"
)

// StockRoundErrorRepository Clearance 记录调货时错误信息的仓库
type StockRoundErrorRepository struct{}

// Save ...
func (StockRoundErrorRepository) Save(brandCode, shipmentLocationCode, receiptLocationCode, waybillNo, errMsg string, errType clrConst.ClrErrorType) error {
	has, tranError, err := StockRoundErrorRepository{}.GetByWaybillNo(brandCode, shipmentLocationCode, receiptLocationCode, waybillNo)
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	tranError.BrandCode = brandCode
	tranError.ShipmentLocationCode = shipmentLocationCode
	tranError.ReceiptLocationCode = receiptLocationCode
	tranError.WaybillNo = waybillNo
	tranError.Type = errType
	tranError.ErrorMessage = errMsg
	tranError.IsProcessed = false
	tranError.CreatedBy = "Clearance"

	affected, err := factory.GetClrEngine().Insert(&tranError)
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	if affected == 0 {
		err = fmt.Errorf("StockRoundErrorRepository.Save failed affected is 0")
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
