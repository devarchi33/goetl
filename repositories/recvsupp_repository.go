package repositories

import (
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"clearance-adapter/models"
	"clearance-adapter/repositories/entities"
	"errors"
	"log"
	"time"

	"github.com/go-xorm/xorm"
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
		Where("RecvSuppMst.BrandCode = ? AND RecvSuppMst.ShopCode = ? AND RecvSuppMst.WayBillNo = ?",
			brandCode, shopCode, waybillNo).
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

// GetShopCodeByChiefShopCodeAndBrandCode 根据主卖场Code和品牌获取子卖场的Code
func (RecvSuppRepository) GetShopCodeByChiefShopCodeAndBrandCode(chiefShopCode, brandCode string) (string, error) {
	sql := `
		SELECT ShopCode
		FROM ComplexShopMapping
		WHERE BrandCode = ?
		AND ChiefShopCode = ?
		AND DelChk = 0
	`

	result, err := factory.GetCSLEngine().Query(sql, brandCode, chiefShopCode)
	if err != nil {
		return "", err
	}

	shop := infra.ConvertByteResult(result)[0]["ShopCode"]

	return shop, nil
}

// CreateReturnToWarehouseOrder 创建退仓订单，返回RecvSuppNo
func (RecvSuppRepository) CreateReturnToWarehouseOrder(brandCode, shopCode, waybillNo, empID, deliveryOrderNo string) (string, error) {
	sql := `
		EXEC [dbo].[up_CSLK_IOM_InsertReturnGoodReservation_RecvSuppMst_C1_Clearance]
			@BrandCode				= ?
			,@ShopCode				= ?
			,@Dates					= ?
			,@WayBillNo				= ?
			,@ShippingTypeCode		= '41'
			,@ShippingCompanyCode  	= 'SR'
			,@EmpID 				= ?
			,@DeliveryID 			= ?
			,@DeliveryOrderNo 		= ?
	`
	today := time.Now().Format("20060102")
	result, err := factory.GetCSLEngine().Query(sql, brandCode, shopCode, today, waybillNo, empID, waybillNo, deliveryOrderNo)
	if err != nil {
		return "", err
	}

	master := infra.ConvertByteResult(result)
	if len(master) == 0 {
		return "", errors.New("exec up_CSLK_IOM_InsertReturnGoodReservation_RecvSuppMst_C1_Clearance failed")
	}
	recvSuppNo := master[0]["RecvSuppNo"]

	return recvSuppNo, nil
}

// AddReturnToWarehouseOrderItem 向退仓单中添加商品
func (RecvSuppRepository) AddReturnToWarehouseOrderItem(brandCode, shopCode, recvSuppNo, skuCode string, qty int, empID string) error {
	sql := `
		EXEC [dbo].[up_CSLK_IOM_InsertReturnGoodReservation_RecvSuppDtl_C1_Clearance]
			@RecvSuppNo			        = ?
			,@RecvSuppSeqNo		        = NULL
			,@BrandCode					= ?
			,@ShopCode					= ?
			,@Dates						= ?
			,@ProdCode					= ?
			,@RecvSuppQty				= ?
			,@AbnormalProdReasonCode	= NULL
			,@EmpID						= ?
			,@AbnormalChkCode    		= 0
			,@AbnormalSerialNo   		= NULL
	`
	today := time.Now().Format("20060102")
	_, err := factory.GetCSLEngine().Exec(sql, recvSuppNo, brandCode, shopCode, today, skuCode, qty, empID)
	if err != nil {
		log.Println("up_CSLK_IOM_InsertReturnGoodReservation_RecvSuppDtl_C1_Clearance params:")
		log.Printf("recvSuppNo: %v", recvSuppNo)
		log.Printf("brandCode: %v", brandCode)
		log.Printf("shopCode: %v", shopCode)
		log.Printf("today: %v", today)
		log.Printf("skuCode: %v", skuCode)

		return err
	}

	return nil
}

// CreateTransferOrder 创建调货出库单
func (RecvSuppRepository) CreateTransferOrder(brandCode, shipmentLocationCode, receiptLocationCode, waybillNo, boxNo, shippingCompanyCode, deliveryOrderNo, empID string) (string, error) {
	session := factory.GetCSLEngine().NewSession()
	return RecvSuppRepository{}.createTransferOrder(session, brandCode, shipmentLocationCode, receiptLocationCode, waybillNo, boxNo, shippingCompanyCode, deliveryOrderNo, empID)
}

func (RecvSuppRepository) createTransferOrder(session *xorm.Session, brandCode, shipmentLocationCode, receiptLocationCode, waybillNo, boxNo, shippingCompanyCode, deliveryOrderNo, empID string) (string, error) {
	sql := `
		EXEC [dbo].[up_CSLK_IOM_InsertRotationOuterReg_RecvSuppMst_C1_Clearance]
			@BrandCode				= ?
			,@ShopCode				= ?
			,@TargetShopCode		= ?
			,@Dates					= ?
			,@WayBillNo				= ?
			,@BoxNo					= ?
			,@ShippingCompanyCode  	= ?
			,@DeliveryOrderNo 		= ?
			,@EmpID 				= ?
			,@IsBigSize 			= 0
			,@ExpressNo 			= NULL
			,@BoxAmount 			= NULL
			,@StockOutUseAmt 		= NULL
			,@ProvinceCode 			= NULL
			,@CityCode 				= NULL
			,@DistrictCode 			= NULL
			,@Area 					= NULL
			,@ShopManagerName 		= NULL
			,@MobilePhone 			= NULL
	`
	today := time.Now().Format("20060102")
	result, err := factory.GetCSLEngine().Query(sql, brandCode, shipmentLocationCode, receiptLocationCode, today, waybillNo, boxNo, shippingCompanyCode, deliveryOrderNo, empID)
	if err != nil {
		return "", errors.New("exec up_CSLK_IOM_InsertRotationOuterReg_RecvSuppMst_C1_Clearance failed " + err.Error())
	}

	master := infra.ConvertByteResult(result)
	if len(master) == 0 {
		return "", errors.New("exec up_CSLK_IOM_InsertRotationOuterReg_RecvSuppMst_C1_Clearance failed")
	}
	recvSuppNo := master[0]["RecvSuppNo"]

	return recvSuppNo, nil
}

// AddTransferOrderItem 向调货出库单中添加商品
func (RecvSuppRepository) AddTransferOrderItem(brandCode, shopCode, recvSuppNo, skuCode string, qty int, empID string) error {
	session := factory.GetCSLEngine().NewSession()
	return RecvSuppRepository{}.addTransferOrderItem(session, brandCode, shopCode, recvSuppNo, skuCode, qty, empID)
}

func (RecvSuppRepository) addTransferOrderItem(session *xorm.Session, brandCode, shopCode, recvSuppNo, skuCode string, qty int, empID string) error {
	sql := `
		EXEC [dbo].[up_CSLK_IOM_InsertRotationOuterReg_RecvSuppDtl_C1_Clearance]
			@RecvSuppNo   	= ?
			,@RecvSuppSeqNo	= NULL
			,@BrandCode  	= ?
			,@ShopCode  	= ?
			,@Dates  		= ?
			,@ProdCode    	= ?
			,@RecvSuppQty 	=  ?
			,@EmpID 		= ?
	`
	today := time.Now().Format("20060102")
	_, err := factory.GetCSLEngine().Exec(sql, recvSuppNo, brandCode, shopCode, today, skuCode, qty, empID)
	if err != nil {
		log.Println("up_CSLK_IOM_InsertRotationOuterReg_RecvSuppDtl_C1_Clearance params:")
		log.Printf("recvSuppNo: %v", recvSuppNo)
		log.Printf("brandCode: %v", brandCode)
		log.Printf("shopCode: %v", shopCode)
		log.Printf("today: %v", today)
		log.Printf("skuCode: %v", skuCode)

		return err
	}

	return nil
}

// ConfirmTransferOrder 调进确认
func (RecvSuppRepository) ConfirmTransferOrder(brandCode, receiptLocationCode, shipmentLocationCode, waybillNo, boxNo, roundRecvSuppNo, empID string) (string, error) {
	session := factory.GetCSLEngine().NewSession()
	return RecvSuppRepository{}.confirmTransferOrder(session, brandCode, receiptLocationCode, shipmentLocationCode, waybillNo, boxNo, roundRecvSuppNo, empID)
}

func (RecvSuppRepository) confirmTransferOrder(session *xorm.Session, brandCode, receiptLocationCode, shipmentLocationCode, waybillNo, boxNo, roundRecvSuppNo, empID string) (string, error) {
	sql := `
		EXEC [dbo].[up_CSLK_IOM_InsertRotationEnterConfirm_RecvSuppMst_C1_Clearance]
			@BrandCode			= ?
			,@ShopCode			= ?
			,@TargetShopCode	= ?
			,@WayBillNo			= ?
			,@BoxNo				= ?
			,@RoundRecvSuppNo  	= ?
			,@EmpID 			= ?
	`

	result, err := factory.GetCSLEngine().Query(sql, brandCode, receiptLocationCode, shipmentLocationCode, waybillNo, boxNo, roundRecvSuppNo, empID)
	if err != nil {
		return "", errors.New("exec up_CSLK_IOM_InsertRotationEnterConfirm_RecvSuppMst_C1_Clearance failed " + err.Error())
	}

	master := infra.ConvertByteResult(result)
	if len(master) == 0 {
		return "", errors.New("exec up_CSLK_IOM_InsertRotationEnterConfirm_RecvSuppMst_C1_Clearance failed")
	}
	recvSuppNo := master[0]["RecvSuppNo"]

	return recvSuppNo, nil
}

// CreateTransferOrderSet 创建调货出库单 + 调货入库确认
func (RecvSuppRepository) CreateTransferOrderSet(order entities.TransferOrderSet) error {
	session := factory.GetCSLEngine().NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		log.Println(err.Error())
		return err
	}

	outRecvSuppNo, err := RecvSuppRepository{}.createTransferOrder(session, order.BrandCode, order.ShipmentShopCode, order.ReceiptShopCode, order.WaybillNo,
		order.BoxNo, order.ShippingCompanyCode, order.DeliveryOrderNo, order.OutEmpID)
	if err != nil {
		session.Rollback()
		log.Println(err.Error())
		return err
	}

	for _, item := range order.Items {
		err = RecvSuppRepository{}.addTransferOrderItem(session, order.BrandCode, order.ShipmentShopCode, outRecvSuppNo, item.SkuCode, item.Qty, order.OutEmpID)
		if err != nil {
			session.Rollback()
			log.Println(err.Error())
			return err
		}
	}

	inRecvSuppNo, err := RecvSuppRepository{}.confirmTransferOrder(session, order.BrandCode, order.ReceiptShopCode, order.ShipmentShopCode, order.WaybillNo, order.BoxNo, outRecvSuppNo, order.InEmpID)
	if err != nil {
		session.Rollback()
		log.Println(err.Error())
		return err
	}

	if err = session.Commit(); err != nil {
		log.Println(err.Error())
		return err
	}

	log.Printf("运单号为：%v 的调货单同步完成，调出登记的RecvSuppNo为：%v，调入确认的RecvSuppNo为：%v", order.WaybillNo, outRecvSuppNo, inRecvSuppNo)

	return nil
}
