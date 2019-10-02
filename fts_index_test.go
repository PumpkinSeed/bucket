package bucket

import (
	"testing"
)

func TestCreateFullTextSearchIndex(t *testing.T) {
	if err := createFullTextSearchIndex("order_fts_idx", true); err != nil {
		t.Fatal(err)
	}
}
