package odatas

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestSearchQuery(t *testing.T) {
	sq := SearchQuery{
		Query: "card",
	}

	sqjso, err := json.Marshal(sq)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(string(sqjso))

	placeholderInit()
}

func TestCreateFullTextSearchIndex(t *testing.T) {
	err := DeleteFullTextSearchIndex("order_fts_idx")
	if err != nil {
		t.Fatal(err)
	}

	def, err := DefaultFullTextSearchIndexDefinition(FullTextSearchIndexMeta{
		Name:                 "order_fts_idx",
		SourceType:           "couchbase",
		SourceName:           "company",
		DocIDPrefixDelimiter: "::",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = CreateFullTextSearchIndex(def)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSimpleSearchMatch(t *testing.T) {
	placeholderInit()

	for i := 0; i< 10; i++ {
		order := NewTestStruct1()
		_, err := placeholderBucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	handler := New(&Configuration{})
	searchMatch := "Talia Hudson"
	mes := time.Now()
	err := handler.SimpleSearch("order_fts_idx", &SearchQuery{
		Match: searchMatch,
		Field: "CardHolderName",
	})
	fmt.Println(time.Since(mes))
	if err != nil {
		t.Fatal(err)
	}
}