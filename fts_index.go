package bucket

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/volatiletech/null"
)

const (
	ftsEndpoint = "/_p/fts/api/index"
	statAddress = "http://localhost:8094"
)

type apiResponse struct {
	Status    string    `json:"status"`
	IndexDefs IndexDefs `json:"indexDefs"`
	Error     string    `json:"error"`
}

// IndexDefs ...
type IndexDefs struct {
	UUID      string                     `json:"uuid"`
	IndexDefs map[string]IndexDefinition `json:"indexDefs"`
}

// IndexDefinition ...
type IndexDefinition struct {
	Type         string          `json:"type"`
	Name         string          `json:"name"`
	UUID         string          `json:"uuid"`
	SourceType   string          `json:"sourceType"`
	SourceName   string          `json:"sourceName"`
	SourceUUID   string          `json:"sourceUUID"`
	SourceParams interface{}     `json:"sourceParams"` // TODO
	PlanParams   IndexPlanParams `json:"planParams"`
	Params       IndexParams     `json:"params"`
}

// IndexPlanParams ...
type IndexPlanParams struct {
	MaxPartitionsPerPIndex int64 `json:"maxPartitionsPerPIndex"`
	NumReplicas            int64 `json:"numReplicas"`
}

// IndexParams ...
type IndexParams struct {
	DocConfig IndexDocConfig `json:"doc_config"`
	Mapping   IndexMapping   `json:"mapping"`
	Store     IndexStore     `json:"store"`
}

// IndexDocConfig ...
type IndexDocConfig struct {
	DocIDPrefixDelimiter string `json:"docid_prefix_delim"`
	DocIDRegexp          string `json:"docid_regexp"`
	Mode                 string `json:"mode"`
	TypeField            string `json:"type_field"`
}

// IndexMapping ...
type IndexMapping struct {
	DefaultAnalyzer       string               `json:"default_analyzer"`
	DefaultDatetimeParser string               `json:"default_datetime_parser"`
	DefaultField          string               `json:"default_field"`
	DefaultMapping        IndexDefaultMapping  `json:"default_mapping"`
	DefaultType           string               `json:"default_type"`
	DocvaluesDynamic      bool                 `json:"docvalues_dynamic"`
	IndexDynamic          bool                 `json:"index_dynamic"`
	StoreDynamic          bool                 `json:"store_dynamic"`
	TypeField             string               `json:"type_field"`
	Types                 map[string]IndexType `json:"types"`
}

// IndexDefaultMapping ...
type IndexDefaultMapping struct {
	Dynamic bool `json:"dynamic"`
	Enabled bool `json:"enabled"`
}

// IndexType ...
type IndexType struct {
	Dynamic         bool                       `json:"dynamic"`
	Enabled         bool                       `json:"enabled"`
	DefaultAnalyzer string                     `json:"default_analyzer,omitempty"`
	Properties      map[string]IndexProperties `json:"properties"`
}

// IndexProperties ...
type IndexProperties struct {
	Dynamic bool         `json:"dynamic"`
	Enabled bool         `json:"enabled"`
	Fields  []IndexField `json:"fields"`
}

// IndexField ...
type IndexField struct {
	Analyzer           string `json:"analyzer"`
	IncludeInAll       bool   `json:"include_in_all"`
	IncludeTermVectors bool   `json:"include_term_vectors"`
	Index              bool   `json:"index"`
	Name               string `json:"name"`
	Store              bool   `json:"store"`
	Type               string `json:"type"`
}

// IndexStore ...
type IndexStore struct {
	IndexType   string `json:"indexType"`
	KVStoreName string `json:"kvStoreName"`
}

// IndexMeta ...
type IndexMeta struct {
	Name                 string
	SourceType           string
	SourceName           string
	DocIDPrefixDelimiter string
	DocIDRegexp          string
	TypeField            string
}

// IndexCount represents index count response
type IndexCount struct {
	Status  string      `json:"status"`
	Count   null.Uint   `json:"count,omitempty"`
	Error   null.String `json:"error,omitempty"`
	Request null.String `json:"request,omitempty"`
}

// IndexStat represents the statistics of the search index
type IndexStat struct {
	Status     null.String `json:"status,omitempty"`
	Error      null.String `json:"error,omitempty"`
	Request    null.String `json:"request,omitempty"`
	AggStats   null.JSON   `json:"aggStats,omitempty"`
	DocCount   null.Uint   `json:"docCount,omitempty"`
	NodesStats null.JSON   `json:"nodesStats"`
}

// DefaultFullTextSearchIndexDefinition creates a default index def
// for full-text search and return it in purpose to change default
// values manually
func DefaultFullTextSearchIndexDefinition(meta IndexMeta) (*IndexDefinition, error) {
	if meta.Name == "" {
		return nil, ErrEmptyIndex
	}
	if meta.SourceType == "" {
		return nil, ErrEmptyType
	}
	if meta.SourceName == "" {
		return nil, ErrEmptySource
	}

	var ftsDef = &IndexDefinition{
		Type:       "fulltext-index",
		Name:       meta.Name,
		SourceType: meta.SourceType,
		SourceName: meta.SourceName,
		PlanParams: IndexPlanParams{
			MaxPartitionsPerPIndex: 171,
		},
		Params: IndexParams{
			Mapping: IndexMapping{
				DefaultAnalyzer:       "standard",
				DefaultDatetimeParser: "dateTimeOptional",
				DefaultField:          "_all",
				DefaultMapping: IndexDefaultMapping{
					Dynamic: false,
					Enabled: false,
				},
				DefaultType:      "_default",
				DocvaluesDynamic: true,
				IndexDynamic:     true,
				StoreDynamic:     true,
				TypeField:        "_type",
			},
			Store: IndexStore{
				IndexType:   "scorch",
				KVStoreName: "",
			},
		},
	}

	switch {
	case meta.DocIDPrefixDelimiter != "":
		ftsDef.Params.DocConfig = IndexDocConfig{
			DocIDPrefixDelimiter: meta.DocIDPrefixDelimiter,
			Mode:                 "docid_prefix",
			DocIDRegexp:          "",
			TypeField:            "type",
		}
	case meta.DocIDRegexp != "":
		ftsDef.Params.DocConfig = IndexDocConfig{
			DocIDPrefixDelimiter: "",
			Mode:                 "docid_regexp",
			DocIDRegexp:          meta.DocIDRegexp,
			TypeField:            "type",
		}
	case meta.TypeField != "":
		ftsDef.Params.DocConfig = IndexDocConfig{
			DocIDPrefixDelimiter: "",
			Mode:                 "type_field",
			DocIDRegexp:          "",
			TypeField:            meta.TypeField,
		}
	}

	return ftsDef, nil
}

// CreateFullTextSearchIndex ...
func (h *Handler) CreateFullTextSearchIndex(ctx context.Context, def *IndexDefinition) error {
	body, err := json.Marshal(def)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("PUT", h.fullTextSearchURL(ctx, def.Name), bytes.NewBuffer(body))
	setupBasicAuth(req)
	req.Header.Add("Content-Type", "application/json")
	resp, err := h.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var ar apiResponse
	if err := json.Unmarshal(respbody, &ar); err != nil {
		return err
	}
	if ar.Status == "fail" {
		return errors.New(ar.Error)
	}

	return nil
}

// DeleteFullTextSearchIndex ...
func (h *Handler) DeleteFullTextSearchIndex(ctx context.Context, indexName string) error {
	req, _ := http.NewRequest("DELETE", h.fullTextSearchURL(ctx, indexName), nil)
	setupBasicAuth(req)
	req.Header.Add("Content-Type", "application/json")

	resp, err := h.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var ar apiResponse
	err = json.Unmarshal(respbody, &ar)
	if err != nil {
		return err
	}
	if ar.Status == "fail" {
		return errors.New(ar.Error)
	}

	return nil
}

// InspectFullTextSearchIndex checks the availability of the index
// and returns it if exists
func (h *Handler) InspectFullTextSearchIndex(ctx context.Context, indexName string) (bool, *IndexDefinition, error) {
	req, _ := http.NewRequest("GET", h.fullTextSearchURL(ctx, ""), nil)
	setupBasicAuth(req)
	req.Header.Add("Content-Type", "application/json")

	resp, err := h.http.Do(req)
	if err != nil {
		return false, nil, err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, nil, err
	}
	var ar apiResponse
	err = json.Unmarshal(respbody, &ar)
	if err != nil {
		return false, nil, err
	}
	if v, ok := ar.IndexDefs.IndexDefs[indexName]; ok {
		return true, &v, nil
	}
	return false, nil, nil
}

func (h *Handler) fullTextSearchURL(ctx context.Context, indexName string) string {
	if indexName == "" {
		return fmt.Sprintf("%s%s", h.httpAddress, ftsEndpoint)
	}
	return fmt.Sprintf("%s%s/%s", h.httpAddress, ftsEndpoint, indexName)
}

func (h *Handler) CountIndex(ctx context.Context, indexName string) (*IndexCount, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/index/%s/count", statAddress, indexName), nil)
	setupBasicAuth(req)
	req.Header.Add("Content-Type", "application/json")

	resp, err := h.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var indexCount IndexCount
	if err := json.Unmarshal(respbody, &indexCount); err != nil {
		return nil, err
	}

	return &indexCount, nil
}

func (h *Handler) IndexStat(ctx context.Context, indexName string) (*IndexStat, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/stats/sourceStats/%s", statAddress, indexName), nil)
	setupBasicAuth(req)
	req.Header.Add("Content-Type", "application/json")

	resp, err := h.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var indexStat IndexStat
	if err := json.Unmarshal(respbody, &indexStat); err != nil {
		return nil, err
	}

	return &indexStat, nil
}
