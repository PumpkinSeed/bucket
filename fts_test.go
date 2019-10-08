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

	err := cr.setup()
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

	err := rq.setup()
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

//TODO
func TestSimpleSearchMatch(t *testing.T) {
	for i := 0; i < 10; i++ {
		order := generate()
		_, err := th.state.bucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			t.Fatal(err)
		}
	}
	_ = createFullTextSearchIndex("order_fts_simple_idx", false)
	searchMatch := "processed"
	mes := time.Now()
	_, err := th.SimpleSearch(context.Background(), "order_fts_simple_idx", &SearchQuery{
		Query: searchMatch,
		//Field: "CardHolderName",
	})
	fmt.Println(time.Since(mes))
	if err != nil {
		t.Fatal(err)
	}
	//for _, a := range res {
	//	fmt.Println(a.Id, a.Score)
	//	//resp = append(resp, a.Id)
	//}
}

func TestSimpleSearchMatchInvalidIndex(t *testing.T) {
	_, err := th.SimpleSearch(context.Background(), "order_fts_simple_invalid_asdaas_idx", &SearchQuery{
		Query: "processed",
		//Field: "CardHolderName",
	})

	assert.NotNil(t, err)
}

func TestSimpleSearchWithoutField(t *testing.T) {
	_, err := th.SimpleSearch(
		context.Background(),
		"",
		&SearchQuery{
			Query: "",
		})

	assert.Equal(t, ErrEmptyField, err)
}

func TestSimpleSearchWithoutIndex(t *testing.T) {
	_, err := th.SimpleSearch(
		context.Background(),
		"",
		&SearchQuery{
			Query: "asd",
		})

	assert.Equal(t, ErrEmptyIndex, err)
}

//TODO
func TestSimpleSearchMatchWithFacet(t *testing.T) {
	for i := 0; i < 10; i++ {
		order := generate()
		_, err := th.state.bucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			t.Fatal(err)
		}
	}
	_ = createFullTextSearchIndex("order_fts_simple_facet_idx", false)
	searchMatch := "Talia"
	mes := time.Now()
	_, _, err := th.SimpleSearchWithFacets(
		context.Background(),
		"order_fts_simple_facet_idx",
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

func TestSimpleSearchMatchWithFacetInvalidIndex(t *testing.T) {
	_, _, err := th.SimpleSearchWithFacets(
		context.Background(),
		"order_fts_simple_facet_random_asdadasd_idx",
		&SearchQuery{
			Query: "Talia",
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
	assert.NotNil(t, err)
}

func TestSimpleSearchWithFacetsWithoutField(t *testing.T) {
	_, _, err := th.SimpleSearchWithFacets(
		context.Background(),
		"",
		&SearchQuery{
			Query: "",
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

	assert.Equal(t, ErrEmptyField, err)
}

func TestSimpleSearchWithFacetsWithoutIndex(t *testing.T) {
	_, _, err := th.SimpleSearchWithFacets(
		context.Background(),
		"",
		&SearchQuery{
			Query: "asd",
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
	assert.Equal(t, ErrEmptyIndex, err)
}

//TODO
func TestCompoundSearchConjunction(t *testing.T) {
	for i := 0; i < 10; i++ {
		order := generate()
		_, err := th.state.bucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	_ = createFullTextSearchIndex("order_compound_conj_fts_idx", false)
	mes := time.Now()

	res, err := th.CompoundSearch(context.Background(),
		"order_compound_conj_fts_idx",
		&CompoundQueries{
			Conjunction: []SearchQuery{
				{
					Query: "card",
				},
				{
					Query: "processed",
				},
			},
		})
	fmt.Println(time.Since(mes))

	assert.Nil(t, err)
	t.Logf("%+v", res)
}

//TODO
func TestCompoundSearchDisjunction(t *testing.T) {
	for i := 0; i < 10; i++ {
		order := generate()
		_, err := th.state.bucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	_ = createFullTextSearchIndex("order_compound_disj_fts_idx", false)
	mes := time.Now()

	res, err := th.CompoundSearch(context.Background(),
		"order_compound_disj_fts_idx",
		&CompoundQueries{
			Disjunction: []SearchQuery{
				{
					Query: "processed",
				},
				{
					Query: "failed",
				},
			},
		})
	fmt.Println(time.Since(mes))

	assert.Nil(t, err)
	t.Logf("%+v", res)
}

func TestCompoundSearchDisjunctionInvalidIndex(t *testing.T) {
	_, err := th.CompoundSearch(context.Background(),
		"order_compound_disj_fts_asdsad_random_idx",
		&CompoundQueries{
			Disjunction: []SearchQuery{
				{
					Query: "processed",
				},
				{
					Query: "failed",
				},
			},
		})
	assert.NotNil(t, err)
}

func TestCompoundSearchMissingQuery(t *testing.T) {
	_, err := th.CompoundSearch(context.Background(),
		"order_compound_conj_fts_idx",
		&CompoundQueries{})
	assert.NotNil(t, err)
}

func TestCompoundSearchWithoutIndex(t *testing.T) {
	_, err := th.CompoundSearch(context.Background(),
		"",
		&CompoundQueries{
			Conjunction: []SearchQuery{
				{
					Query: "card",
				},
				{
					Query: "processed",
				},
			},
		})
	assert.NotNil(t, err)
}

//TODO
func TestCompoundSearchWithFacetDisjunction(t *testing.T) {
	for i := 0; i < 10; i++ {
		order := generate()
		_, err := th.state.bucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	_ = createFullTextSearchIndex("order_compound_disj_facet_fts_idx", false)
	mes := time.Now()

	_, _, err := th.CompoundSearchWithFacets(context.Background(),
		"order_compound_disj_facet_fts_idx",
		&CompoundQueries{
			Disjunction: []SearchQuery{
				{
					Query: "processed",
				},
			},
		},
		[]FacetDef{
			{
				Name:  "status",
				Type:  FacetTerm,
				Field: "status",
				Size:  10,
			},
		})
	fmt.Println(time.Since(mes))

	assert.Nil(t, err)
}

func TestCompoundSearchWithFacetDisjunctionInvalidFacet(t *testing.T) {
	for i := 0; i < 10; i++ {
		order := generate()
		_, err := th.state.bucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	_ = createFullTextSearchIndex("order_compound_disj_facet_fts_idx", false)
	mes := time.Now()

	_, _, err := th.CompoundSearchWithFacets(context.Background(),
		"order_compound_disj_facet_fts_idxasd",
		&CompoundQueries{
			Disjunction: []SearchQuery{
				{
					Query: "processed",
				},
			},
		},
		[]FacetDef{
			{
				Name:  "BillingAddressAddress1",
				Type:  FacetTerm,
				Field: "BillingAddressAddress1",
				Size:  10,
			},
		})
	fmt.Println(time.Since(mes))

	assert.NotNil(t, err)
}

func TestCompoundSearchWithFacetMissingQuery(t *testing.T) {
	_, _, err := th.CompoundSearchWithFacets(context.Background(),
		"order_compound_disj_facet_fts_idx",
		&CompoundQueries{},
		[]FacetDef{
			{
				Name:  "status",
				Type:  FacetTerm,
				Field: "status",
				Size:  10,
			},
		})
	assert.NotNil(t, err)
}

func TestCompoundSearchWithFacetWithoutIndex(t *testing.T) {
	_, _, err := th.CompoundSearchWithFacets(context.Background(),
		"",
		&CompoundQueries{
			Conjunction: []SearchQuery{
				{
					Query: "card",
				},
				{
					Query: "processed",
				},
			},
		}, []FacetDef{
			{
				Name:  "status",
				Type:  FacetTerm,
				Field: "status",
				Size:  10,
			},
		})
	assert.NotNil(t, err)
}

//TODO
func TestRangeSearch(t *testing.T) {
	for i := 0; i < 10; i++ {
		order := generate()
		_, err := th.state.bucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := createFullTextSearchIndex("order_fts_range_idx", true); err != nil {
		t.Fatal(err)
	}
	mes := time.Now()
	_, err := th.RangeSearch(context.Background(), "order_fts_range_idx", &RangeQuery{
		StartAsTime: time.Now().Add(-2000 * time.Hour),
		EndAsTime:   time.Now().Add(-500 * time.Hour),
		Field:       "something",
		//Field: "CardHolderName",
	})
	fmt.Println(time.Since(mes))

	assert.Nil(t, err)
}

func TestRangeSearchInvalidIndex(t *testing.T) {
	_, err := th.RangeSearch(context.Background(), "random_range_index", &RangeQuery{
		StartAsTime: time.Now().Add(-2000 * time.Hour),
		EndAsTime:   time.Now().Add(-500 * time.Hour),
		Field:       "something",
		//Field: "CardHolderName",
	})
	assert.NotNil(t, err)
}

func TestRangeSearchWithoutField(t *testing.T) {
	_, err := th.RangeSearch(context.Background(),
		"",
		&RangeQuery{
			StartAsTime: time.Now().Add(-2000 * time.Hour),
			EndAsTime:   time.Now().Add(-500 * time.Hour),
			Field:       "",
			//Field: "CardHolderName",
		})

	assert.Equal(t, ErrEmptyField, err)
}

func TestRangeSearchWithoutIndex(t *testing.T) {
	_, err := th.RangeSearch(context.Background(),
		"",
		&RangeQuery{
			StartAsTime: time.Now().Add(-2000 * time.Hour),
			EndAsTime:   time.Now().Add(-500 * time.Hour),
			Field:       "something",
			//Field: "CardHolderName",
		})

	assert.Equal(t, ErrEmptyIndex, err)
}

func TestRangeSearchWithFacet(t *testing.T) {
	for i := 0; i < 10; i++ {
		order := generate()
		_, err := th.state.bucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			t.Fatal(err)
		}
	}
	_ = createFullTextSearchIndex("order_fts_range_facet_idx", false)
	mes := time.Now()
	_, _, err := th.RangeSearchWithFacets(context.Background(),
		"order_fts_range_idx",
		&RangeQuery{
			StartAsTime: time.Now().Add(-2000 * time.Hour),
			EndAsTime:   time.Now().Add(-500 * time.Hour),
			Field:       "something",
			//Field: "CardHolderName",
		},
		[]FacetDef{
			{
				Name:  "BillingAddressAddress1",
				Type:  FacetTerm,
				Field: "BillingAddressAddress1",
				Size:  10,
			},
		})
	fmt.Println(time.Since(mes))

	assert.Nil(t, err)
}

func TestRangeSearchWithFacetInvalidIndex(t *testing.T) {
	_, _, err := th.RangeSearchWithFacets(context.Background(),
		"order_fts_range_facet_alksdja_idx",
		&RangeQuery{
			StartAsTime: time.Now().Add(-2000 * time.Hour),
			EndAsTime:   time.Now().Add(-500 * time.Hour),
			Field:       "something",
			//Field: "CardHolderName",
		},
		[]FacetDef{
			{
				Name:  "BillingAddressAddress1",
				Type:  FacetTerm,
				Field: "BillingAddressAddress1",
				Size:  10,
			},
		})

	assert.NotNil(t, err)
}

func TestRangeSearchWithFacetWithoutField(t *testing.T) {
	_, _, err := th.RangeSearchWithFacets(context.Background(),
		"",
		&RangeQuery{
			StartAsTime: time.Now().Add(-2000 * time.Hour),
			EndAsTime:   time.Now().Add(-500 * time.Hour),
			Field:       "",
		},
		[]FacetDef{
			{
				Name:  "BillingAddressAddress1",
				Type:  FacetTerm,
				Field: "BillingAddressAddress1",
				Size:  10,
			},
		})

	assert.Equal(t, ErrEmptyField, err)
}

func TestRangeSearchWithFacetWithoutIndex(t *testing.T) {
	_, _, err := th.RangeSearchWithFacets(context.Background(),
		"",
		&RangeQuery{
			StartAsTime: time.Now().Add(-2000 * time.Hour),
			EndAsTime:   time.Now().Add(-500 * time.Hour),
			Field:       "something",
		},
		[]FacetDef{
			{
				Name:  "BillingAddressAddress1",
				Type:  FacetTerm,
				Field: "BillingAddressAddress1",
				Size:  10,
			},
		})

	assert.Equal(t, ErrEmptyIndex, err)
}

func TestSetupMatch(t *testing.T) {
	query := &SearchQuery{Match: "asd", Field: "field"}
	assert.Nil(t, query.setup())
}

func TestSetupMatchPharse(t *testing.T) {
	query := &SearchQuery{MatchPhrase: "asd", Field: "field"}
	assert.Nil(t, query.setup())
}

func TestSetupTerm(t *testing.T) {
	query := &SearchQuery{Term: "asd", Field: "field"}
	assert.Nil(t, query.setup())
}

func TestSetupPrefix(t *testing.T) {
	query := &SearchQuery{Prefix: "asd", Field: "field"}
	assert.Nil(t, query.setup())
}

func TestSetupRegexp(t *testing.T) {
	query := &SearchQuery{Regexp: "asd\\d*", Field: "field"}
	assert.Nil(t, query.setup())
}

func TestSetupWildcard(t *testing.T) {
	query := &SearchQuery{Wildcard: "*", Field: "field"}
	assert.Nil(t, query.setup())
}
