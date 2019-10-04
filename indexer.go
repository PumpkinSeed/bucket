package bucket

import (
	"context"
	"log"
	"reflect"

	"github.com/couchbase/gocb"
)

const (
	tagJSON       = "json"
	tagIndexable  = "cb_indexable"
	tagReferenced = "cb_referenced" // referenced tag represents external types for id-s
)

//Index runs trough the given interface v and creates secondary indexes for all the with indexable:"true" tags
func (h *Handler) Index(ctx context.Context, v interface{}) error {
	if err := h.GetManager(ctx).CreatePrimaryIndex("", true, false); err != nil {
		log.Fatalf("Error when create primary index %+v", err)
	}

	t := reflect.TypeOf(v)

	indexables := make(map[string][]string)
	goDeep(t, indexables)

	for key, val := range indexables {
		if err := makeIndex(h.GetManager(ctx), key, val); err != nil {
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
			if json := removeOmitempty(f.Tag.Get(tagJSON)); json != "" && json != "-" {
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
			if err := manager.DropIndex(indexName, true); err != nil {
				log.Printf("Error when dropping index[%s] %+v", indexName, err)
				return err
			}
			if err := manager.CreateIndex(indexName, indexedFields, false, false); err != nil {
				log.Printf("Error when create secondary index %+v", err)
				return err
			}
		} else {
			log.Printf("Error when create secondary index %+v", err)
			return err
		}
	}
	return nil
}
