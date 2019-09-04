package main

import (
	"clearance-adapter/domain/services"
	"clearance-adapter/factory"
	"context"
	"flag"
	"log"
	"time"
)

func main() {
	flag.Parse()
	factory.Init()

	start := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("开始同步: %s \n", start)

	var err error

	// 自动入库
	autoDistETL := services.AutoDistributionETL{}.New()
	err = autoDistETL.Run(context.Background())
	if err != nil {
		log.Println(err.Error())
	}

	distributionETL := services.DistributionETL{}.New()
	err = distributionETL.Run(context.Background())
	if err != nil {
		log.Println(err.Error())
	}

	returnToWarehouseETL := services.ReturnToWarehouseETL{}.New()
	err = returnToWarehouseETL.Run(context.Background())
	if err != nil {
		log.Println(err.Error())
	}

	// 调货自动入库
	autoTransferETL := services.AutoTransferETL{}.New()
	err = autoTransferETL.Run(context.Background())
	if err != nil {
		log.Println(err.Error())
	}

	transferETL := services.TransferETL{}.New()
	err = transferETL.Run(context.Background())
	if err != nil {
		log.Println(err.Error())
	}

	end := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("结束同步: %s \n", end)
}
