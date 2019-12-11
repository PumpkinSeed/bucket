package bucket

import (
	"context"
	"reflect"

	"github.com/couchbase/gocb"
)

// Get retrieves a document from the bucket
func (h *Handler) Get(ctx context.Context, typ, id string, ptr interface{}) error {
	kv, err := h.get(ctx, typ, id, ptr)
	if err != nil {
		return err
	}

	var ops []gocb.BulkOp
	for k, v := range kv {
		ops = append(ops, &gocb.GetOp{Key: k.Key, Value: v})
	}

	return h.state.bucket.Do(ops)
}

// GetAndTouch retrieves a document and simultaneously updates its expiry times
func (h *Handler) GetAndTouch(ctx context.Context, typ, id string, ptr interface{}, ttl uint32) error {
	kv, err := h.get(ctx, typ, id, ptr)
	if err != nil {
		return err
	}

	var ops []gocb.BulkOp
	for k, v := range kv {
		ops = append(ops, &gocb.GetAndTouchOp{Key: k.Key, Value: v, Expiry: ttl})
	}

	return h.state.bucket.Do(ops)
}

func (h *Handler) get(ctx context.Context, typ, id string, ptr interface{}) (map[documentMeta]interface{}, error) {
	// checks for invalid input, ptr must be a pointer
	if err := h.inputcheck(ptr); err != nil {
		return nil, err
	}

	// returns the list of available documents must setup by the BulkOp
	kv, err := h.availableDocuments(ctx, typ, id, ptr)
	if err != nil {
		return nil, err
	}

	// prepare the resultset what the Get looks for
	lookFor := h.prepareResultSet(kv)

	fields, err := h.lookForNestedFields(ptr, lookFor)
	if err != nil {
		return nil, err
	}

	// setup document key and value pointer pairs for GetOp
	for k := range kv {
		if field, ok := fields[k.Type]; ok && field != nil {
			kv[k] = field
		}
	}

	return kv, nil
}

func (h *Handler) inputcheck(ptr interface{}) error {
	if reflect.ValueOf(ptr).Kind() != reflect.Ptr {
		return ErrInputStructPointer
	}

	return nil
}

// getAllMeta read the document meta field and
func (h *Handler) availableDocuments(tx context.Context, typ, id string, ptr interface{}) (map[documentMeta]interface{}, error) {
	var kv = make(map[documentMeta]interface{})
	dk := h.state.getDocumentKey(typ, id)

	key := documentMeta{
		Key:  dk,
		Type: typ,
		ID:   id,
	}
	kv[key] = ptr

	m, err := h.getMeta(typ, id)
	if err != nil {
		return nil, err
	}

	for _, rdm := range m.ChildDocuments {
		kv[rdm] = nil
	}

	return kv, nil
}

func (h *Handler) prepareResultSet(kv map[documentMeta]interface{}) map[string]interface{} {
	var result = make(map[string]interface{})
	for rdm, elem := range kv {
		if elem == nil {
			result[rdm.Type] = nil
		}
	}

	return result
}

func (h *Handler) lookForNestedFields(ptr interface{}, fields map[string]interface{}) (map[string]interface{}, error) {
	// get reflection of result
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
			var cont bool
			var err error
			fields, cont, err = h.gfield(rvQField, rtQField, fields)
			if err != nil {
				return nil, err
			}
			if cont {
				continue
			}
		}
	}

	return fields, nil
}

// gfieldcheck checks the certain field in Get method
func (h *Handler) gfieldcheck(rvQField reflect.Value, rtQField reflect.StructField, fields map[string]interface{}) (string, bool, error) {
	refTag, hasRefTag := rtQField.Tag.Lookup(tagReferenced)

	// if the struct isn't referenced or it's referenced but it's not a struct
	// see more: Rule #1
	if !hasRefTag || rvQField.Type().Elem().Kind() != reflect.Struct {
		return "", true, nil
	}

	// if it referenced and struct then must be a reference tag filled with value
	if refTag == "" {
		return "", false, ErrEmptyRefTag
	}

	// if a referenced struct wasn't added at insert then continue, prevent nil overwrites
	if _, ok := fields[refTag]; !ok {
		return "", true, nil
	}

	return refTag, false, nil
}

func (h *Handler) gfield(rvQField reflect.Value, rtQField reflect.StructField, fields map[string]interface{}) (map[string]interface{}, bool, error) {
	// check field
	refTag, cont, err := h.gfieldcheck(rvQField, rtQField, fields)
	if err != nil {
		return fields, false, err
	}
	if cont {
		return fields, true, err
	}

	// rvQField initialized with it's own type
	rvQField.Set(reflect.New(rvQField.Type().Elem()))

	// passed to the fields to be set up by BulkOp
	fields[refTag] = rvQField.Addr().Interface()

	// look for nested fields in struct
	fields, err = h.lookForNestedFields(rvQField.Interface(), fields)
	if err != nil {
		return nil, false, err
	}

	// then return the constructed fields
	return fields, false, nil
}
