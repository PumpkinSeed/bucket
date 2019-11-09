package bucket

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/volatiletech/null"

	"github.com/stretchr/testify/assert"
)

func TestCreateFullTextSearchIndex(t *testing.T) {
	assert.Nil(t, createFullTextSearchIndex("custom_fts_index", true, "custom"))
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
		TypeField:            "type",
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

func TestHandler_indexStat(t *testing.T) {
	type args struct {
		indexName string
	}
	tests := []struct {
		name        string
		args        args
		want        *IndexStat
		wantErr     bool
		createIndex bool
		ctx         context.Context
	}{
		{
			name: "Missing index",
			ctx:  context.Background(),
			args: args{
				indexName: "noname_random_index",
			},
			want: &IndexStat{
				Status:  null.StringFrom("fail"),
				Error:   null.StringFrom("rest_auth: preparePerms, err: index not found"),
				Request: null.StringFrom(""),
			},
			wantErr: false,
		},
		{
			name: "Currently created",
			args: args{
				indexName: "existing_index",
			},
			want: &IndexStat{
				DocCount: null.UintFrom(0),
			},
			wantErr:     false,
			createIndex: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createIndex {
				indexDef, _ := DefaultFullTextSearchIndexDefinition(IndexMeta{
					Name:                 tt.args.indexName,
					SourceType:           "couchbase",
					SourceName:           "company",
					DocIDPrefixDelimiter: "::",
					TypeField:            "type",
				})
				_ = th.CreateFullTextSearchIndex(tt.ctx, indexDef)
				time.Sleep(1 * time.Second)
			}
			got, err := th.indexStat(tt.ctx, tt.args.indexName)
			if (err != nil) != tt.wantErr {
				t.Errorf("indexStat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Fatal("missing response")
			}
			assert.NotNil(t, tt.want)
			assert.Equal(t, tt.want.Error.String, got.Error.String)
			assert.Equal(t, tt.want.Status.String, got.Status.String)
			assert.Equal(t, tt.want.Request.String, got.Request.String)
			assert.Equal(t, tt.want.DocCount.Valid, got.DocCount.Valid)

			if tt.createIndex {
				_ = th.DeleteFullTextSearchIndex(tt.ctx, tt.args.indexName)
			}
		})
	}
}

func TestHandler_countIndex(t *testing.T) {
	type args struct {
		indexName string
	}
	tests := []struct {
		name        string
		args        args
		want        *IndexCount
		wantErr     bool
		createIndex bool
		ctx         context.Context
	}{
		{
			name: "Missing index",
			args: args{
				indexName: "noname_random_index",
			},
			want: &IndexCount{
				Status:  "fail",
				Error:   null.StringFrom("rest_auth: preparePerms, err: index not found"),
				Request: null.StringFrom(""),
			},
			wantErr:     false,
			createIndex: false,
			ctx:         context.Background(),
		},
		{
			name: "Currently created",
			args: args{
				indexName: "existing_index",
			},
			want: &IndexCount{
				Status: "ok",
				Count:  null.UintFrom(0),
			},
			wantErr:     false,
			createIndex: true,
			ctx:         context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createIndex {
				indexDef, _ := DefaultFullTextSearchIndexDefinition(IndexMeta{
					Name:                 tt.args.indexName,
					SourceType:           "couchbase",
					SourceName:           "company",
					DocIDPrefixDelimiter: "::",
					TypeField:            "type",
				})
				_ = th.CreateFullTextSearchIndex(tt.ctx, indexDef)
				time.Sleep(1 * time.Second)
			}

			got, err := th.countIndex(tt.ctx, tt.args.indexName)
			if (err != nil) != tt.wantErr {
				t.Errorf("countIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Fatal("missing response")
			}
			assert.NotNil(t, tt.want)
			assert.Equal(t, tt.want.Status, got.Status)
			assert.Equal(t, tt.want.Error.String, got.Error.String)
			assert.Equal(t, tt.want.Request.String, got.Request.String)
			assert.Equal(t, tt.want.Count.Valid, got.Count.Valid)

			if tt.createIndex {
				_ = th.DeleteFullTextSearchIndex(tt.ctx, tt.args.indexName)
			}
		})
	}
}
