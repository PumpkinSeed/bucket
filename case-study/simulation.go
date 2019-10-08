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

	profileSelection()
}

func preloadAll() {
	preload(ProfileType, func() interface{} {
		return models.GenerateProfile()
	})

	log.Println("Wait 2 seconds to get ready...")
	time.Sleep(2 * time.Second)

	preload(EventType, func() interface{} {
		_ = th.SetDocumentType(context.Background(), "event", "sim_event")
		_ = th.SetDocumentType(context.Background(), "event_location", "sim_event_location")
		_ = th.SetDocumentType(context.Background(), "event_photo", "sim_event_photo")
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

func profileSelection() {
	var quantity, affiliationLimit = 10000, 55
	ids := make([]string, quantity)
	ctx := context.Background()
	log.Printf("Start inserting new profiles... \n")
	for i := 0; i < quantity/2; i++ {
		profile := models.GenerateProfile()
		_, id, err := th.Insert(ctx, ProfileType, "", profile, 0)
		if err != nil {
			log.Fatal(err)
		}
		ids[i] = id
	}
	for i := 0; i < quantity/2; i++ {
		profile := models.GenerateProfile()
		_, id, err := th.Upsert(ctx, ProfileType, "", profile, 0)
		if err != nil {
			log.Fatal(err)
		}
		ids[quantity/2+i] = id
	}
	log.Printf("Start profile selection... \n")
	upsertTimeSum, removeTimeSum := 0, 0
	upsertNum, removeNum := 0, 0
	for i := 0; i < quantity; i++ {
		profile := models.Profile{}
		err := th.Get(ctx, ProfileType, ids[i], &profile)
		if err != nil {
			log.Fatal(err)
		}
		if profile.AffiliationCount < uint64(affiliationLimit) {
			profile.Status = models.GenerateStatus()
			mes := time.Now()
			_, _, err := th.Upsert(ctx, ProfileType, ids[i], profile, 0)
			if err != nil {
				log.Fatal(err)
			}
			upsertTimeSum += int(time.Since(mes))
			upsertNum++
		} else {
			mes := time.Now()
			err := th.Remove(ctx, ProfileType, ids[i], &profile)
			if err != nil {
				log.Fatal(err)
			}
			removeTimeSum += int(time.Since(mes))
			removeNum++
		}
	}
	log.Printf("Number of upserts: %d \n", upsertNum)
	log.Printf("Number of removes: %d \n", removeNum)
	log.Printf("Total spent %s \n", time.Duration(removeTimeSum+upsertTimeSum))
	log.Printf("Single upsert operation spent AVG: %s \n", time.Duration(upsertTimeSum/upsertNum))
	log.Printf("Single remove operation spent AVG: %s \n", time.Duration(removeTimeSum/removeNum))

}
