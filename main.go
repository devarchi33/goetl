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

	start := time.Now().Format("2006-01-02") + " 00:00:00"
	end := time.Now().Format("2006-01-02") + " 23:59:59"
	log.Printf("起始日期: %s, 结束日期: %s \n", start, end)

	distributionETL := services.DistributionETL{}.New(start, end)
	err := distributionETL.Run(context.Background())
	if err != nil {
		log.Println(err.Error())
	}

	returnToWarehouseETL := services.ReturnToWarehouseETL{}.New(start, end)
	err = returnToWarehouseETL.Run(context.Background())
	if err != nil {
		log.Println(err.Error())
	}
}
