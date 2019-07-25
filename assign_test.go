package main

import (
	"clearance-adapter/models"
	_ "clearance-adapter/test"
	"context"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

// TransactionID:
// WaybillNo:
// BoxNo:
// SkuCode:
// Qty:
func TestTransform(t *testing.T) {
	Convey("测试AssignETL的Transform方法", t, func() {
		Convey("如果map中的Key都是正确的，应该得到正确的Transaction数组", func() {
			master := models.RecvSuppMst{
				RecvSuppNo: "EECDVQ201907220001",
				WayBillNo:  "SR201907220001",
				BoxNo:      "01",
			}
			detail := models.RecvSuppDtl{
				ProdCode:    "SPJA948S2230095",
				RecvSuppQty: 2,
			}
			recvSupp := models.RecvSupp{
				RecvSuppMst: master,
				RecvSuppDtl: detail,
			}
			transactions, err := AssignETL{}.Transform(context.Background(), []models.RecvSupp{recvSupp})
			So(err, ShouldBeNil)
			transaction := transactions.([]models.Transaction)[0]
			So(transaction.ID, ShouldEqual, 0)
			So(transaction.TransactionID, ShouldEqual, "EECDVQ201907220001")
			So(transaction.WaybillNo, ShouldEqual, "SR201907220001")
			So(transaction.BoxNo, ShouldEqual, "01")
			So(transaction.SkuCode, ShouldEqual, "SPJA948S2230095")
			So(transaction.Qty, ShouldEqual, 2)
		})
	})
}

func TestAssignETL(t *testing.T) {
	Convey("测试AssignETL的Run方法", t, func() {
		Convey("可以把入库预约从CSL导入到MSL", func() {
			etl := AssignETL{}.New()
			err := etl.Run(context.Background())
			So(err, ShouldBeNil)
		})
	})
}
