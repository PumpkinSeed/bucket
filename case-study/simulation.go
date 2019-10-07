package main

import (
	"context"
	"log"
	"time"

	"github.com/PumpkinSeed/bucket"
	"github.com/PumpkinSeed/bucket/case-study/models"
)

const (
	EventType   = "event"
	ProfileType = "profile"
	OrderType   = "order"
)

var th *bucket.Handler

func init() {
	var err error
	th, err = bucket.New(&bucket.Configuration{
		Username:         "Administrator",
		Password:         "password",
		BucketName:       "company",
		BucketPassword:   "",
		ConnectionString: "couchbase://localhost",
		Separator:        "::",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	preloadAll()
}

func preloadAll() {
	preload(ProfileType, func() interface{} {
		return models.GenerateProfile()
	})

	log.Println("Wait 2 seconds to get ready...")
	time.Sleep(2 * time.Second)

	preload(EventType, func() interface{} {
		return models.GenerateEvent()
	})

	log.Println("Wait 2 seconds to get ready...")
	time.Sleep(2 * time.Second)

	preload(OrderType, func() interface{} {
		return models.GenerateOrder()
	})
}

func preload(typ string, generator func() interface{}) {
	var quantity = 10000

	log.Printf("Start to load %d new %s.\n", quantity, typ)
	mesSum := 0
	for i := 0; i < quantity; i++ {
		data := generator()
		mes := time.Now()
		_, _, err := th.Insert(context.Background(), typ, "", data, 0)
		if err != nil {
			log.Fatal(err)
		}
		mesSum += int(time.Since(mes))
	}

	mesAvg := time.Duration(mesSum / quantity)
	log.Printf("Total spent: %s\n", time.Duration(mesSum))
	log.Printf("Single operation spent Avg: %s\n", mesAvg)
}
