package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/PumpkinSeed/bucket"
	"github.com/PumpkinSeed/bucket/case-study/models"
	"github.com/go-chi/chi"
)

var han *bucket.Handler

func init() {
	var err error
	han, err = bucket.New(&bucket.Configuration{
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
	rand.Seed(time.Now().UTC().UnixNano())
	r := chi.NewRouter()
	Routes(r)

	srv := &http.Server{
		Addr:         ":8001",
		Handler:      r,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

	go insertProfile(randInt(1, 5))
}

func randInt(min int, max int) uint {
	return uint(min + rand.Intn(max-min))
}

/*
	Profile section
*/

func insertProfile(interval uint) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		profile := models.GenerateProfile()

		profileJSON, _ := json.Marshal(profile)

		resp, err := http.Post("localhost:8001/api/profile", "application/json", bytes.NewBuffer(profileJSON))
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
	}
}

func getProfile(interval uint) {
	var i = 0
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		if len(profiles) <= i {
			continue
		}

		url := "localhost:8001/api/profile/" + profiles[i]

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
	}
}

func updateProfile(interval uint) {
	var i = 0
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		if len(profiles) <= i {
			continue
		}

		url := "localhost:8001/api/profile/" + profiles[i]

		profile := models.GenerateProfile()

		profileJSON, _ := json.Marshal(profile)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(profileJSON))
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
	}
}

/*
	Event section
*/

func insertEvent(interval uint) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		event := models.GenerateEvent()

		eventJSON, _ := json.Marshal(event)

		resp, err := http.Post("localhost:8001/api/event", "application/json", bytes.NewBuffer(eventJSON))
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
	}
}

func getEvent(interval uint) {
	var i = 0
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		if len(events) <= i {
			continue
		}

		url := "localhost:8001/api/profile/" + events[i]

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
	}
}

func updateEvent(interval uint) {
	var i = 0
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		if len(events) <= i {
			continue
		}

		url := "localhost:8001/api/event/" + events[i]

		event := models.GenerateEvent()

		eventJSON, _ := json.Marshal(event)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(eventJSON))
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
	}
}

/*
	Order section
*/

func insertOrder(interval uint) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		order := models.GenerateOrder()

		orderJSON, _ := json.Marshal(order)

		resp, err := http.Post("localhost:8001/api/order", "application/json", bytes.NewBuffer(orderJSON))
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
	}
}

func getOrder(interval uint) {
	var i = 0
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		if len(orders) <= i {
			continue
		}

		url := "localhost:8001/api/profile/" + orders[i]

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
	}
}

func updateOrder(interval uint) {
	var i = 0
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		if len(orders) <= i {
			continue
		}

		url := "localhost:8001/api/order/" + orders[i]

		order := models.GenerateOrder()

		orderJSON, _ := json.Marshal(order)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(orderJSON))
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
	}
}