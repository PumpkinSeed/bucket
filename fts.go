package odatas

import (
	"errors"
	"fmt"
	"github.com/couchbase/gocb"
	"github.com/volatiletech/null"
)

var (
	ErrEmptyField     = errors.New("field must be filled")
	ErrEmptyIndex = errors.New("index must be filled")

	placeholderBucket *gocb.Bucket
)

type SearchQuery struct {
	Match       null.String `json:"match,omitempty"`
	MatchPhrase null.String `json:"match_phrase,omitempty"`
	Term        null.String `json:"term,omitempty"`
	Prefix      null.String `json:"prefix,omitempty"`
	Regexp      null.String `json:"regexp,omitempty"`
	Wildcard    null.String `json:"wildcard,omitempty"`
	Bool        null.Bool   `json:"bool,omitempty"`

	Field        string      `json:"field,omitempty"`
	Analyzer     null.String `json:"analyzer,omitempty"`
	Fuzziness    null.Int64  `json:"fuzziness,omitempty"`
	PrefixLength null.Int64  `json:"prefix_length,omitempty"`
}

type CompoundQueries struct {
	Conjunction []SearchQuery `json:"conjuction,omitempty"`
	Disjunction []SearchQuery `json:"disjunction,omitempty"`
}

type RangeQuery struct {
	StartAsTime null.Time `json:"-"`
	EndAsTime   null.Time `json:"-"`
	Start       string    `json:"start,omitempty"`
	End         string    `json:"end,omitempty"`

	Min null.Int64 `json:"min,omitempty"`
	Max null.Int64 `json:"max,omitempty"`

	InclusiveStart null.Bool `json:"inclusive_start,omitempty"`
	InclusiveEnd   null.Bool `json:"inclusive_end,omitempty"`
	InclusiveMin   null.Bool `json:"inclusive_min,omitempty"`
	InclusiveMax   null.Bool `json:"inclusive_max,omitempty"`

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
		cb, err := gocb.Connect("couchbase://localhost")
		if err != nil {
			panic(err)
		}

		err = cb.Authenticate(gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		})
		if err != nil {
			panic(err)
		}

		placeholderBucket, err = cb.OpenBucket("company", "")
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
		fmt.Printf("%s\n", hit.Id)
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
	case s.Match.Valid && s.Match.String != "":
		s.MatchPhrase = EmptyString()
		s.Term = EmptyString()
		s.Prefix = EmptyString()
		s.Regexp = EmptyString()
		s.Wildcard = EmptyString()
		s.Bool = EmptyBool()
	case s.MatchPhrase.Valid && s.MatchPhrase.String != "":
		s.Match = EmptyString()
		s.Term = EmptyString()
		s.Prefix = EmptyString()
		s.Regexp = EmptyString()
		s.Wildcard = EmptyString()
		s.Bool = EmptyBool()
	case s.Term.Valid && s.Term.String != "":
		s.Match = EmptyString()
		s.MatchPhrase = EmptyString()
		s.Prefix = EmptyString()
		s.Regexp = EmptyString()
		s.Wildcard = EmptyString()
		s.Bool = EmptyBool()
	case s.Prefix.Valid && s.Prefix.String != "":
		s.Match = EmptyString()
		s.MatchPhrase = EmptyString()
		s.Term = EmptyString()
		s.Regexp = EmptyString()
		s.Wildcard = EmptyString()
		s.Bool = EmptyBool()
	case s.Regexp.Valid && s.Regexp.String != "":
		s.Match = EmptyString()
		s.MatchPhrase = EmptyString()
		s.Term = EmptyString()
		s.Prefix = EmptyString()
		s.Wildcard = EmptyString()
		s.Bool = EmptyBool()
	case s.Wildcard.Valid && s.Wildcard.String != "":
		s.Match = EmptyString()
		s.MatchPhrase = EmptyString()
		s.Term = EmptyString()
		s.Prefix = EmptyString()
		s.Regexp = EmptyString()
		s.Bool = EmptyBool()
	case s.Bool.Valid:
		s.Match = EmptyString()
		s.MatchPhrase = EmptyString()
		s.Term = EmptyString()
		s.Prefix = EmptyString()
		s.Regexp = EmptyString()
		s.Wildcard = EmptyString()
	}
	return nil
}
