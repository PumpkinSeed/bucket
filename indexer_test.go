package bucket

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/couchbase/gocb"
)

func TestIndexCreate(t *testing.T) {
	fmt.Println("hola")
	type webshopWithNonPointerNestedStruct struct {
		webshop
		Something  string `json:"something" cb_indexable:"true"`
		NestedData struct {
			Data1 int `json:"data_1"`
		}
	}
	instance := webshopWithNonPointerNestedStruct{}

	if err := th.Index(context.Background(), instance); err != nil {
		t.Fatal(err)
	}

	indexes, err := th.GetManager(context.Background()).GetIndexes()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 5, len(indexes))

	for _, ind := range indexes {
		t.Logf("%+v", ind.Name)
	}
}

func TestPrimaryIndexCreateError(t *testing.T) {
	h := defaultHandler()

	_ = h.state.bucket.Close()
	assert.NotNil(t, h.Index(context.Background(), webshop{}))
}

func TestExistingIndex(t *testing.T) {
	if err := th.Index(context.Background(), webshop{}); err != nil {
		t.Fatal(err)
	}
	if err := th.Index(context.Background(), webshop{}); err != nil {
		t.Fatal(err)
	}
}

func TestMakeIndex(t *testing.T) {
	assert.Nil(t, makeIndex(th.GetManager(context.Background()), "randomIndexName", []string{"randomField"}))
	assert.Nil(t, th.GetManager(context.Background()).DropIndex("randomIndexName", true))
}

func TestMakeIndexMissingIndexName(t *testing.T) {
	h := defaultHandler()
	assert.NotNil(t, makeIndex(h.GetManager(context.Background()), "", nil))
}

func TestDropAndCreateMissingIndexName(t *testing.T) {
	h := defaultHandler()
	assert.NotNil(t, dropAndCreateIndex(h.GetManager(context.Background()), "", nil))
}

func BenchmarkCreateIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		instance := webshop{}

		if err := th.Index(context.Background(), instance); err != nil {
			b.Fatal(err)
		}
		//indexes, _ := th.GetManager(context.Background()).GetIndexes()
		//for _, ind := range indexes {
		//	_ = th.GetManager(context.Background()).DropIndex(ind.Name, true)
		//}
	}
}

func BenchmarkWithIndex(b *testing.B) {
	if err := th.Index(context.Background(), webshop{}); err != nil {
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
	if err := th.Index(context.Background(), webshop{}); err != nil {
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
	resp, err := th.state.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery("select * from `company` where CONTAINS(email, $1)"), []interface{}{"a"})
	if err != nil {
		return start, nil, err
	}
	_ = resp.Close()
	return start, resp, nil
}

func searchNotIndexedProperty(t *testing.T) (time.Time, gocb.QueryResults, error) {
	start := time.Now()
	resp, err := th.state.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery("select * from `company` where CONTAINS(billing_address_address_2, $1)"), []interface{}{"a"})
	if err != nil {
		return start, nil, err
	}
	_ = resp.Close()
	return start, resp, nil
}
