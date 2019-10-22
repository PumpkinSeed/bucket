package main

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/couchbase/gocb"

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
	profileLoad()
	touch()
	ping()
}

func preloadAll() {
	preload(ProfileType, func() interface{} {
		return models.GenerateProfile()
	})

	log.Println("Wait 2 seconds to get ready...")
	time.Sleep(2 * time.Second)

	preload(EventType, func() interface{} {
		th.SetDocumentType(context.Background(), "event", "sim_event")
		th.SetDocumentType(context.Background(), "event_location", "sim_event_location")
		th.SetDocumentType(context.Background(), "event_photo", "sim_event_photo")
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

func profileLoad() {
	var quantity, loadTimeSum = 10000, 0

	var profiles = make(map[string]*models.Profile)
	ctx := context.Background()
	log.Printf("Start inserting new profiles... \n")
	start := time.Now()
	for i := 0; i < quantity; i++ {
		profile := models.GenerateProfile()
		_, id, err := th.Insert(ctx, ProfileType, "", profile, 0)
		if err != nil {
			log.Fatal(err)
		}
		profiles[id] = profile
	}
	for id, v := range profiles {
		loadTime := time.Now()
		profile := &models.Profile{}
		if err := th.Get(ctx, ProfileType, id, profile); err != nil {
			log.Fatal(err)
		}
		loadTimeSum += int(time.Since(loadTime))

		if !reflect.DeepEqual(v, profile) {
			log.Printf("%+v\n%+v\n are not equals", v, profile)
		}
	}

	log.Printf("All profile load time: %v", time.Since(start))
	log.Printf("Single load operation spent AVG: %v\n", time.Duration(loadTimeSum/len(profiles)))
}

func touch() {
	var quantity = 10000
	var timeToLive = uint32(15)
	ctx := context.Background()
	var store = make(map[string]*models.Order)
	var c = 1

	log.Printf("Start inserting new orders... \n")
	start := time.Now()
	insert := time.Now()

	for i := 0; i < quantity; i++ {
		order := models.GenerateOrder()
		_, id, err := th.Insert(ctx, OrderType, "", order, timeToLive)
		if err != nil {
			log.Fatal(err)
		}
		store[id] = order
	}
	log.Printf("New orders inserted in: %v", time.Since(insert))

	tTouch := time.Now()
	var newStore = make(map[string]*models.Order)
	for id, v := range store {
		if c%2 == 0 {
			err := th.Touch(ctx, OrderType, id, v, 0)
			if err != nil {
				log.Fatal(err)
			}

			newStore[id] = v
		}
		c++
	}
	log.Printf("Orders touched in: %v", time.Since(tTouch))

	time.Sleep(time.Duration(timeToLive))

	var results []string
	for id, v := range newStore {
		err := th.Get(ctx, OrderType, id, v)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, id)
	}

	if len(results) != quantity/2 {
		log.Fatal("touch failed")
	}

	log.Printf("Touch and check finished in: %v", time.Since(start)-time.Duration(timeToLive))
}

func ping() {
	ctx := context.Background()
	var services []gocb.ServiceType
	services = append(services, gocb.FtsService)
	report, err := th.Ping(ctx, services)
	if err != nil {
		log.Fatal(err)
	}

	for _, ser := range report.Services {
		if ser.Service != gocb.FtsService && ser.Success != true {
			log.Fatal("Full text search not available")
		}
	}
}
