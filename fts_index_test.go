package bucket

import (
	"context"
	"testing"
)

func TestCreateFullTextSearchIndex(t *testing.T) {
	indexName := "order_fts_idx"

	if ok, _, _ := th.InspectFullTextSearchIndex(context.Background(), indexName); ok {
		err := th.DeleteFullTextSearchIndex(context.Background(), indexName)
		if err != nil {
			t.Fatal(err)
		}
	}

	def, err := DefaultFullTextSearchIndexDefinition(IndexMeta{
		Name:                 indexName,
		SourceType:           "couchbase",
		SourceName:           "company",
		DocIDPrefixDelimiter: "::",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = th.CreateFullTextSearchIndex(context.Background(), def)
	if err != nil {
		t.Fatal(err)
	}
}
