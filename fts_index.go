package odatas

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ftsEndpoint = "/_p/fts/api/index"
)

type apiResponse struct {
	Status    string    `json:"status"`
	IndexDefs IndexDefs `json:"indexDefs"`
	Error     string    `json:"error"`
}

type IndexDefs struct {
	UUID      string                     `json:"uuid"`
	IndexDefs map[string]IndexDefinition `json:"indexDefs"`
}

type IndexDefinition struct {
	Type       string          `json:"type"`
	Name       string          `json:"name"`
	SourceType string          `json:"sourceType"`
	SourceName string          `json:"sourceName"`
	PlanParams IndexPlanParams `json:"planParams"`
	Params     IndexParams     `json:"params"`
}

type IndexPlanParams struct {
	MaxPartitionsPerPIndex int64 `json:"maxPartitionsPerPIndex"`
}

type IndexParams struct {
	DocConfig IndexDocConfig `json:"doc_config"`
	Mapping   IndexMapping   `json:"mapping"`
	Store     IndexStore     `json:"store"`
}

type IndexDocConfig struct {
	DocIDPrefixDelimiter string `json:"docid_prefix_delim"`
	DocIDRegexp          string `json:"docid_regexp"`
	Mode                 string `json:"mode"`
	TypeField            string `json:"type_field"`
}

type IndexMapping struct {
	DefaultAnalyzer       string              `json:"default_analyzer"`
	DefaultDatetimeParser string              `json:"default_datetime_parser"`
	DefaultField          string              `json:"default_field"`
	DefaultMapping        IndexDefaultMapping `json:"default_mapping"`
	DefaultType           string              `json:"default_type"`
	DocvaluesDynamic      bool                `json:"docvalues_dynamic"`
	IndexDynamic          bool                `json:"index_dynamic"`
	StoreDynamic          bool                `json:"store_dynamic"`
	TypeField             string              `json:"type_field"`
}

type IndexDefaultMapping struct {
	Dynamic bool `json:"dynamic"`
	Enabled bool `json:"enabled"`
}

type IndexStore struct {
	IndexType   string `json:"indexType"`
	KVStoreName string `json:"kvStoreName"`
}

type IndexMeta struct {
	Name                 string
	SourceType           string
	SourceName           string
	DocIDPrefixDelimiter string
	DocIDRegexp          string
	TypeField            string
}

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
					Dynamic: true,
					Enabled: true,
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
			TypeField:            "",
		}
	case meta.DocIDRegexp != "":
		ftsDef.Params.DocConfig = IndexDocConfig{
			DocIDPrefixDelimiter: "",
			Mode:                 "docid_regexp",
			DocIDRegexp:          meta.DocIDRegexp,
			TypeField:            "",
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
	err = json.Unmarshal(respbody, &ar)
	if err != nil {
		return err
	}
	if ar.Status == "fail" {
		return errors.New(ar.Error)
	}

	return nil
}

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
