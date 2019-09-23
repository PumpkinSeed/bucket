package odatas

import (
	"errors"
	"gopkg.in/couchbase/gocb.v1"
)

type State struct {
	DocumentTypes map[string]string

	bucket *gocb.Bucket
}


type Wrapper struct {
	ID string
	Data interface {}
}

func NewState(bucket *gocb.Bucket) *State {
	var s State
	s.DocumentTypes = make(map[string]string)
	s.bucket = bucket

	return &s
}

func LoadState(bucket *gocb.Bucket) *State {
	var s State
	s.DocumentTypes = make(map[string]string)
	//s.bucket = bucket.


	return &s
}

func (s *State) NewType (name, prefix string) error {
	if _, ok := s.DocumentTypes[name]; ok {
		return errors.New("document type already exists")
	}

	s.DocumentTypes[name] = prefix
	return nil
}

func (s *State) GetType (name string) (string, error) {
	if v, ok := s.DocumentTypes[name]; ok {
		return v, nil
	}
	return "", errors.New("document type doesn't exist")
}

func (s *State) DeleteType (name string) error {
	if _, ok := s.DocumentTypes[name]; ok {
		delete(s.DocumentTypes, name)
		return nil
	}

	return errors.New("document type doesn't exist")
}