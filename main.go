package main

import (
	"clearance-adapter/domain/services"
	"clearance-adapter/factory"
	"context"
	"flag"
	"log"
)

func main() {
	flag.Parse()
	factory.Init()
	etl := services.InStorageETL{}.New("2019-07-01 00:00:00", "2019-07-31 23:59:59")
	err := etl.Run(context.Background())
	if err != nil {
		log.Println(err.Error())
	}
}
