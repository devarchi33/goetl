package repositories

import (
	clrConst "clearance-adapter/domain/clr-constants"
	"clearance-adapter/factory"
	"clearance-adapter/models"
	"fmt"
	"log"
)

// StockDistributionErrorRepository Clearance 记录物流分配入库时错误信息的仓库
type StockDistributionErrorRepository struct{}

// Save ...
func (StockDistributionErrorRepository) Save(brandCode, receiptLocationCode, waybillNo, errMsg string, errType clrConst.ClrErrorType) error {
	has, distError, err := StockDistributionErrorRepository{}.GetByWaybillNo(brandCode, receiptLocationCode, waybillNo)
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	distError.BrandCode = brandCode
	distError.ReceiptLocationCode = receiptLocationCode
	distError.WaybillNo = waybillNo
	distError.Type = errType
	distError.ErrorMessage = errMsg
	distError.IsProcessed = false
	distError.CreatedBy = "Clearance"

	affected, err := factory.GetClrEngine().Insert(&distError)
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	if affected == 0 {
		err = fmt.Errorf("StockDistributionErrorRepository.Save failed affected is 0")
		log.Printf(err.Error())
		return err
	}

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
