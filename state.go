package odatas

import (
	"errors"

	"gopkg.in/couchbase/gocb.v1"
)

type State struct {
	DocumentTypes map[string]string

	bucket     *gocb.Bucket
	bucketName string
	separator string
}

func NewState(bucket *gocb.Bucket, bucketName, separator string) *State {
	var s State
	s.DocumentTypes = make(map[string]string)
	s.bucket = bucket
	s.bucketName = bucketName
	s.separator = separator

	return &s
}

func (s *State) Load() error {
	q := gocb.NewN1qlQuery("SELECT * FROM " + s.bucketName)
	rows, err := s.bucket.ExecuteN1qlQuery(q, nil)
	var val interface{}
	if err != nil {
		return err
	}
	for rows.Next(&val) {
		
	}

	return nil
}

func (s *State) NewType(name, prefix string) error {
	if _, ok := s.DocumentTypes[name]; ok {
		return errors.New("document type already exists")
	}

	s.DocumentTypes[name] = prefix + s.separator
	return nil
}

func (s *State) GetType(name string) (string, error) {
	if v, ok := s.DocumentTypes[name]; ok {
		return v, nil
	}
	return "", errors.New("document type doesn't exist")
}

func (s *State) DeleteType(name string) error {
	if _, ok := s.DocumentTypes[name]; ok {
		delete(s.DocumentTypes, name)
		return nil
	}

	return errors.New("document type doesn't exist")
}
