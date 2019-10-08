package bucket

import (
	"context"
	"reflect"

	"github.com/couchbase/gocb"
	"github.com/rs/xid"
)

type writerF func(string, string, interface{}, uint32) (gocb.Cas, error)
type readerF func(string, string, interface{}, uint32) (gocb.Cas, error)

// Cas is the container of Cas operation of all documents
type Cas map[string]gocb.Cas

// Insert inserts a new record into couchbase bucket
func (h *Handler) Insert(ctx context.Context, typ, id string, q interface{}, ttl uint32) (Cas, string, error) {
	cas := make(map[string]gocb.Cas)
	if id == "" {
		id = xid.New().String()
	}
	id, err := h.write(ctx, typ, id, q, func(typ, id string, ptr interface{}, ttl uint32) (gocb.Cas, error) {
		documentID, err := h.state.getDocumentKey(typ, id)
		if err != nil {
			return 0, err
		}
		return h.state.bucket.Insert(documentID, ptr, ttl)
	}, ttl, cas)
	if err != nil {
		return nil, "", err
	}
	return cas, id, nil
}

func (h *Handler) write(ctx context.Context, typ, id string, q interface{}, f writerF, ttl uint32, cas Cas) (string, error) {
	if !h.state.inspect(typ) {
		err := h.state.setType(typ, typ)
		if err != nil {
			return "", err
		}
	}
	fields := make(map[string]interface{})

	rvQ := reflect.ValueOf(q)
	rtQ := rvQ.Type()

	if rtQ.Kind() == reflect.Ptr {
		rvQ = reflect.Indirect(rvQ)
		rtQ = rvQ.Type()
	}

	if rtQ.Kind() == reflect.Struct {
		for i := 0; i < rvQ.NumField(); i++ {
			rvQField := rvQ.Field(i)
			rtQField := rtQ.Field(i)
			refTag, hasRefTag := rtQField.Tag.Lookup(tagReferenced)

			if rvQField.Kind() == reflect.Ptr && rvQField.IsNil() && !hasRefTag {
				if tag, ok := rtQField.Tag.Lookup(tagJSON); ok {
					fields[removeOmitempty(tag)] = nil
				}
			} else {
				if rvQField.Kind() == reflect.Ptr {
					rvQField = reflect.Indirect(rvQField)
				}
				if !rvQField.IsValid() {
					continue
				}
				if rvQField.Kind() == reflect.Struct && hasRefTag {
					if refTag == "" {
						return "", ErrEmptyRefTag
					}
					if _, err := h.write(ctx, refTag, id, rvQField.Interface(), f, ttl, cas); err != nil {
						return id, err
					}
				} else {
					if tag, ok := rtQField.Tag.Lookup(tagJSON); ok {
						fields[removeOmitempty(tag)] = rvQField.Interface()
					}
				}
			}
		}
	}
	c, err := f(typ, id, fields, ttl)
	cas[typ] = c

	return id, err
}

// Get retrieves a document from the bucket
func (h *Handler) Get(ctx context.Context, typ, id string, ptr interface{}) error {
	if err := h.read(ctx, typ, id, ptr, 0, func(typ, id string, ptr interface{}, ttl uint32) (gocb.Cas, error) {
		documentID, err := h.state.getDocumentKey(typ, id)
		if err != nil {
			return 0, err
		}
		return h.state.bucket.Get(documentID, ptr)
	}); err != nil {
		return err
	}
	return nil
}

func (h *Handler) read(ctx context.Context, typ, id string, ptr interface{}, ttl uint32, f readerF) error {
	_, err := f(typ, id, ptr, ttl)
	if err != nil {
		return err
	}

	rvQ := reflect.ValueOf(ptr)
	rtQ := rvQ.Type()

	if rtQ.Kind() == reflect.Ptr {
		rvQ = reflect.Indirect(rvQ)
		rtQ = rvQ.Type()
		if rtQ.Kind() == reflect.Struct {
			for i := 0; i < rvQ.NumField(); i++ {
				rvQField := rvQ.Field(i)
				rtQField := rtQ.Field(i)
				if rvQField.Kind() == reflect.Ptr {
					refTag, hasRefTag := rtQField.Tag.Lookup(tagReferenced)
					if !hasRefTag || rvQField.Type().Elem().Kind() != reflect.Struct {
						continue
					}
					rvQField.Set(reflect.New(rvQField.Type().Elem()))
					if refTag == "" {
						return ErrEmptyRefTag
					}
					if er := h.read(ctx, refTag, id, rvQField.Interface(), ttl, f); er != nil {
						if er != gocb.ErrKeyNotFound {
							return er
						}
						rvQField.Set(reflect.Zero(rvQField.Type()))
					}
				}
			}
		}
	} else {
		return ErrInputStructPointer
	}
	return nil
}

// Remove removes a document from the bucket
func (h *Handler) Remove(ctx context.Context, typ, id string, ptr interface{}) error {
	typs, e := getDocumentTypes(ptr)
	if e != nil {
		return e
	}

	for _, typ := range typs {
		documentID, err := h.state.getDocumentKey(typ, id)
		if err != nil {
			return err
		}
		if _, err := h.state.bucket.Remove(documentID, 0); err != nil {
			return err
		}
	}
	return nil
}

// Upsert inserts or replaces a document in the bucket
func (h *Handler) Upsert(ctx context.Context, typ, id string, q interface{}, ttl uint32) (Cas, string, error) {
	cas := make(map[string]gocb.Cas)
	if id == "" {
		id = xid.New().String()
	}
	id, err := h.write(ctx, typ, id, q, func(typ, id string, q interface{}, ttl uint32) (gocb.Cas, error) {
		documentID, err := h.state.getDocumentKey(typ, id)
		if err != nil {
			return 0, err
		}
		return h.state.bucket.Upsert(documentID, q, ttl)
	}, ttl, cas)
	if err != nil {
		return nil, "", err
	}
	return cas, id, nil

}

// Touch touches documents, specifying a new expiry time for it
// The Cas value must be 0
func (h *Handler) Touch(ctx context.Context, typ, id string, ptr interface{}, ttl uint32) error {
	typs, e := getDocumentTypes(ptr)
	if e != nil {
		return e
	}

	for _, typ := range typs {
		documentID, err := h.state.getDocumentKey(typ, id)
		if err != nil {
			return err
		}
		if _, err := h.state.bucket.Touch(documentID, 0, ttl); err != nil {
			return err
		}
	}
	return nil
}

// GetAndTouch retrieves a document and simultaneously updates its expiry time
func (h *Handler) GetAndTouch(ctx context.Context, typ, id string, ptr interface{}, ttl uint32) error {
	if err := h.read(ctx, typ, id, ptr, ttl, func(typ, id string, ptr interface{}, ttl uint32) (gocb.Cas, error) {
		documentID, err := h.state.getDocumentKey(typ, id)
		if err != nil {
			return 0, err
		}
		return h.state.bucket.GetAndTouch(documentID, uint32(ttl), ptr)
	}); err != nil {
		return err
	}
	return nil
}

// Ping will ping a list of services and verify they are active and responding in an acceptable period of time
func (h *Handler) Ping(ctx context.Context, services []gocb.ServiceType) (*gocb.PingReport, error) {
	report, err := h.state.bucket.Ping(services)
	if err != nil {
		return nil, err
	}
	return report, nil
}
