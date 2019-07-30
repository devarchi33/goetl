package repositories

import (
	"clearance-adapter/factory"
	"clearance-adapter/models"
	"log"
)

// RecvSuppRepository RecvSupp仓库，包括Master和Detail
type RecvSuppRepository struct{}

// SaveMasters 保存Master
func (RecvSuppRepository) SaveMasters(masters []models.RecvSuppMst) {
	for _, master := range masters {
		if (RecvSuppRepository{}.validateRecvSuppMst(master)) {
			RecvSuppRepository{}.saveRecvSuppMst(master)
		}
	}
}

func (RecvSuppRepository) validateRecvSuppMst(master models.RecvSuppMst) bool {
	sql := `SELECT RecvSuppNo
				FROM RecvSuppMst
				WHERE RecvSuppNo = ?
			`
	result, err := factory.GetCSLEngine().Query(sql, master.RecvSuppNo)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	if len(result) == 0 {
		return true
	}

	return false
}

func (RecvSuppRepository) saveRecvSuppMst(master models.RecvSuppMst) error {
	_, err := factory.GetCSLEngine().Insert(&master)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// SaveDetails 保存Details
func (RecvSuppRepository) SaveDetails(details []models.RecvSuppDtl) {
	for _, detail := range details {
		if (RecvSuppRepository{}.validateRecvSuppDtl(detail)) {
			RecvSuppRepository{}.saveRecvSuppDtl(detail)
		}
	}
}

func (RecvSuppRepository) validateRecvSuppDtl(detail models.RecvSuppDtl) bool {
	sql := `SELECT RecvSuppNo
				FROM RecvSuppDtl
				WHERE RecvSuppNo = ? AND ProdCode = ?
			`
	result, err := factory.GetCSLEngine().Query(sql, detail.RecvSuppNo, detail.ProdCode)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if len(result) == 0 {
		return true
	}
	return false
}

func (RecvSuppRepository) saveRecvSuppDtl(detail models.RecvSuppDtl) error {
	_, err := factory.GetCSLEngine().Insert(&detail)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
