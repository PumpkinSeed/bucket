package odatas

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/couchbase/gocb"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	ErrEmptyField = errors.New("field must be filled")
	ErrEmptyIndex = errors.New("index must be filled")

	placeholderBucket  *gocb.Bucket
	placeholderCluster *gocb.Cluster
)

type SearchQuery struct {
	Query       string `json:"query,omitempty"`
	Match       string `json:"match,omitempty"`
	MatchPhrase string `json:"match_phrase,omitempty"`
	Term        string `json:"term,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
	Regexp      string `json:"regexp,omitempty"`
	Wildcard    string `json:"wildcard,omitempty"`
	Bool        bool   `json:"bool,omitempty"`

	Field        string `json:"field,omitempty"`
	Analyzer     string `json:"analyzer,omitempty"`
	Fuzziness    int64  `json:"fuzziness,omitempty"`
	PrefixLength int64  `json:"prefix_length,omitempty"`
}

type CompoundQueries struct {
	Conjunction []SearchQuery `json:"conjuction,omitempty"`
	Disjunction []SearchQuery `json:"disjunction,omitempty"`
}

type RangeQuery struct {
	StartAsTime time.Time `json:"-"`
	EndAsTime   time.Time `json:"-"`
	Start       string    `json:"start,omitempty"`
	End         string    `json:"end,omitempty"`

	Min int64 `json:"min,omitempty"`
	Max int64 `json:"max,omitempty"`

	InclusiveStart bool `json:"inclusive_start,omitempty"`
	InclusiveEnd   bool `json:"inclusive_end,omitempty"`
	InclusiveMin   bool `json:"inclusive_min,omitempty"`
	InclusiveMax   bool `json:"inclusive_max,omitempty"`

	Field string `json:"field,omitempty"`
}

// time RFC-3339
//{
//"start": "2001-10-09T10:20:30-08:00",
//"end": "2016-10-31",
//"inclusive_start": false,
//"inclusive_end": false,
//"field": "review_date"
//}
//{
//"min": 100, "max": 1000,
//"inclusive_min": false,
//"inclusive_max": false,
//"field": "id"
//}

func placeholderInit() {
	if placeholderBucket == nil {
		var err error
		placeholderCluster, err = gocb.Connect("couchbase://localhost")
		if err != nil {
			panic(err)
		}

		err = placeholderCluster.Authenticate(gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		})
		if err != nil {
			panic(err)
		}

		placeholderBucket, err = placeholderCluster.OpenBucket("company", "")
		if err != nil {
			panic(err)
		}
	}
}

func (h *Handler) SimpleSearch(index string, q *SearchQuery) error {
	placeholderInit()

	if err := q.Setup(); err != nil {
		return err
	}

	if index == "" {
		return ErrEmptyIndex
	}

	query := gocb.NewSearchQuery(index, q)
	res, err := placeholderBucket.ExecuteSearchQuery(query)
	if err != nil {
		return err
	}
	fmt.Println(res.Status())
	for _, hit := range res.Hits() {
		fmt.Printf("%+v\n", hit)
	}
	for _, facet := range res.Facets() {
		fmt.Printf("%+v\n", facet)
	}

	return nil
}

func (h *Handler) CompoundSearch(doc string, q *CompoundQueries) {

}

func (h *Handler) RangeSearch(doc string, q *RangeQuery) {

}

func (s *SearchQuery) Setup() error {
	if s.Field == "" {
		return ErrEmptyField
	}

	switch {
	case s.Match != "":
		s.MatchPhrase = emptyString()
		s.Term = emptyString()
		s.Prefix = emptyString()
		s.Regexp = emptyString()
		s.Wildcard = emptyString()
		s.Bool = emptyBool()
	case s.MatchPhrase != "":
		s.Match = emptyString()
		s.Term = emptyString()
		s.Prefix = emptyString()
		s.Regexp = emptyString()
		s.Wildcard = emptyString()
		s.Bool = emptyBool()
	case s.Term != "":
		s.Match = emptyString()
		s.MatchPhrase = emptyString()
		s.Prefix = emptyString()
		s.Regexp = emptyString()
		s.Wildcard = emptyString()
		s.Bool = emptyBool()
	case s.Prefix != "":
		s.Match = emptyString()
		s.MatchPhrase = emptyString()
		s.Term = emptyString()
		s.Regexp = emptyString()
		s.Wildcard = emptyString()
		s.Bool = emptyBool()
	case s.Regexp != "":
		s.Match = emptyString()
		s.MatchPhrase = emptyString()
		s.Term = emptyString()
		s.Prefix = emptyString()
		s.Wildcard = emptyString()
		s.Bool = emptyBool()
	case s.Wildcard != "":
		s.Match = emptyString()
		s.MatchPhrase = emptyString()
		s.Term = emptyString()
		s.Prefix = emptyString()
		s.Regexp = emptyString()
		s.Bool = emptyBool()
		//case s.Bool.Valid:
		//	s.Match = emptyString()
		//	s.MatchPhrase = emptyString()
		//	s.Term = emptyString()
		//	s.Prefix = emptyString()
		//	s.Regexp = emptyString()
		//	s.Wildcard = emptyString()
	}
	return nil
}

/*
	Index of FTS
*/

const (
	ftsEndpoint = "/_p/fts/api/index"
)

type FullTextSearchIndexDefinition struct {
	Type       string                        `json:"type"`
	Name       string                        `json:"name"`
	SourceType string                        `json:"sourceType"`
	SourceName string                        `json:"sourceName"`
	PlanParams FullTextSearchIndexPlanParams `json:"planParams"`
	Params     FullTextSearchIndexParams     `json:"params"`
}

type FullTextSearchIndexPlanParams struct {
	MaxPartitionsPerPIndex int64 `json:"maxPartitionsPerPIndex"`
}

type FullTextSearchIndexParams struct {
	DocConfig FullTextSearchIndexDocConfig `json:"doc_config"`
	Mapping   FullTextSearchIndexMapping   `json:"mapping"`
	Store     FullTextSearchIndexStore     `json:"store"`
}

type FullTextSearchIndexDocConfig struct {
	DocIDPrefixDelimiter string `json:"docid_prefix_delim"`
	DocIDRegexp          string `json:"docid_regexp"`
	Mode                 string `json:"mode"`
	TypeField            string `json:"type_field"`
}

type FullTextSearchIndexMapping struct {
	DefaultAnalyzer       string                            `json:"default_analyzer"`
	DefaultDatetimeParser string                            `json:"default_datetime_parser"`
	DefaultField          string                            `json:"default_field"`
	DefaultMapping        FullTextSearchIndexDefaultMapping `json:"default_mapping"`
	DefaultType           string                            `json:"default_type"`
	DocvaluesDynamic      bool                              `json:"docvalues_dynamic"`
	IndexDynamic          bool                              `json:"index_dynamic"`
	StoreDynamic          bool                              `json:"store_dynamic"`
	TypeField             string                            `json:"type_field"`
}

type FullTextSearchIndexDefaultMapping struct {
	Dynamic bool `json:"dynamic"`
	Enabled bool `json:"enabled"`
}

type FullTextSearchIndexStore struct {
	IndexType   string `json:"indexType"`
	KVStoreName string `json:"kvStoreName"`
}

type FullTextSearchIndexMeta struct {
	Name                 string
	SourceType           string
	SourceName           string
	DocIDPrefixDelimiter string
	DocIDRegexp          string
	TypeField            string
}

func DefaultFullTextSearchIndexDefinition(meta FullTextSearchIndexMeta) (*FullTextSearchIndexDefinition, error) {
	if meta.Name == "" {
		return nil, errors.New("index name must set")
	}
	if meta.SourceType == "" {
		return nil, errors.New("source type must set")
	}
	if meta.SourceName == "" {
		return nil, errors.New("source name must set")
	}

	var ftsDef = &FullTextSearchIndexDefinition{
		Type:       "fulltext-index",
		Name:       meta.Name,
		SourceType: meta.SourceType,
		SourceName: meta.SourceName,
		PlanParams: FullTextSearchIndexPlanParams{
			MaxPartitionsPerPIndex: 171,
		},
		Params: FullTextSearchIndexParams{
			Mapping: FullTextSearchIndexMapping{
				DefaultAnalyzer:       "standard",
				DefaultDatetimeParser: "dateTimeOptional",
				DefaultField:          "_all",
				DefaultMapping: FullTextSearchIndexDefaultMapping{
					Dynamic: true,
					Enabled: true,
				},
				DefaultType:      "_default",
				DocvaluesDynamic: true,
				IndexDynamic:     true,
				StoreDynamic:     true,
				TypeField:        "_type",
			},
			Store: FullTextSearchIndexStore{
				IndexType:   "scorch",
				KVStoreName: "",
			},
		},
	}

	switch {
	case meta.DocIDPrefixDelimiter != "":
		ftsDef.Params.DocConfig = FullTextSearchIndexDocConfig{
			DocIDPrefixDelimiter: meta.DocIDPrefixDelimiter,
			Mode:                 "docid_prefix",
			DocIDRegexp:          "",
			TypeField:            "",
		}
	case meta.DocIDRegexp != "":
		ftsDef.Params.DocConfig = FullTextSearchIndexDocConfig{
			DocIDPrefixDelimiter: "",
			Mode:                 "docid_regexp",
			DocIDRegexp:          meta.DocIDRegexp,
			TypeField:            "",
		}
	case meta.TypeField != "":
		ftsDef.Params.DocConfig = FullTextSearchIndexDocConfig{
			DocIDPrefixDelimiter: "",
			Mode:                 "type_field",
			DocIDRegexp:          "",
			TypeField:            meta.TypeField,
		}
	}

	return ftsDef, nil
}

func (h *Handler) CreateFullTextSearchIndex(def *FullTextSearchIndexDefinition) error {
	body, err := json.Marshal(def)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("PUT", h.fullTestSearchURL(def.Name), bytes.NewBuffer(body))
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
	log.Println(string(respbody))

	return nil
}

func (h *Handler) DeleteFullTextSearchIndex(indexName string) error {
	req, _ := http.NewRequest("DELETE", h.fullTestSearchURL(indexName), nil)
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
	log.Println(string(respbody))

	return nil
}

func (h *Handler) InspectFullTextSearchIndex(indexName string) (bool, error) {
	req, _ := http.NewRequest("GET", h.fullTestSearchURL(""), nil)
	setupBasicAuth(req)
	req.Header.Add("Content-Type", "application/json")

	resp, err := h.http.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	log.Println(string(respbody))
	return false, nil
}

func (h *Handler) fullTestSearchURL(indexName string) string {
	if indexName == "" {
		return fmt.Sprintf("%s%s", h.httpAddress, ftsEndpoint)
	}
	return fmt.Sprintf("%s%s/%s", h.httpAddress, ftsEndpoint, indexName)
}