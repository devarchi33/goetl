package repositories

import (
	cslConst "clearance-adapter/domain/csl-constants"
	"clearance-adapter/factory"
	"clearance-adapter/infra"
	"clearance-adapter/models"
	"clearance-adapter/repositories/entities"
	"errors"
	"log"

	"github.com/go-xorm/xorm"
)

// RecvSuppRepository RecvSupp仓库，包括Master和Detail
type RecvSuppRepository struct{}

// GetUnconfirmedDistributionOrdersByDeadline 已到截止日期仍未入库的出库单
func (RecvSuppRepository) GetUnconfirmedDistributionOrdersByDeadline(deadline string) ([]models.RecvSupp, error) {
	details := make([]models.RecvSupp, 0)
	engine := factory.GetCSLEngine()

	distributionCodes := []string{cslConst.TypLogisticsToShop, cslConst.TypFactoryToShop, cslConst.TypAncillaryProduct}
	err := engine.Join("INNER", "RecvSuppMst",
		`RecvSuppMst.RecvSuppNo = RecvSuppDtl.RecvSuppNo 
		AND RecvSuppMst.BrandCode = RecvSuppDtl.BrandCode 
		AND RecvSuppMst.ShopCode = RecvSuppDtl.ShopCode`).
		Where("RecvSuppMst.RecvSuppStatusCode = ? AND RecvSuppMst.RecvChk = 0 AND RecvSuppMst.DelChk = 0 AND RecvSuppMst.BrandSuppRecvDate <= ?",
			cslConst.StsSentOut, deadline).
		In("RecvSuppMst.ShippingTypeCode", distributionCodes).
		Find(&details)

	if err != nil {
		return nil, err
	}

	return details, nil
}

// GetUnconfirmedTransferOrdersByDeadline 已到截止日期仍未入库的调货出库单
func (RecvSuppRepository) GetUnconfirmedTransferOrdersByDeadline(deadline string) ([]models.RecvSupp, error) {
	details := make([]models.RecvSupp, 0)
	engine := factory.GetCSLEngine()

	err := engine.Join("INNER", "RecvSuppMst",
		`RecvSuppMst.RecvSuppNo = RecvSuppDtl.RecvSuppNo 
		AND RecvSuppMst.BrandCode = RecvSuppDtl.BrandCode 
		AND RecvSuppMst.ShopCode = RecvSuppDtl.ShopCode`).
		Where(`RecvSuppMst.RecvSuppStatusCode = ? 
			AND RecvSuppMst.RecvChk = 0 
			AND RecvSuppMst.DelChk = 0 
			AND RecvSuppMst.ShopSuppRecvDate <= ?
			AND RecvSuppMst.ShippingTypeCode = ?`,
			cslConst.StsSentOut, deadline, cslConst.TypShopToShop).
		Find(&details)

	if err != nil {
		return nil, err
	}

	return details, nil
}

// PutInStorage 入库
func (RecvSuppRepository) PutInStorage(brandCode, shopCode, waybillNo, inDate, empID string) error {
	sql := `
		EXEC [up_CSLK_IOM_UpdateStockInEnterConfirmSave_RecvSuppMst_R1_Clearance_By_WaybillNo]
				@BrandCode = ?,
				@ShopCode = ?,
				@WaybillNo = ?,
				@InDate = ?,
				@EmpID = ?
		`
	_, err := factory.GetCSLEngine().Query(sql, brandCode, shopCode, waybillNo, inDate, empID)
	if err != nil {
		return errors.New("PutInStorage error: " + err.Error())
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
		return errors.New("WriteDownStockMiss error: " + err.Error())
	}

	return nil
}

// GetShopCodeByChiefShopCodeAndBrandCode 根据主卖场Code和子品牌获取子卖场的Code
func (RecvSuppRepository) GetShopCodeByChiefShopCodeAndBrandCode(chiefShopCode, brandCode string) (string, error) {
	sql := `
		SELECT ShopCode
		FROM ComplexShopMapping
		WHERE BrandCode = ?
		AND ChiefShopCode = ?
		AND DelChk = 0

		UNION

		SELECT ShopCode
		FROM Shop
		WHERE BrandCode = ?
		AND  ShopCode = ?
	`

	result, err := factory.GetCSLEngine().Query(sql, brandCode, chiefShopCode, brandCode, chiefShopCode)
	if err != nil {
		return "", err
	}

	if result == nil || len(result) == 0 {
		log.Printf("brandCode: %v", brandCode)
		log.Printf("chiefShopCode: %v", chiefShopCode)
		return "", errors.New("result is nil, get sub shop code failed")
	}

	shop := infra.ConvertByteResult(result)[0]["ShopCode"]

	return shop, nil
}

// GetChiefShopCodeByShopCodeAndBrandCode 根据子卖场Code和子品牌获取主卖场的Code
func (RecvSuppRepository) GetChiefShopCodeByShopCodeAndBrandCode(shopCode, brandCode string) (string, error) {
	sql := `
		SELECT ChiefShopCode
			FROM ComplexShopMapping
			WHERE BrandCode = ?
			AND ShopCode = ?
			AND DelChk = 0
	`

	result, err := factory.GetCSLEngine().Query(sql, brandCode, shopCode)
	if err != nil {
		return "", err
	}

	// ComplexShopMapping 表中没有结果，说明不是复合卖场
	if result == nil || len(result) == 0 {
		return shopCode, nil
	}

	shop := infra.ConvertByteResult(result)[0]["ChiefShopCode"]

	return shop, nil
}

// CreateReturnToWarehouseOrder 创建退仓订单，返回RecvSuppNo
func (RecvSuppRepository) CreateReturnToWarehouseOrder(brandCode, shopCode, waybillNo, outDate, empID, deliveryOrderNo string) (string, error) {
	sql := `
		EXEC [dbo].[up_CSLK_IOM_InsertReturnGoodReservation_RecvSuppMst_C1_Clearance]
			@BrandCode				= ?
			,@ShopCode				= ?
			,@OutDate				= ?
			,@WayBillNo				= ?
			,@ShippingTypeCode		= '41'
			,@ShippingCompanyCode  	= 'SR'
			,@EmpID 				= ?
			,@DeliveryID 			= ?
			,@DeliveryOrderNo 		= ?
	`

	result, err := factory.GetCSLEngine().Query(sql, brandCode, shopCode, outDate, waybillNo, empID, waybillNo, deliveryOrderNo)
	if err != nil {
		log.Printf("brandCode: %v", brandCode)
		log.Printf("shopCode: %v", shopCode)
		log.Printf("outDate: %v", outDate)
		log.Printf("waybillNo: %v", waybillNo)
		return "", errors.New("CreateReturnToWarehouseOrder error: " + err.Error())
	}

	master := infra.ConvertByteResult(result)
	if len(master) == 0 {
		return "", errors.New("CreateReturnToWarehouseOrder error(master is nil): " + err.Error())
	}
	recvSuppNo := master[0]["RecvSuppNo"]

	return recvSuppNo, nil
}

// AddReturnToWarehouseOrderItem 向退仓单中添加商品
func (RecvSuppRepository) AddReturnToWarehouseOrderItem(brandCode, shopCode, outDate, recvSuppNo, skuCode string, qty int, empID string) error {
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

	_, err := factory.GetCSLEngine().Exec(sql, recvSuppNo, brandCode, shopCode, outDate, skuCode, qty, empID)
	if err != nil {
		log.Println("up_CSLK_IOM_InsertReturnGoodReservation_RecvSuppDtl_C1_Clearance params:")
		log.Printf("recvSuppNo: %v", recvSuppNo)
		log.Printf("brandCode: %v", brandCode)
		log.Printf("shopCode: %v", shopCode)
		log.Printf("outDate: %v", outDate)
		log.Printf("skuCode: %v", skuCode)

		return errors.New("AddReturnToWarehouseOrderItem error: " + err.Error())
	}

	return nil
}

// CreateTransferOrder 创建调货出库单
func (RecvSuppRepository) CreateTransferOrder(brandCode, shipmentLocationCode, receiptLocationCode, outDate, waybillNo, boxNo, shippingCompanyCode, deliveryOrderNo, empID string) (string, error) {
	session := factory.GetCSLEngine().NewSession()
	return RecvSuppRepository{}.createTransferOrder(session, brandCode, shipmentLocationCode, receiptLocationCode, outDate, waybillNo, boxNo, shippingCompanyCode, deliveryOrderNo, empID)
}

func (RecvSuppRepository) createTransferOrder(session *xorm.Session, brandCode, shipmentLocationCode, receiptLocationCode, outDate, waybillNo, boxNo, shippingCompanyCode, deliveryOrderNo, empID string) (string, error) {
	sql := `
		EXEC [dbo].[up_CSLK_IOM_InsertRotationOuterReg_RecvSuppMst_C1_Clearance]
			@BrandCode				= ?
			,@ShopCode				= ?
			,@TargetShopCode		= ?
			,@OutDate				= ?
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

	result, err := factory.GetCSLEngine().Query(sql, brandCode, shipmentLocationCode, receiptLocationCode, outDate, waybillNo, boxNo, shippingCompanyCode, deliveryOrderNo, empID)
	if err != nil {
		return "", errors.New("createTransferOrder error: " + err.Error())
	}

	master := infra.ConvertByteResult(result)
	if len(master) == 0 {
		return "", errors.New("createTransferOrder error(master is nil): " + err.Error())
	}
	recvSuppNo := master[0]["RecvSuppNo"]

	return recvSuppNo, nil
}

// AddTransferOrderItem 向调货出库单中添加商品
func (RecvSuppRepository) AddTransferOrderItem(brandCode, shopCode, outDate, recvSuppNo, skuCode string, qty int, empID string) error {
	session := factory.GetCSLEngine().NewSession()
	return RecvSuppRepository{}.addTransferOrderItem(session, brandCode, shopCode, outDate, recvSuppNo, skuCode, qty, empID)
}

func (RecvSuppRepository) addTransferOrderItem(session *xorm.Session, brandCode, shopCode, outDate, recvSuppNo, skuCode string, qty int, empID string) error {
	sql := `
		EXEC [dbo].[up_CSLK_IOM_InsertRotationOuterReg_RecvSuppDtl_C1_Clearance]
			@RecvSuppNo   	= ?
			,@RecvSuppSeqNo	= NULL
			,@BrandCode  	= ?
			,@ShopCode  	= ?
			,@OutDate  		= ?
			,@ProdCode    	= ?
			,@RecvSuppQty 	=  ?
			,@EmpID 		= ?
	`

	_, err := factory.GetCSLEngine().Exec(sql, recvSuppNo, brandCode, shopCode, outDate, skuCode, qty, empID)
	if err != nil {
		log.Println("up_CSLK_IOM_InsertRotationOuterReg_RecvSuppDtl_C1_Clearance params:")
		log.Printf("recvSuppNo: %v", recvSuppNo)
		log.Printf("brandCode: %v", brandCode)
		log.Printf("shopCode: %v", shopCode)
		log.Printf("outDate: %v", outDate)
		log.Printf("skuCode: %v", skuCode)

		return errors.New("addTransferOrderItem error: " + err.Error())
	}

	return nil
}

// ConfirmTransferOrder 调进确认
func (RecvSuppRepository) ConfirmTransferOrder(brandCode, receiptLocationCode, shipmentLocationCode, inDate, waybillNo, boxNo, roundRecvSuppNo, empID string) (string, error) {
	session := factory.GetCSLEngine().NewSession()
	return RecvSuppRepository{}.confirmTransferOrder(session, brandCode, receiptLocationCode, shipmentLocationCode, inDate, waybillNo, boxNo, roundRecvSuppNo, empID)
}

func (RecvSuppRepository) confirmTransferOrder(session *xorm.Session, brandCode, receiptLocationCode, shipmentLocationCode, inDate, waybillNo, boxNo, roundRecvSuppNo, empID string) (string, error) {
	sql := `
		EXEC [dbo].[up_CSLK_IOM_InsertRotationEnterConfirm_RecvSuppMst_C1_Clearance]
			@BrandCode			= ?
			,@ShopCode			= ?
			,@TargetShopCode	= ?
			,@InDate			= ?
			,@WayBillNo			= ?
			,@BoxNo				= ?
			,@RoundRecvSuppNo  	= ?
			,@EmpID 			= ?
	`

	result, err := factory.GetCSLEngine().Query(sql, brandCode, receiptLocationCode, shipmentLocationCode, inDate, waybillNo, boxNo, roundRecvSuppNo, empID)
	if err != nil {
		return "", errors.New("confirmTransferOrder error: " + err.Error())
	}

	master := infra.ConvertByteResult(result)
	if len(master) == 0 {
		return "", errors.New("confirmTransferOrder error(master is nil): " + err.Error())
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

	outRecvSuppNo, err := RecvSuppRepository{}.createTransferOrder(session, order.BrandCode, order.ShipmentShopCode, order.ReceiptShopCode, order.OutDate, order.WaybillNo,
		order.BoxNo, order.ShippingCompanyCode, order.DeliveryOrderNo, order.OutEmpID)
	if err != nil {
		session.Rollback()
		log.Println(err.Error())
		return err
	}

	for _, item := range order.Items {
		err = RecvSuppRepository{}.addTransferOrderItem(session, order.BrandCode, order.ShipmentShopCode, order.OutDate, outRecvSuppNo, item.SkuCode, item.Qty, order.OutEmpID)
		if err != nil {
			session.Rollback()
			log.Println(err.Error())
			return err
		}
	}

	inRecvSuppNo, err := RecvSuppRepository{}.confirmTransferOrder(session, order.BrandCode, order.ReceiptShopCode, order.ShipmentShopCode, order.InDate, order.WaybillNo, order.BoxNo, outRecvSuppNo, order.InEmpID)
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
