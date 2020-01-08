package bucket

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rs/xid"
)

func TestGetBulk(t *testing.T) {
	const webshopStatus = "success"
	for i := 0; i < 10; i++ {
		ws := generate()
		ws.Status = "success"
		_, _, err := th.EInsert(context.Background(), "webshop", xid.New().String(), ws, 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	waitUntilFtsIndexCompleted(context.Background(), "webshop_fts_index")

	searchMatch := "success"
	res, err := th.SimpleSearch(context.Background(), "webshop_fts_index", &SearchQuery{
		Query: searchMatch,
	})
	if err != nil {
		t.Fatal(err)
	}

	var ws = make([]webshop, len(res))
	err = th.GetBulk(context.Background(), res, &ws)
	if err != nil {
		t.Error(err)
	}
	if len(ws) > 0 {
		assert.Equal(t, webshopStatus, ws[0].Status)
		assert.Equal(t, "Free shipping", ws[0].ShippingMethod)
		assert.Equal(t, "active", ws[0].Product.Status)
		assert.Equal(t, "productshop", ws[0].Store.Name)
	} else {
		t.Errorf("Empty resultset of the search")
	}
}

func BenchmarkGetBulk(b *testing.B) {
	b.StopTimer()
	for i := 0; i < 10; i++ {
		ws := generate()
		_, _, err := th.EInsert(context.Background(), "webshop", xid.New().String(), ws, 0)
		//_, err := th.state.bucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			b.Fatal(err)
		}
	}

	if err := createFullTextSearchIndex("webshop_fts_index", false, "webshop"); err != nil {
		b.Fatal(err)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		searchMatch := "processed"
		res, err := th.SimpleSearch(context.Background(), "webshop_fts_idx", &SearchQuery{
			Query: searchMatch,
		})
		if err != nil {
			b.Fatal(err)
		}

		var ws = make([]webshop, len(res))
		err = th.GetBulk(context.Background(), res, &ws)
		if err != nil {
			b.Error(err)
		}
	}
}
