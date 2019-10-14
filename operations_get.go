package bucket

import (
	"context"
	"reflect"

	"github.com/couchbase/gocb"
)

// Get retrieves a document from the bucket
func (h *Handler) Get(ctx context.Context, typ, id string, ptr interface{}) error {
	if _, err := h.state.getDocumentKey(typ, id); err != nil {
		return err
	}

	// checks for invalid input, ptr must be a pointer
	if err := h.inputcheck(ptr); err != nil {
		return err
	}

	// getAllMeta
	kv, err := h.getAllMeta(ctx, typ, id, ptr)
	if err != nil {
		return err
	}

	lookFor := h.getTypesWhereValueIsNil(kv)

	fields, err := h.lookForNestedFields(ptr, lookFor)
	if err != nil {
		return err
	}

	for k := range kv {
		if field, ok := fields[k.Type]; ok && field != nil {
			kv[k] = field
		}
	}

	var ops []gocb.BulkOp
	for k, v := range kv {
		ops = append(ops, &gocb.GetOp{Key: k.Key, Value: v})
	}

	return h.state.bucket.Do(ops)

	//return nil
}

func (h *Handler) inputcheck(ptr interface{}) error {
	if reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		return ErrInputStructPointer
	}

	return nil
}

// getAllMeta read the document meta field and
func (h *Handler) getAllMeta(tx context.Context, typ, id string, ptr interface{}) (map[referencedDocumentMeta]interface{}, error) {
	var kv = make(map[referencedDocumentMeta]interface{})
	dk, err := h.state.getDocumentKey(typ, id)
	if err != nil {
		return nil, err
	}

	key := referencedDocumentMeta{
		Key:  dk,
		Type: typ,
		ID:   id,
	}
	kv[key] = ptr

	m, err := h.getMeta(typ, id)
	if err != nil {
		return nil, err
	}

	for _, rdm := range m.ReferencedDocuments {
		kv[rdm] = nil
	}

	return kv, nil
}

func (h *Handler) getTypesWhereValueIsNil(kv map[referencedDocumentMeta]interface{}) map[string]interface{} {
	var result = make(map[string]interface{})
	for rdm, elem := range kv {
		if elem == nil {
			result[rdm.Type] = nil
		}
	}

	return result
}

func (h *Handler) lookForNestedFields(ptr interface{}, fields map[string]interface{}) (map[string]interface{}, error) {
	rv := reflect.ValueOf(ptr)
	rt := rv.Type()

	if rt.Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
		rt = rv.Type()
	}

	for i := 0; i < rt.NumField(); i++ {
		rvQField := rv.Field(i)
		rtQField := rt.Field(i)
		if rvQField.Kind() == reflect.Ptr {
			refTag, hasRefTag := rtQField.Tag.Lookup(tagReferenced)
			if !hasRefTag || rvQField.Type().Elem().Kind() != reflect.Struct {
				continue
			}
			if refTag == "" {
				return nil, ErrEmptyRefTag
			}
			if _, ok := fields[refTag]; !ok {
				continue
			}
			rvQField.Set(reflect.New(rvQField.Type().Elem()))
			fields[refTag] = rvQField.Addr().Interface()
			var err error
			fields, err = h.lookForNestedFields(rvQField.Interface(), fields)
			if err != nil {
				return fields, err
			}
		}
	}

	return fields, nil
}
