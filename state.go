package odatas

import (
	"errors"

	"github.com/couchbase/gocb"
)

type state struct {
	documentTypes map[string]string

	bucket    *gocb.Bucket
	separator string
}

func newState(bucket *gocb.Bucket, separator string) *state {
	var s state
	s.documentTypes = make(map[string]string)
	s.bucket = bucket
	s.separator = separator

	return &s
}

func (s *state) load() error {
	_, err := s.bucket.Get("documentType", s.documentTypes)
	if err != nil {
		return err
	}
	return nil
}

func (s *state) newType(name, prefix string) error {
	if _, ok := s.documentTypes[name]; ok {
		return errors.New("document type already exists")
	}

	s.documentTypes[name] = prefix + s.separator
	_, err := s.bucket.Upsert("documentType", s.documentTypes, 0)
	if err != nil {
		return err
	}

	return nil
}

func (s *state) getType(name string) (string, error) {
	if v, ok := s.documentTypes[name]; ok {
		return v, nil
	}
	return "", errors.New("document type doesn't exist")
}

func (s *state) deleteType(name string) error {
	if _, ok := s.documentTypes[name]; ok {
		delete(s.documentTypes, name)

		_, err := s.bucket.Upsert("documentType", s.documentTypes, 0)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("document type doesn't exist")
}
