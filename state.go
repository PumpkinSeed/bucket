package odatas

import (
	"github.com/couchbase/gocb"
)

const (
	stateDocumentKey = "bucket_state"
)

type state struct {
	DocumentTypes map[string]string `json:"document_types"`

	bucket    *gocb.Bucket `json:"-"`
	Separator string `json:"separator"`
}

func newState(c *Configuration) (*state, error) {
	cluster, err := gocb.Connect(c.ConnectionString)
	if err != nil {
		return nil, err
	}
	err = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: c.Username,
		Password: c.Password,
	})
	if err != nil {
		return nil, err
	}

	bucket, err := cluster.OpenBucket(c.BucketName, c.BucketPassword)
	if err != nil {
		return nil, err
	}

	var s = &state{}
	s.DocumentTypes = make(map[string]string)
	s.bucket = bucket
	s.Separator = c.Separator

	return s, nil
}

func (s *state) load() error {
	_, err := s.bucket.Get(stateDocumentKey, s)
	return err
}

func (s *state) newType(name, prefix string) error {
	if _, ok := s.DocumentTypes[name]; ok {
		return ErrDocumentTypeAlredyExists
	}

	s.DocumentTypes[name] = prefix + s.Separator
	err := s.updateState()
	if err != nil {
		return err
	}

	return nil
}

func (s *state) getType(name string) (string, error) {
	if v, ok := s.DocumentTypes[name]; ok {
		return v, nil
	}
	return "", ErrDocumentTypeDoesntExists
}

func (s *state) deleteType(name string) error {
	if _, ok := s.DocumentTypes[name]; ok {
		delete(s.DocumentTypes, name)

		err := s.updateState()
		if err != nil {
			return err
		}

		return nil
	}

	return ErrDocumentTypeDoesntExists
}

func (s *state) updateState() error {
	_, err := s.bucket.Upsert(stateDocumentKey, s, 0)
	return err
}
