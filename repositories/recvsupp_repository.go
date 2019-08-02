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

// // SaveMasters 保存Master
// func (RecvSuppRepository) SaveMasters(masters []models.RecvSuppMst) {
// 	for _, master := range masters {
// 		if (RecvSuppRepository{}.validateRecvSuppMst(master)) {
// 			RecvSuppRepository{}.saveRecvSuppMst(master)
// 		}
// 	}
// }

// func (RecvSuppRepository) validateRecvSuppMst(master models.RecvSuppMst) bool {
// 	sql := `SELECT RecvSuppNo
// 				FROM RecvSuppMst
// 				WHERE RecvSuppNo = ?
// 			`
// 	result, err := factory.GetCSLEngine().Query(sql, master.RecvSuppNo)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return false
// 	}

// 	if len(result) == 0 {
// 		return true
// 	}

// 	return false
// }

// func (RecvSuppRepository) saveRecvSuppMst(master models.RecvSuppMst) error {
// 	_, err := factory.GetCSLEngine().Insert(&master)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return err
// 	}

// 	return nil
// }

// // SaveDetails 保存Details
// func (RecvSuppRepository) SaveDetails(details []models.RecvSuppDtl) {
// 	for _, detail := range details {
// 		if (RecvSuppRepository{}.validateRecvSuppDtl(detail)) {
// 			RecvSuppRepository{}.saveRecvSuppDtl(detail)
// 		}
// 	}
// }

// func (RecvSuppRepository) validateRecvSuppDtl(detail models.RecvSuppDtl) bool {
// 	sql := `SELECT RecvSuppNo
// 				FROM RecvSuppDtl
// 				WHERE RecvSuppNo = ? AND ProdCode = ?
// 			`
// 	result, err := factory.GetCSLEngine().Query(sql, detail.RecvSuppNo, detail.ProdCode)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return false
// 	}
// 	if len(result) == 0 {
// 		return true
// 	}
// 	return false
// }

// func (RecvSuppRepository) saveRecvSuppDtl(detail models.RecvSuppDtl) error {
// 	_, err := factory.GetCSLEngine().Insert(&detail)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return err
// 	}
// 	return nil
// }
