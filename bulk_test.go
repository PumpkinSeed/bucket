package bucket

import (
	"context"
	"fmt"
	"testing"

	"github.com/rs/xid"
)

func TestGetBulk(t *testing.T) {
	for i := 0; i < 10; i++ {
		ws := generate()
		_, _, err := th.Insert(context.Background(), "webshop", xid.New().String(), ws, 0)
		//_, err := th.state.bucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := createFullTextSearchIndex("webshop_fts_idx", false); err != nil {
		t.Fatal(err)
	}

	searchMatch := "processed"
	res, err := th.SimpleSearch(context.Background(), "webshop_fts_idx", &SearchQuery{
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
	fmt.Println(ws)
}
