package odatas

import (
	"log"
	"reflect"

	"github.com/couchbase/gocb"
)

const (
	BUCKET_USERNAME = ""
	BUCKET_PASSWORD = ""
	TAG_JSON        = "json"
	TAG_INDEXED     = "indexed"
	TAG_REFERENCED  = "referenced"
)

type indexer struct {
	bucket        *gocb.Bucket
	bucketManager *gocb.BucketManager
}

func NewIndexer(b *gocb.Bucket) *indexer {
	return &indexer{bucket: b, bucketManager: b.Manager(BUCKET_USERNAME, BUCKET_PASSWORD)}
}

func (i *indexer) Index(v interface{}) error {
	if err := i.bucketManager.CreatePrimaryIndex("", true, false); err != nil {
		log.Fatalf("Error when create primary index %+v", err)
	}

	t := reflect.TypeOf(v)
	indexed, _ := goDeep(t)

	if err := makeIndex(i.bucketManager, t.Name(), indexed); err != nil {
		return err
	}

	return nil
}

func (i *indexer) Reindex(v interface{}) error {
	if err := i.bucketManager.CreatePrimaryIndex("", true, false); err != nil {
		log.Fatalf("Error when create primary index %+v", err)
	}

	t := reflect.TypeOf(v)
	indexName := t.Name()
	if err := i.bucketManager.DropIndex(indexName, false); err != nil {
		log.Printf("Error when dropping old secondary index %+v", err)
	}

	indexed, _ := goDeep(t)

	if err := makeIndex(i.bucketManager, indexName, indexed); err != nil {
		return err
	}

	return nil
}

func goDeep(t reflect.Type) (indexed []string, referenced []string) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type.Kind() == reflect.Struct {
			goDeep(f.Type)
		}
		if f.Tag != "" {
			if json := f.Tag.Get(TAG_JSON); json != "" && json != "-" {
				if f.Tag.Get(TAG_INDEXED) != "" {
					indexed = append(indexed, json)
				}
				if f.Tag.Get(TAG_REFERENCED) != "" {
					referenced = append(referenced, json)
				}
			}
		}
	}
	return indexed, referenced
}

func makeIndex(manager *gocb.BucketManager, indexName string, indexedFields []string) error {
	if err := manager.CreateIndex(indexName, indexedFields, false, false); err != nil {
		log.Printf("Error when create secondary index %+v", err)
		return err
	}
	return nil
}
