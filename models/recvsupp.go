package models

import "time"

// RecvSuppMst ...
type RecvSuppMst struct {
	RecvSuppNo          string    `xorm:"char(14) not null pk "`
	BrandCode           string    `xorm:"varchar(4)"`
	ShopCode            string    `xorm:"char(4)"`
	Dates               string    `xorm:"char(8) not null"`
	RecvSuppType        string    `xorm:"char(1)"`
	ShippingTypeCode    string    `xorm:"char(2)"`
	WayBillNo           string    `xorm:"varchar(50) not null"`
	RecvSuppStatusCode  string    `xorm:"char(1) not null"`
	SeqNo               int       `xorm:"int not null"`
	BoxAmount           int       `xorm:"int"`
	SAPDeliveryNo       string    `xorm:"char(10)"`
	SAPDeliveryDate     string    `xorm:"char(8)"`
	NormalProductType   string    `xorm:"char"`
	ShopSuppRecvDate    string    `xorm:"char(8)"`
	TransTypeCode       string    `xorm:"char"`
	RequestNo           string    `xorm:"char(14)"`
	BoxNo               string    `xorm:"char(20)"`
	PlantCode           string    `xorm:"char(4)"`
	RoundRecvSuppNo     string    `xorm:"char(14)"`
	RoundSAPDeliveryNo  string    `xorm:"char(10)"`
	TargetShopCode      string    `xorm:"char(4)"`
	RecvEmpID           string    `xorm:"char(10)"`
	SuppEmpID           string    `xorm:"char(10)"`
	SAPMenuType         string    `xorm:"char"`
	OrderControlNo      string    `xorm:"char(12)"`
	SendFlag            string    `xorm:"char not null default 'R'"`
	InvtBaseDate        string    `xorm:"char(8)"`
	ProvinceCode        string    `xorm:"char(3)"`
	CityCode            string    `xorm:"char(5)"`
	DistrictCode        string    `xorm:"char(8)"`
	BoxType             string    `xorm:"char(2)"`
	ShippingCompanyCode string    `xorm:"char(2)"`
	RecvChk             bool      `xorm:"bit"`
	DelChk              bool      `xorm:"bit not null default 0"`
	ShopDesc            string    `xorm:"nvarchar(400)"`
	BrandDesc           string    `xorm:"nvarchar(400)"`
	RecvEmpName         string    `xorm:"nvarchar(100)"`
	SuppEmpName         string    `xorm:"nvarchar(200)"`
	VolumeType          string    `xorm:"nvarchar(20)"`
	VolumesUnit         string    `xorm:"nvarchar(10)"`
	Area                string    `xorm:"nvarchar(100)"`
	ShopManagerName     string    `xorm:"nvarchar(10)"`
	BrandSuppRecvDate   string    `xorm:"varchar(8)"`
	InUserID            string    `xorm:"varchar(20) not null"`
	ModiUserID          string    `xorm:"varchar(20) not null"`
	SendState           string    `xorm:"varchar(2) DEFAULT '' not null"`
	ExpressNo           string    `xorm:"varchar(13)"`
	DeliveryID          string    `xorm:"varchar(250)"`
	DeliveryOrderNo     string    `xorm:"varchar(250)"`
	VolumesSize         string    `xorm:"varchar(20)"`
	Channel             string    `xorm:"varchar(20)"`
	MobilePhone         string    `xorm:"varchar(25)"`
	ModiDateTime        time.Time `xorm:"datetime"`
	SendDateTime        time.Time `xorm:"datetime"`
	InDateTime          time.Time `xorm:"datetime"`
	DeliverySendTime    time.Time `xorm:"datetime"`
	DeliveryReceiveTime time.Time `xorm:"datetime"`
	StockOutUseAmt      float64   `xorm:"decimal(9,2)"`
	BoxGram             float64   `xorm:"decimal(18,3)"`
	// SendSeqNo           int       `xorm:"bigint not null"` 这个字段是自增的，数据库会自动插入
}

// RecvSuppDtl ...
type RecvSuppDtl struct {
	RecvSuppNo             string    `xorm:"char(14) not null"`
	BrandCode              string    `xorm:"varchar(4)"`
	ShopCode               string    `xorm:"char(4)"`
	Dates                  string    `xorm:"char(8)"`
	ProdCode               string    `xorm:"varchar(18)"`
	RecvSuppQty            int       `xorm:"int"`
	RecvSuppFixedQty       int       `xorm:"int"`
	DelChk                 bool      `xorm:"bit not null default 0"`
	SalePrice              float64   `xorm:"decimal(19,2)"`
	RecvSuppSeqNo          int       `xorm:"int not null"`
	SeqNo                  int       `xorm:"int"`
	RoundRecvSuppDtSeq     int       `xorm:"int"`
	SupGroupCode           string    `xorm:"char(2)"`
	SAPDeliveryNo          string    `xorm:"char(10)"`
	SAPDeliveryItemNo      string    `xorm:"char(10)"`
	RoundRecvSuppNo        string    `xorm:"char(14)"`
	RoundSAPDeliveryNo     string    `xorm:"char(10)"`
	RoundSAPDeliveryItemNo string    `xorm:"char(10)"`
	PriceTypeCode          string    `xorm:"char(2)"`
	SaipType               string    `xorm:"char(2)"`
	AbnormalProdReasonCode string    `xorm:"char(2)"`
	SendFlag               string    `xorm:"char not null default 'R'"`
	AbnormalChkCode        string    `xorm:"char(2)"`
	AbnormalSerialNo       string    `xorm:"varchar(7)"`
	ApplyID                string    `xorm:"varchar(30)"`
	InUserID               string    `xorm:"varchar(20) not null"`
	ModiUserID             string    `xorm:"varchar(20) not null"`
	SendState              string    `xorm:"varchar(2) not null default ''"`
	ModiReason             string    `xorm:"nvarchar(800)"`
	ModiDateTime           time.Time `xorm:"datetime"`
	SendDateTime           time.Time `xorm:"datetime"`
	InDateTime             time.Time `xorm:"datetime"`
	// SendSeqNo           int       `xorm:"bigint not null"` 这个字段是自增的，数据库会自动插入
}

// RecvSupp ...
type RecvSupp struct {
	RecvSuppMst `xorm:"extends"`
	RecvSuppDtl `xorm:"extends"`
}

// TableName ..
func (RecvSupp) TableName() string {
	return "RecvSuppDtl"
}
