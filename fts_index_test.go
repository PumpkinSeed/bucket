package bucket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFullTextSearchIndex(t *testing.T) {
	assert.Nil(t, createFullTextSearchIndex("order_fts_idx", true))
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
