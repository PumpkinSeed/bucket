package odatas

import (
	"log"
	"reflect"

	"github.com/couchbase/gocb"
)

const (
	TagJson      = "json"
	TagIndexable = "indexable"
	//TagReferenced = "referenced" // referenced tag represents external id-s
)

type indexer struct {
	bucket        *gocb.Bucket
	bucketManager *gocb.BucketManager
}

func (h *Handler) Index(v interface{}) error {
	if err := h.bucketManager.CreatePrimaryIndex("", true, false); err != nil {
		log.Fatalf("Error when create primary index %+v", err)
	}

	t := reflect.TypeOf(v)

	indexables := goDeep(t)

	if err := makeIndex(h.bucketManager, t.Name(), indexables); err != nil {
		return err
	}
	return nil
}

func goDeep(t reflect.Type) (indexed []string) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type.Kind() == reflect.Struct {
			goDeep(f.Type)
		}
		if f.Tag != "" {
			if json := f.Tag.Get(TagJson); json != "" && json != "-" {
				if f.Tag.Get(TagIndexable) != "" {
					indexed = append(indexed, json)
				}
			}
		}
	}
	return indexed
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
