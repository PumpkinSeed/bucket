package odatas

import (
	"errors"
	"fmt"
	"github.com/couchbase/gocb"
	"github.com/couchbase/gocb/cbft"
	"time"
)

const (
	FacetDate = iota
	FacetNumeric
	FacetTerm
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

	Limit int `json:"-"`
	Offset int `json:"-"`
}

type FacetDef struct {
	Name string
	Type int
	Field string
	Size int
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

	query := gocb.NewSearchQuery(index, q).Limit(q.Limit).Skip(q.Offset)
	return h.doSimpleSearch(query)
}

func (h *Handler) SimpleSearchWithFacets(index string, q *SearchQuery, facets []FacetDef) error {
	placeholderInit()

	if err := q.Setup(); err != nil {
		return err
	}

	if index == "" {
		return ErrEmptyIndex
	}

	query := gocb.NewSearchQuery(index, q).Limit(q.Limit).Skip(q.Offset)
	for _, facet := range facets {
		switch facet.Type {
		case FacetDate:
			query.AddFacet(facet.Name, cbft.NewDateFacet(facet.Field, facet.Size))
		case FacetNumeric:
			query.AddFacet(facet.Name, cbft.NewNumericFacet(facet.Field, facet.Size))
		case FacetTerm:
			query.AddFacet(facet.Name, cbft.NewTermFacet(facet.Field, facet.Size))
		}
	}

	return h.doSimpleSearch(query)
}

func (h *Handler) doSimpleSearch(query *gocb.SearchQuery) error {
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
	if s.Query != "" {
		s.Match = emptyString()
		s.MatchPhrase = emptyString()
		s.Term = emptyString()
		s.Prefix = emptyString()
		s.Regexp = emptyString()
		s.Wildcard = emptyString()
		s.Bool = emptyBool()
		return nil
	}

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