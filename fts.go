package odatas

import (
	"errors"
	"fmt"
	"time"

	"github.com/couchbase/gocb"
	"github.com/couchbase/gocb/cbft"
)

const (
	FacetDate = iota
	FacetNumeric
	FacetTerm
)

var (
	ErrEmptyField = errors.New("field must be filled")
	ErrEmptyIndex = errors.New("index must be filled")
)

type SearchQuery struct {
	Query       string `json:"query,omitempty"`
	Match       string `json:"match,omitempty"`
	MatchPhrase string `json:"match_phrase,omitempty"`
	Term        string `json:"term,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
	Regexp      string `json:"regexp,omitempty"`
	Wildcard    string `json:"wildcard,omitempty"`

	Field        string `json:"field,omitempty"`
	Analyzer     string `json:"analyzer,omitempty"`
	Fuzziness    int64  `json:"fuzziness,omitempty"`
	PrefixLength int64  `json:"prefix_length,omitempty"`

	Limit  int `json:"-"`
	Offset int `json:"-"`
}

type FacetDef struct {
	Name  string
	Type  int
	Field string
	Size  int
}

type CompoundQueries struct {
	Conjunction []SearchQuery `json:"conjuncts,omitempty"`
	Disjunction []SearchQuery `json:"disjuncts,omitempty"`

	Limit  int `json:"-"`
	Offset int `json:"-"`
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

	Limit  int `json:"-"`
	Offset int `json:"-"`
}

func (h *Handler) SimpleSearch(index string, q *SearchQuery) ([]gocb.SearchResultHit, error) {
	if err := q.Setup(); err != nil {
		return nil, err
	}

	if index == "" {
		return nil, ErrEmptyIndex
	}

	query := gocb.NewSearchQuery(index, q).Limit(q.Limit).Skip(q.Offset)
	hits, _, err := h.doSearch(query)
	return hits, err
}

func (h *Handler) SimpleSearchWithFacets(index string, q *SearchQuery, facets []FacetDef) ([]gocb.SearchResultHit, map[string]gocb.SearchResultFacet, error) {
	if err := q.Setup(); err != nil {
		return nil, nil, err
	}

	if index == "" {
		return nil, nil, ErrEmptyIndex
	}

	query := gocb.NewSearchQuery(index, q).Limit(q.Limit).Skip(q.Offset)
	h.addFacets(query, facets)

	return h.doSearch(query)
}

func (h *Handler) CompoundSearch(index string, q *CompoundQueries) ([]gocb.SearchResultHit, error) {
	if err := q.Setup(); err != nil {
		return nil, err
	}

	if index == "" {
		return nil, ErrEmptyIndex
	}

	query := gocb.NewSearchQuery(index, q).Limit(q.Limit).Skip(q.Offset)
	result, _, err := h.doSearch(query)
	return result, err
}

func (h *Handler) CompoundSearchWithFacets(index string, q *CompoundQueries, facets []FacetDef) ([]gocb.SearchResultHit, map[string]gocb.SearchResultFacet, error) {
	if err := q.Setup(); err != nil {
		return nil, nil, err
	}

	if index == "" {
		return nil, nil, ErrEmptyIndex
	}

	query := gocb.NewSearchQuery(index, q).Limit(q.Limit).Skip(q.Offset)
	h.addFacets(query, facets)
	result, facetResult, err := h.doSearch(query)
	return result, facetResult, err
}

func (h *Handler) RangeSearch(index string, q *RangeQuery) ([]gocb.SearchResultHit, error) {
	if err := q.Setup(); err != nil {
		return nil, err
	}

	if index == "" {
		return nil, ErrEmptyIndex
	}

	query := gocb.NewSearchQuery(index, q).Limit(q.Limit).Skip(q.Offset)
	result, _, err := h.doSearch(query)
	return result, err
}

func (h *Handler) RangeSearchWithFacets(index string, q *RangeQuery, facets []FacetDef) ([]gocb.SearchResultHit, map[string]gocb.SearchResultFacet, error) {
	if err := q.Setup(); err != nil {
		return nil, nil, err
	}

	if index == "" {
		return nil, nil, ErrEmptyIndex
	}

	query := gocb.NewSearchQuery(index, q).Limit(q.Limit).Skip(q.Offset)
	h.addFacets(query, facets)
	result, facetResult, err := h.doSearch(query)
	return result, facetResult, err
}

func (h *Handler) doSearch(query *gocb.SearchQuery) ([]gocb.SearchResultHit, map[string]gocb.SearchResultFacet, error) {
	res, err := h.state.bucket.ExecuteSearchQuery(query)
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("%+v\n", res.Status())
	for i, v := range res.Hits() {
		fmt.Printf("%d ---- %+v\n", i, v)
	}

	return res.Hits(), res.Facets(), nil
}

func (h *Handler) addFacets(query *gocb.SearchQuery, facets []FacetDef) {
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
}

func (s *SearchQuery) Setup() error {
	if s.Query != "" {
		s.Match = ""
		s.MatchPhrase = ""
		s.Term = ""
		s.Prefix = ""
		s.Regexp = ""
		s.Wildcard = ""
		return nil
	}

	if s.Field == "" {
		return ErrEmptyField
	}

	switch {
	case s.Match != "":
		s.Query = ""
		s.MatchPhrase = ""
		s.Term = ""
		s.Prefix = ""
		s.Regexp = ""
		s.Wildcard = ""
	case s.MatchPhrase != "":
		s.Query = ""
		s.Match = ""
		s.Term = ""
		s.Prefix = ""
		s.Regexp = ""
		s.Wildcard = ""
	case s.Term != "":
		s.Query = ""
		s.Match = ""
		s.MatchPhrase = ""
		s.Prefix = ""
		s.Regexp = ""
		s.Wildcard = ""
	case s.Prefix != "":
		s.Query = ""
		s.Match = ""
		s.MatchPhrase = ""
		s.Term = ""
		s.Regexp = ""
		s.Wildcard = ""
	case s.Regexp != "":
		s.Query = ""
		s.Match = ""
		s.MatchPhrase = ""
		s.Term = ""
		s.Prefix = ""
		s.Wildcard = ""
	case s.Wildcard != "":
		s.Query = ""
		s.Match = ""
		s.MatchPhrase = ""
		s.Term = ""
		s.Prefix = ""
		s.Regexp = ""
	}
	return nil
}

func (c *CompoundQueries) Setup() error {
	if c.Conjunction == nil && c.Disjunction == nil {
		return errors.New("")
	}

	if c.Conjunction != nil {
		c.Disjunction = nil
		for _, sq := range c.Conjunction {
			err := sq.Setup()
			if err != nil {
				return err
			}
		}
	} else {
		for _, sq := range c.Disjunction {
			err := sq.Setup()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *RangeQuery) Setup() error {
	if d.Field == "" {
		return ErrEmptyField
	}

	if !d.StartAsTime.IsZero() {
		if d.EndAsTime.IsZero() {
			return errors.New("")
		}
		d.Start = d.StartAsTime.Format(time.RFC3339)
		d.End = d.EndAsTime.Format(time.RFC3339)

		d.Min = 0
		d.Max = 0
		return nil
	}

	d.Start = ""
	d.End = ""

	return nil
}
