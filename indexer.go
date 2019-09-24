package odatas

import (
	"log"
	"reflect"

	"github.com/couchbase/gocb"
)

const (
	tagJson      = "json"
	tagIndexable = "indexable"
	//TagReferenced = "referenced" // referenced tag represents external id-s
)

func (h *Handler) Index(v interface{}) error {
	if err := h.GetManager().CreatePrimaryIndex("", true, false); err != nil {
		log.Fatalf("Error when create primary index %+v", err)
	}

	t := reflect.TypeOf(v)

	indexables := make(map[string][]string)
	goDeep(t, indexables)

	for key, val := range indexables {
		if err := makeIndex(h.GetManager(), key, val); err != nil {
			return err
		}
	}
	return nil
}

func goDeep(t reflect.Type, indexables map[string][]string) {
	indexables[t.Name()] = []string{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type.Kind() == reflect.Struct {
			goDeep(f.Type, indexables)
		}
		if f.Tag != "" {
			if json := f.Tag.Get(tagJson); json != "" && json != "-" {
				if f.Tag.Get(tagIndexable) != "" {
					indexables[t.Name()] = append(indexables[t.Name()], json)
				}
			}
		}
	}
	if len(indexables[t.Name()]) == 0 {
		delete(indexables, t.Name())
	}
}

func makeIndex(manager *gocb.BucketManager, indexName string, indexedFields []string) error {
	if err := manager.CreateIndex(indexName, indexedFields, false, false); err != nil {
		if err == gocb.ErrIndexAlreadyExists {
			_ = manager.DropIndex(indexName, true)
			return makeIndex(manager, indexName, indexedFields)
		} else {
			log.Printf("Error when create secondary index %+v", err)
			return err
		}
	}
	return nil
}
