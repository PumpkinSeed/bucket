package bucket

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/couchbase/gocb"
)

const (
	stateDocumentKey = "bucket_state"
)

func (h *Handler) SetDocumentType(ctx context.Context, name, prefix string) error {
	return h.state.setType(name, prefix)
}

type state struct {
	sync.RWMutex
	DocumentTypes map[string]string `json:"document_types"`

	bucket        *gocb.Bucket
	configuration *Configuration
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
	s.bucket = bucket
	s.configuration = c
	_ = s.load()

	if s.DocumentTypes == nil {
		s.DocumentTypes = make(map[string]string)
	}

	return s, nil
}

func (s *state) load() error {
	_, err := s.bucket.Get(stateDocumentKey, s)
	return err
}

func (s *state) inspect(name string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.DocumentTypes[name]
	return ok
}

func (s *state) setType(name, prefix string) error {
	s.Lock()
	defer s.Unlock()
	s.DocumentTypes[name] = prefix + s.configuration.Separator
	err := s.updateState()
	if err != nil {
		return err
	}

	return nil
}

func (s *state) getType(name string) (string, error) {
	s.RLock()
	defer s.RUnlock()
	if v, ok := s.DocumentTypes[name]; ok {
		return v, nil
	}
	return "", ErrDocumentTypeDoesntExists
}

func (s *state) fetchDocumentIdentifier(documentKey string) string {
	elems := strings.Split(documentKey, s.configuration.Separator)
	if len(elems) > 0 {
		return elems[len(elems)-1]
	}

	return ""
}

func (s *state) getDocumentKey(name, id string) (string, error) {
	typ, err := s.getType(name)
	if err != nil {
		return "", err
	}
	return typ + id, nil
}

func (s *state) deleteType(name string) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.DocumentTypes[name]; ok {
		delete(s.DocumentTypes, name)

		return s.updateState()
	}

	return ErrDocumentTypeDoesntExists
}

func (s *state) validate() (bool, error) {
	key := "doc_type"
	queryStr := fmt.Sprintf(`SELECT SPLIT(META().id, "%s")[0] %s FROM %s GROUP BY SPLIT(META().id, "%s")[0];`, s.configuration.Separator, key, s.configuration.BucketName, s.configuration.Separator)
	query := gocb.NewN1qlQuery(queryStr)
	rows, err := s.bucket.ExecuteN1qlQuery(query, nil)
	if err != nil {
		return false, err
	}
	var row map[string]string
	var docTypesInMemory = make(map[string]bool)
	for _, availableDocTypes := range s.DocumentTypes {
		docTypesInMemory[strings.Replace(availableDocTypes, s.configuration.Separator, "", -1)] = true
	}
	var missingDocTypes []string
	for rows.Next(&row) {
		if docType := row[key]; docType != stateDocumentKey {
			if _, ok := docTypesInMemory[docType]; !ok {
				missingDocTypes = append(missingDocTypes, docType)
			}
		}
	}

	if len(missingDocTypes) > 0 {
		return false, fmt.Errorf("missing doc types: [%s]", strings.Join(missingDocTypes, ", "))
	}
	if err = rows.Close(); err != nil {
		return true, err
	}
	return true, nil
}

func (s *state) updateState() error {
	_, err := s.bucket.Upsert(stateDocumentKey, s, 0)
	return err
}
