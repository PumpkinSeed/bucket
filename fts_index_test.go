package odatas

import "testing"

func TestCreateFullTextSearchIndex(t *testing.T) {
	indexName := "order_fts_idx"

	if ok, _, _ := th.InspectFullTextSearchIndex(indexName); ok {
		err := th.DeleteFullTextSearchIndex(indexName)
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
	err = th.CreateFullTextSearchIndex(def)
	if err != nil {
		t.Fatal(err)
	}
}