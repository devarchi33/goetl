package repositories

import (
	"clearance-adapter/factory"
	"clearance-adapter/models"
)

// RecvSuppRepository RecvSupp仓库，包括Master和Detail
type RecvSuppRepository struct{}

// PutInStorage 入库
func (RecvSuppRepository) PutInStorage(brandCode, shopCode, waybillNo, empID string) error {
	sql := `
		EXEC [up_CSLK_IOM_UpdateStockInEnterConfirmSave_RecvSuppMst_R1_Clearance_By_WaybillNo]
				@BrandCode = ?,
				@ShopCode = ?,
				@WaybillNo = ?,
				@EmpID = ?
		`
	_, err := factory.GetCSLEngine().Query(sql, brandCode, shopCode, waybillNo, empID)
	if err != nil {
		return err
	}

	return nil
}

// GetByWaybillNo 根据运单号获取出库单
func (RecvSuppRepository) GetByWaybillNo(brandCode, shopCode, waybillNo string) ([]models.RecvSupp, error) {
	details := make([]models.RecvSupp, 0)
	engine := factory.GetCSLEngine()
	err := engine.Join("INNER", "RecvSuppMst",
		`RecvSuppMst.RecvSuppNo = RecvSuppDtl.RecvSuppNo 
		AND RecvSuppMst.BrandCode = RecvSuppDtl.BrandCode 
		AND RecvSuppMst.ShopCode = RecvSuppDtl.ShopCode`).
		Where("RecvSuppMst.BrandCode = ? AND RecvSuppMst.ShopCode = ? AND RecvSuppMst.WayBillNo = ? AND RecvSuppMst.BoxNo = ?",
			brandCode, shopCode, waybillNo, waybillNo).
		Find(&details)
	if err != nil {
		return nil, err
	}

	return details, nil
}

// WriteDownStockMiss 记录误差
func (RecvSuppRepository) WriteDownStockMiss(brandCode, shopCode, inDate, waybillNo, skuCode, empID string, outQty, inQty int) error {
	sql := `
		EXEC [up_CSLK_IOM_InsertStockInMissSave_StockMisDtl_C1_Clearance]
				@BrandCode = ?,
				@ShopCode = ?,
				@WaybillNo = ?,
				@ProdCode = ?,
				@ShopRecvSuppQty = ?,
				@ShopInFixQty = ?,
				@ErrorRegEmpID = ?,
				@RecvDate = ?
		`

	_, err := factory.GetCSLEngine().Query(sql, brandCode, shopCode, waybillNo, skuCode, outQty, inQty, empID, inDate)
	if err != nil {
		return err
	}

	return nil
}
