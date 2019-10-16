package bucket

import (
	"context"
	"reflect"

	"github.com/couchbase/gocb"
	"github.com/rs/xid"
)

func (h *Handler) EInsert(ctx context.Context, typ, id string, q interface{}, ttl uint32) (Cas, string, error) {
	if id == "" {
		id = xid.New().String()
	}

	var ops []gocb.BulkOp
	//for k, v := range kv {
	//	ops = append(ops, &gocb.InsertOp{Key: k.Key, Value: v, Expiry: ttl})
	//}

	err := h.state.bucket.Do(ops)
	return nil, "", err
}

func (h *Handler) getSubDocuments(typ string, q interface{}) map[string]map[string]interface{} {
	var documents = make(map[string]map[string]interface{})

	var rv = reflect.ValueOf(q)
	var rt = rv.Type()

	if rv.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
		rt = rv.Type()
	}

	var fields = make(map[string]interface{})
	for i := 0; i < rt.NumField(); i++ {
		rvField := rv.Field(i)
		rtField := rt.Field(i)
		if tag, ok := rtField.Tag.Lookup(tagReferenced); ok {
			subDocuments := h.getSubDocuments(tag, rvField.Interface())
			for k, v := range subDocuments {
				documents[k] = v
			}
		} else {
			if j, ok := rtField.Tag.Lookup(tagJSON); ok && j != "-" {
				fields[j] = rvField.Interface()
			}
		}
	}
	documents[typ] = fields

	return documents
}
