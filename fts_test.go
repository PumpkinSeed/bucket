package odatas

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
)

func init() {
	gofakeit.Seed(time.Now().UnixNano())
}

func TestSearchQuery(t *testing.T) {
	sq := SearchQuery{
		Query: "card",
	}

	sqjso, err := json.Marshal(sq)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"query":"card"}`
	if string(sqjso) != expected {
		t.Errorf("Query should be %s, instead of %s", expected, string(sqjso))
	}
}

func TestConjuncts(t *testing.T) {
	cr := CompoundQueries{
		Conjunction: []SearchQuery{
			{
				Query: "card",
			},
			{
				Query: "processed",
			},
		},
	}

	err := cr.Setup()
	if err != nil {
		t.Fatal(err)
	}

	json.Marshal(cr)
	crjso, err := json.Marshal(cr)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"conjuncts":[{"query":"card"},{"query":"processed"}]}`
	if string(crjso) != expected {
		t.Errorf("Query should be %s, instead of %s", expected, string(crjso))
	}
}

func TestRangeQuery(t *testing.T) {
	rq := RangeQuery{
		StartAsTime: time.Now().Add(-2000 * time.Hour),
		EndAsTime:   time.Now().Add(-500 * time.Hour),
		Field:       "something",
	}

	err := rq.Setup()
	if err != nil {
		t.Fatal(err)
	}

	json.Marshal(rq)
	rqjso, err := json.Marshal(rq)
	if err != nil {
		t.Fatal(err)
	}
	expected := fmt.Sprintf(`{"start":"%s","end":"%s","field":"something"}`, rq.Start, rq.End)
	if string(rqjso) != expected {
		t.Errorf("Query should be %s, instead of %s", expected, string(rqjso))
	}
}

func TestSimpleSearchMatch(t *testing.T) {
	placeholderInit()

	for i := 0; i < 10; i++ {
		order := newTestStruct1()
		_, err := placeholderBucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	handler, _ := New(&Configuration{})
	searchMatch := "Talia"
	mes := time.Now()
	_, err := handler.SimpleSearch("order_fts_idx", &SearchQuery{
		Query: searchMatch,
		//Field: "CardHolderName",
	})
	fmt.Println(time.Since(mes))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSimpleSearchMatchWithFacet(t *testing.T) {
	placeholderInit()

	for i := 0; i < 10; i++ {
		order := newTestStruct1()
		_, err := placeholderBucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	handler, _ := New(&Configuration{})
	searchMatch := "Talia"
	mes := time.Now()
	_, _, err := handler.SimpleSearchWithFacets(
		"order_fts_idx",
		&SearchQuery{
			Query: searchMatch,
		},
		[]FacetDef{
			{
				Name:  "BillingAddressAddress1",
				Type:  FacetTerm,
				Field: "BillingAddressAddress1",
				Size:  10,
			},
		},
	)
	fmt.Println(time.Since(mes))
	if err != nil {
		t.Fatal(err)
	}
}
