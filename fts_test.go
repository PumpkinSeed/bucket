package bucket

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
	for i := 0; i < 10; i++ {
		webshop := generate()
		_, _, err := th.Insert(context.Background(), "webshop", "", webshop, 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	err := createFullTextSearchIndex("webshop_ftx_idx_simple", false)
	if err != nil {
		t.Fatal(err)
	}

	searchMatch := "processed"
	mes := time.Now()
	res, err := th.SimpleSearch(context.Background(), "webshop_ftx_idx_simple", &SearchQuery{
		Query: searchMatch,
		//Field: "CardHolderName",
	})
	fmt.Println(time.Since(mes))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, len(res), "Length of result set should be 10")
}

func TestSimpleSearchMatchWithFacet(t *testing.T) {
	for i := 0; i < 10; i++ {
		webshop := generate()
		_, _, err := th.Insert(context.Background(), "webshop", "", webshop, 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	err := createFullTextSearchIndex("webshop_ftx_idx_simple_f", false)
	if err != nil {
		t.Fatal(err)
	}

	searchMatch := "processed"
	mes := time.Now()
	res, _, err := th.SimpleSearchWithFacets(
		context.Background(),
		"webshop_ftx_idx_simple_f",
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
	assert.Equal(t, 10, len(res), "Length of result set should be 10")
}
