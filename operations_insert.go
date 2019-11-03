package bucket

import (
	"context"
	"reflect"

	"github.com/couchbase/gocb"
	"github.com/rs/xid"
)

func (h *Handler) Insert(ctx context.Context, typ, id string, q interface{}, ttl uint32) (Cas, string, error) {
	if id == "" {
		id = xid.New().String()
	}

	kv := h.getSubDocuments(typ, id, q, nil)

	var ops []gocb.BulkOp
	for k, v := range kv {
		key := h.state.getDocumentKey(k, id)
		ops = append(ops, &gocb.InsertOp{Key: key, Value: v, Expiry: ttl})
	}

	err := h.state.bucket.Do(ops)
	return nil, id, err
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
		if rv.IsNil() {
			return documents
		}
		rv = reflect.Indirect(rv)
		rt = rv.Type()
	}

	var fields = make(map[string]interface{})
	for i := 0; i < rt.NumField(); i++ {
		rvField := rv.Field(i)
		rtField := rt.Field(i)
		if tag, ok := rtField.Tag.Lookup(tagReferenced); ok {
			h.buildDocuments(typ, id, tag, rvField.Interface(), metaField, documents)
		} else {
			if j, ok := rtField.Tag.Lookup(tagJSON); ok && j != "-" {
				fields[removeOmitempty(j)] = rvField.Interface()
			}
		}
	}
	fields[metaFieldName] = metaField
	documents[typ] = fields

	return documents
}

func (h *Handler) buildDocuments(typ, id, tag string, sub interface{}, metaField *meta, documents map[string]map[string]interface{}) {
	currentKey := h.state.getDocumentKey(typ, id)
	current := documentMeta{
		Type: typ,
		ID:   id,
		Key:  currentKey,
	}
	subDocuments := h.getSubDocuments(tag, id, sub, &current)
	for k, v := range subDocuments {
		childKey := h.state.getDocumentKey(k, id)
		metaField.AddChildDocument(childKey, k, id)
		documents[k] = v
	}
}
