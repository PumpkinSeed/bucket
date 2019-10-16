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

	kv := h.getSubDocuments(typ, id, q, nil)

	var ops []gocb.BulkOp
	for k, v := range kv {
		key, err := h.state.getDocumentKey(k, id)
		if err != nil {
			//return nil, "", err
		}
		ops = append(ops, &gocb.InsertOp{Key: key, Value: v, Expiry: ttl})
	}

	err := h.state.bucket.Do(ops)
	return nil, "", err
}

func (h *Handler) getSubDocuments(typ, id string, q interface{}, parent *documentMeta) map[string]map[string]interface{} {
	var documents = make(map[string]map[string]interface{})
	var metaField = &meta{
		ParentDocument: parent,
		Type:           typ,
	}

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
			currentKey, _ := h.state.getDocumentKey(typ, id)
			current := documentMeta{
				Type: typ,
				ID:   id,
				Key:  currentKey,
			}
			subDocuments := h.getSubDocuments(tag, id, rvField.Interface(), &current)
			for k, v := range subDocuments {
				childKey, _ := h.state.getDocumentKey(tag, id)
				metaField.AddChildDocument(childKey, k, id)
				documents[k] = v
			}
		} else {
			if j, ok := rtField.Tag.Lookup(tagJSON); ok && j != "-" {
				fields[j] = rvField.Interface()
			}
		}
	}
	fields[metaFieldName] = metaField
	documents[typ] = fields

	return documents
}
