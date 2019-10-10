package bucket

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateFullTextSearchIndex(t *testing.T) {
	assert.Nil(t, createFullTextSearchIndex("order_fts_index", true, "product"))
}

func TestDefaultFullTextSearchIndexDefinitionWithDocIDRegexp(t *testing.T) {
	ftsDef, err := DefaultFullTextSearchIndexDefinition(IndexMeta{
		Name:        "name",
		SourceType:  "source_type",
		SourceName:  "soruce_name",
		DocIDRegexp: "id::*",
	})

	assert.Nil(t, err)
	assert.Equal(t, "id::*", ftsDef.Params.DocConfig.DocIDRegexp)
}

func TestDefaultFullTextSearchIndexDefinitionWithTypeField(t *testing.T) {
	ftsDef, err := DefaultFullTextSearchIndexDefinition(IndexMeta{
		Name:       "name",
		SourceType: "source_type",
		SourceName: "soruce_name",
		TypeField:  "type_",
	})

	assert.Nil(t, err)
	assert.Equal(t, "type_", ftsDef.Params.DocConfig.TypeField)
}

func TestDefaultFullTextSearchIndexDefinitionWithoutName(t *testing.T) {
	_, err := DefaultFullTextSearchIndexDefinition(IndexMeta{})
	assert.NotNil(t, err)
}

func TestDefaultFullTextSearchIndexDefinitionWithoutSourceType(t *testing.T) {
	_, err := DefaultFullTextSearchIndexDefinition(IndexMeta{Name: "name"})
	assert.NotNil(t, err)
}

func TestDefaultFullTextSearchIndexDefinitionWithoutSourceName(t *testing.T) {
	_, err := DefaultFullTextSearchIndexDefinition(IndexMeta{Name: "name", SourceType: "source_type"})
	assert.NotNil(t, err)
}

func TestIndexCreationWithAllFieldSetup(t *testing.T) {
	const indexName = "product_fts_index"
	def, err := DefaultFullTextSearchIndexDefinition(IndexMeta{
		Name:                 indexName,
		SourceType:           "couchbase",
		SourceName:           "company",
		TypeField:            "type_",
		DocIDPrefixDelimiter: "::",
	})
	assert.Nil(t, err)
	def.Params.Mapping.Types = map[string]IndexType{
		"product": {
			Dynamic:         false,
			Enabled:         true,
			DefaultAnalyzer: "web",
			Properties: map[string]IndexProperties{
				"name": {
					Dynamic: false,
					Enabled: true,
					Fields: []IndexField{{
						Analyzer:           "web",
						IncludeInAll:       true,
						IncludeTermVectors: true,
						Index:              true,
						Name:               "name",
						Store:              false,
						Type:               "text",
					}},
				},
			},
		},
	}

	if err := th.CreateFullTextSearchIndex(context.Background(), def); err != nil {
		assert.Nil(t, th.DeleteFullTextSearchIndex(context.Background(), indexName))
		assert.Nil(t, th.CreateFullTextSearchIndex(context.Background(), def))
	}

	time.Sleep(2 * time.Second)

	_, ind, err := th.InspectFullTextSearchIndex(context.Background(), indexName)
	if err != nil {
		t.Fatal(err)
	}
	def.UUID = ind.UUID

	indJSON, _ := json.Marshal(ind)
	defindex, _ := json.Marshal(def)
	assert.JSONEq(t, string(defindex), string(indJSON))
}
