package odatas

import (
	"fmt"
	"testing"
	"time"

	"github.com/couchbase/gocb"
)

const (
	bucketName = "company"
)

var h *Handler

func init() {
	h, _ = New(&Configuration{
		Username:       "Administrator",
		Password:       "password",
		BucketName:     bucketName,
		BucketPassword: "",
	})

	start := time.Now()
	if err := h.GetManager().Flush(); err != nil {
		fmt.Printf("Turn on flush in bucket: %+v\n", err)
	}
	fmt.Printf("Bucket flushed: %v\n", time.Since(start))

	for j := 0; j < 10000; j++ {
		instance := newTestStruct1()
		_, _ = h.state.bucket.Insert(instance.Token, instance, 0)
	}
	fmt.Printf("Connection setup, data seeded %v\n", time.Since(start))
}

func TestIndexCreate(t *testing.T) {
	instance := testStructEmbedded{}

	if err := h.Index(instance); err != nil {
		t.Fatal(err)
	}

	indexes, err := h.GetManager().GetIndexes()
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
	if err := h.Index(testStruct1{}); err != nil {
		t.Fatal(err)
	}

	start, resp, err := searchIndexedProperty(t)
	if err != nil {
		t.Fatalf("One search time: %v\n%+v", start, err)
	}
	fmt.Printf("One search time: %v\nFound: %+v", time.Since(start), resp.Metrics())
}

func TestSearchWithoutIndex(t *testing.T) {
	if err := h.Index(testStruct1{}); err != nil {
		t.Fatal(err)
	}

	start, resp, err := searchNotIndexedProperty(t)
	if err != nil {
		t.Fatalf("One search time: %v\n%+v", start, err)
	}
	fmt.Printf("One search time: %v\nFound: %+v", time.Since(start), resp.Metrics())
}

func BenchmarkWithIndex(b *testing.B) {
	if err := h.Index(testStruct1{}); err != nil {
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
	if err := h.Index(testStruct1{}); err != nil {
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
	resp, err := h.state.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery("select * from `company` where CONTAINS(email, $1)"), []interface{}{"a"})
	if err != nil {
		return start, nil, err
	}
	_ = resp.Close()
	return start, resp, nil
}

func searchNotIndexedProperty(t *testing.T) (time.Time, gocb.QueryResults, error) {
	start := time.Now()
	resp, err := h.state.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery("select * from `company` where CONTAINS(billing_address_address_2, $1)"), []interface{}{"a"})
	if err != nil {
		return start, nil, err
	}
	_ = resp.Close()
	return start, resp, nil
}
