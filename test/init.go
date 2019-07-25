package test

import (
	"clearance-adapter/factory"
	"clearance-adapter/models"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	factory.Init()
	setUpCSLDB()
	setUpClrDB()
}

func setUpCSLDB() {
	filename, err := filepath.Abs("test/data/RecvSuppMst_data.csv")
	if err != nil {
		panic(err.Error())
	}
	headers, data := readRecvSuppMstData(filename)
	masters := buildRecvSuppMsts(headers, data)
	createRecvSuppMstTable()
	initRecvSuppMstData(masters)
}

func buildRecvSuppMsts(headers map[int]string, data [][]string) []models.RecvSuppMst {
	masters := make([]models.RecvSuppMst, 0)
	for _, row := range data {
		master := new(models.RecvSuppMst)
		for colIndex, val := range row {
			header := headers[colIndex]
			obj := reflect.ValueOf(master)
			if obj.Kind() == reflect.Ptr && !obj.Elem().CanSet() {
				continue
			}
			field := obj.Elem().FieldByName(header)
			if !field.IsValid() {
				continue
			}
			if field.Kind() == reflect.String {
				field.SetString(val)
			} else if field.Kind() == reflect.Int {
				v, _ := strconv.ParseInt(val, 10, 64)
				field.SetInt(v)
			} else if field.Kind() == reflect.Bool {
				if val == "1" {
					field.SetBool(true)
				} else {

					field.SetBool(false)
				}
			} else if field.Kind() == reflect.Float64 {
				v, _ := strconv.ParseFloat(val, 64)
				field.SetFloat(v)
			}
		}
		masters = append(masters, *master)
	}

	return masters
}

func createRecvSuppMstTable() {
	sql := `
		DROP TABLE IF EXISTS RecvSuppMst;
		CREATE TABLE RecvSuppMst
		(
			RecvSuppNo CHAR(14) PRIMARY KEY NOT NULL,
			BrandCode VARCHAR(4),
			ShopCode CHAR(4),
			Dates CHAR(8) NOT NULL,
			SeqNo INT NOT NULL,
			SAPDeliveryNo CHAR(10),
			SAPDeliveryDate CHAR(8),
			RecvSuppType CHAR,
			NormalProductType CHAR,
			ShopSuppRecvDate CHAR(8),
			BrandSuppRecvDate VARCHAR(8),
			TransTypeCode CHAR,
			ShippingTypeCode CHAR(2),
			WayBillNo VARCHAR(13) NOT NULL,
			RecvSuppStatusCode CHAR NOT NULL,
			RequestNo CHAR(14),
			BoxNo CHAR(20),
			ShopDesc NVARCHAR(400),
			BrandDesc NVARCHAR(400),
			PlantCode CHAR(4),
			RoundRecvSuppNo CHAR(14),
			RoundSAPDeliveryNo CHAR(10),
			TargetShopCode CHAR(4),
			RecvChk BIT,
			OrderControlNo CHAR(12),
			RecvEmpID CHAR(10),
			RecvEmpName NVARCHAR(100),
			SuppEmpID CHAR(10),
			SuppEmpName NVARCHAR(200),
			SAPMenuType CHAR,
			DelChk BIT DEFAULT 0 NOT NULL,
			InUserID VARCHAR(20) NOT NULL,
			InDateTime DATETIME,
			ModiUserID VARCHAR(20) NOT NULL,
			ModiDateTime DATETIME,
			SendState VARCHAR(2) DEFAULT '' NOT NULL,
			SendFlag CHAR DEFAULT 'R' NOT NULL,
			SendSeqNo BIGINT NOT NULL IDENTITY,
			SendDateTime DATETIME,
			InvtBaseDate CHAR(8),
			BoxAmount INT,
			StockOutUseAmt DECIMAL(9,2),
			ExpressNo VARCHAR(13),
			ShippingCompanyCode CHAR(2),
			BoxGram DECIMAL(18,3),
			DeliveryID VARCHAR(250),
			DeliveryOrderNo VARCHAR(250),
			VolumeType NVARCHAR(20),
			VolumesSize VARCHAR(20),
			VolumesUnit NVARCHAR(10),
			Channel VARCHAR(20),
			ProvinceCode CHAR(3),
			CityCode CHAR(5),
			DistrictCode CHAR(8),
			Area NVARCHAR(100),
			ShopManagerName NVARCHAR(10),
			MobilePhone VARCHAR(25),
			DeliverySendTime DATETIME,
			DeliveryReceiveTime DATETIME,
			BoxType CHAR(2)
		);
		CREATE INDEX CL_RecvSuppMst ON RecvSuppMst (BrandCode, ShopCode, InvtBaseDate);
		CREATE INDEX idxf1_RecvSuppMst ON RecvSuppMst (SendFlag, SendSeqNo);
		CREATE INDEX idxf2_RecvSuppMst ON RecvSuppMst (TargetShopCode, ShopSuppRecvDate, ShippingTypeCode, RecvSuppType, RecvSuppStatusCode);
		CREATE INDEX idxf3_RECVSUPPMST ON RecvSuppMst (BrandCode, TargetShopCode, InvtBaseDate, DelChk, ShippingTypeCode, RecvSuppType, RecvSuppStatusCode, RecvSuppNo);
		CREATE INDEX idx4_RECVSUPPMST ON RecvSuppMst (RoundRecvSuppNo);
		CREATE INDEX idx5_RECVSUPPMST ON RecvSuppMst (BrandCode, ShopCode, ShopSuppRecvDate);
		CREATE INDEX idx6_RecvSuppMst ON RecvSuppMst (WayBillNo);
		CREATE INDEX idx7_RecvSuppMst ON RecvSuppMst (ShopCode, RecvSuppType, WayBillNo, BoxNo, RecvSuppNo);
		CREATE INDEX idx8_RecvSuppMst ON RecvSuppMst (RequestNo, ShopSuppRecvDate);
		CREATE INDEX idx9_RecvSuppMst ON RecvSuppMst (BoxNo, TargetShopCode);
		CREATE INDEX idx_DeliveryID_DeliveryOrderNo ON RecvSuppMst (DeliveryID, DeliveryOrderNo);
		CREATE INDEX idx20_RecvSuppMst_DeliveryOrderNo ON RecvSuppMst (DeliveryOrderNo);
		CREATE INDEX idx21_RecvSuppMst_DeliveryID_DeliveryOrderNo ON RecvSuppMst (DeliveryID, DeliveryOrderNo);
		CREATE INDEX idx_RecvSuppMst_BrandSuppRecvDate ON RecvSuppMst (BrandSuppRecvDate DESC);
	`
	if _, err := factory.GetCSLEngine().Exec(sql); err != nil {
		fmt.Printf("createRecvSuppMstTable error: %v", err.Error())
		fmt.Println()
	}
}

func readRecvSuppMstData(filename string) (map[int]string, [][]string) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}
	reader := csv.NewReader(strings.NewReader(string(bytes)))
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err.Error())
	}

	headerRow := records[0]
	headers := make(map[int]string, 0)
	for i, header := range headerRow {
		headers[i] = header
	}

	return headers, records[1:len(records)]
}

func initRecvSuppMstData(masters []models.RecvSuppMst) {
	for _, master := range masters {
		if affected, err := factory.GetCSLEngine().Insert(&master); err != nil {
			fmt.Printf("initRecvSuppMstData error: %v", err.Error())
			fmt.Println()
			fmt.Printf("affected: %v", affected)
			fmt.Println()
		}
	}
}

func setUpClrDB() {
	createTransactionsTable()
}

func createTransactionsTable() {
	sql := `
		DROP TABLE IF EXISTS transactions;
	`
	if _, err := factory.GetClrEngine().Exec(sql); err != nil {
		fmt.Printf("createTransactionsTable error: %v", err.Error())
		fmt.Println()
	}

	sql = `
		CREATE TABLE transactions
		(
			id BIGINT(20) PRIMARY KEY NOT NULL AUTO_INCREMENT,
			transaction_id VARCHAR(14) NOT NULL,
			waybill_no VARCHAR(13) NOT NULL,
			box_no VARCHAR(20) NOT NULL,
			sku_code VARCHAR(18) NOT NULL,
			qty INT NOT NULL
		);
	`
	if _, err := factory.GetClrEngine().Exec(sql); err != nil {
		fmt.Printf("createTransactionsTable error: %v", err.Error())
		fmt.Println()
	}
}