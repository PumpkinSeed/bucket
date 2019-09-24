package odatas

import (
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/rs/xid"

	"github.com/couchbase/gocb"
)

const (
	bucketName = "company"
)

var i *indexer

type indexableUser struct {
	ID          string `json:"id" indexable:"true"`
	Name        string `json:"name" indexable:"true"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Password    string `json:"-"`
}

func init() {
	cluster, _ := gocb.Connect("couchbase://localhost")
	_ = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	bucket, _ := cluster.OpenBucket(bucketName, "")
	i = NewIndexer(bucket, "Administrator", "password")

	start := time.Now()
	if err := i.bucketManager.Flush(); err != nil {
		fmt.Printf("Turn on flush in bucket: %+v\n", err)
	}
	fmt.Printf("Bucket flushed: %v\n", time.Since(start))

	for j := 0; j < 10000; j++ {
		instance := indexableUser{
			ID:          xid.New().String(),
			Name:        gofakeit.Name(),
			PhoneNumber: gofakeit.Phone(),
			Address:     gofakeit.Address().Address,
			Password:    "123",
		}
		_, _ = i.bucket.Insert(instance.ID, instance, 0)
	}
	fmt.Printf("Connection setup, data seeded %v\n", time.Since(start))
}

func TestIndexCreate(t *testing.T) {
	instance := indexableUser{}

	if err := i.Index(instance); err != nil {
		t.Fatal(err)
	}

	indexes, err := i.bucketManager.GetIndexes()
	if err != nil {
		t.Fatal(err)
	}

	if len(indexes) < 2 {
		t.Error("Missing indexes")
	}

	for _, ind := range indexes {
		t.Logf("%+v", ind.Name)
	}
}

func TestSearchWithIndex(t *testing.T) {
	if err := i.Index(indexableUser{}); err != nil {
		t.Fatal(err)
	}

	start, resp, err := searchIndexedProperty(t)
	if err != nil {
		t.Fatalf("One search time: %v\n%+v", start, err)
	}
	fmt.Printf("One search time: %v\nFound: %+v", time.Since(start), resp.Metrics())
}

func TestSearchWithoutIndex(t *testing.T) {
	if err := i.Index(indexableUser{}); err != nil {
		t.Fatal(err)
	}

	start, resp, err := searchNotIndexedProperty(t)
	if err != nil {
		t.Fatalf("One search time: %v\n%+v", start, err)
	}
	fmt.Printf("One search time: %v\nFound: %+v", time.Since(start), resp.Metrics())
}

func BenchmarkWithIndex(b *testing.B) {
	if err := i.Index(indexableUser{}); err != nil {
		b.Fatal(err)
	}

	globalTimer := time.Now()
	for i := 0; i < 100; i++ {
		start, resp, err := searchIndexedProperty(&testing.T{})
		if err != nil {
			b.Fatalf("One search time: %v\n%+v", start, err)
		}
		fmt.Printf("One search time: %v\nFound: %+v\n", time.Since(start), resp.Metrics())
	}
	fmt.Printf("Global time: %v\n", time.Since(globalTimer))
}

func BenchmarkWithoutIndex(b *testing.B) {
	if err := i.Index(indexableUser{}); err != nil {
		b.Fatal(err)
	}

	globalTimer := time.Now()
	for i := 0; i < 100; i++ {
		start, resp, err := searchNotIndexedProperty(&testing.T{})
		if err != nil {
			b.Fatalf("One search time: %v\n%+v", start, err)
		}
		fmt.Printf("One search time: %v\nFound: %+v\n", time.Since(start), resp.Metrics())
	}
	fmt.Printf("Global time: %v\n", time.Since(globalTimer))
}

func searchIndexedProperty(t *testing.T) (time.Time, gocb.QueryResults, error) {
	start := time.Now()
	resp, err := i.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery("select * from `company` where CONTAINS(name, $1)"), []interface{}{"a"})
	if err != nil {
		return start, nil, err
	}
	_ = resp.Close()
	return start, resp, nil
}

func searchNotIndexedProperty(t *testing.T) (time.Time, gocb.QueryResults, error) {
	start := time.Now()
	resp, err := i.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery("select * from `company` where CONTAINS(address, $1)"), []interface{}{"a"})
	if err != nil {
		return start, nil, err
	}
	_ = resp.Close()
	return start, resp, nil
}
