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
	end := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("开始同步: %s \n", start)

	distributionETL := services.DistributionETL{}.New()
	err := distributionETL.Run(context.Background())
	if err != nil {
		log.Println(err.Error())
	}

	returnToWarehouseETL := services.ReturnToWarehouseETL{}.New()
	err = returnToWarehouseETL.Run(context.Background())
	if err != nil {
		log.Println(err.Error())
	}

	transferETL := services.TransferETL{}.New()
	err = transferETL.Run(context.Background())
	if err != nil {
		log.Println(err.Error())
	}

	log.Printf("结束同步: %s \n", end)
}
