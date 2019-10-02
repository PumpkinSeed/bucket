package bucket

import (
	"context"
	"reflect"

	"github.com/couchbase/gocb"
	"github.com/rs/xid"
)

type writerF func(string, string, interface{}, int) (gocb.Cas, error)
type readerF func(string, string, interface{}, int) (gocb.Cas, error)
type Cas map[string]gocb.Cas

func (h *Handler) Insert(ctx context.Context, typ, id string, q interface{}) (Cas, string, error) {
	cas := make(map[string]gocb.Cas)
	if id == "" {
		id = xid.New().String()
	}
	id, err := h.write(ctx, typ, id, q, func(typ, id string, ptr interface{}, ttl int) (gocb.Cas, error) {
		documentID := typ + "::" + id
		return h.state.bucket.Insert(documentID, ptr, 0)
	}, cas)
	if err != nil {
		return nil, "", err
	}
	return cas, id, nil
}

func (h *Handler) write(ctx context.Context, typ, id string, q interface{}, f writerF, cas Cas) (string, error) {
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
				if tag, ok := rtQField.Tag.Lookup(tagJson); ok {
					fields[removeOmitempty(tag)] = nil
				}
			} else {
				if rvQField.Kind() == reflect.Ptr {
					rvQField = reflect.Indirect(rvQField)
				}

				if rvQField.Kind() == reflect.Struct && hasRefTag {
					if refTag == "" {
						return "", ErrEmptyRefTag
					}
					if _, err := h.write(ctx, refTag, id, rvQField.Interface(), f, cas); err != nil {
						return id, err
					}
				} else {
					if tag, ok := rtQField.Tag.Lookup(tagJson); ok {
						fields[removeOmitempty(tag)] = rvQField.Interface()
					}
				}
			}
		}
	}
	c, err := f(typ, id, fields, -1)
	cas[typ] = c

	return id, err
}

func (h *Handler) Get(ctx context.Context, typ, id string, ptr interface{}) error {
	if err := h.read(ctx, typ, id, ptr, -1, func(typ, id string, ptr interface{}, ttl int) (gocb.Cas, error) {
		documentID := typ + "::" + id
		return h.state.bucket.Get(documentID, ptr)
	}); err != nil {
		return err
	}
	return nil
}

func (h *Handler) read(ctx context.Context, typ, id string, ptr interface{}, ttl int, f readerF) error {
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
					if err = h.read(ctx, refTag, id, rvQField.Interface(), ttl, f); err != nil {
						return err
					}
				}
			}
		}
	} else {
		return ErrInputStructPointer
	}
	return nil
}

func (h *Handler) Remove(ctx context.Context, typ, id string, ptr interface{}) error {
	typs := []string{typ}
	e := getDocumentTypes(ptr, typs, id)
	if e != nil {
		return e
	}

	for _, typ := range typs {
		_, err := h.state.bucket.Remove(typ+"::"+id, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) Upsert(ctx context.Context, typ, id string, q interface{}, ttl uint32) (Cas, string, error) {
	cas := make(map[string]gocb.Cas)
	if id == "" {
		id = xid.New().String()
	}
	id, err := h.write(ctx, typ, id, q, func(typ, id string, q interface{}, ttl int) (gocb.Cas, error) {
		documentID := typ + "::" + id
		return h.state.bucket.Upsert(documentID, q, uint32(ttl))
	}, cas)
	if err != nil {
		return nil, "", err
	}
	return cas, id, nil

}

func (h *Handler) Touch(ctx context.Context, typ, id string, ptr interface{}, ttl int) error {
	types := []string{typ}
	e := getDocumentTypes(ptr, types, id)
	if e != nil {
		return e
	}

	for _, typ := range types {
		_, err := h.state.bucket.Touch(typ+"::"+id, 0, uint32(ttl))
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) GetAndTouch(ctx context.Context, typ, id string, ptr interface{}, ttl int) error {
	if err := h.read(ctx, typ, id, ptr, ttl, func(typ, id string, ptr interface{}, ttl int) (gocb.Cas, error) {
		documentID := typ + "::" + id
		return h.state.bucket.GetAndTouch(documentID, uint32(ttl), ptr)
	}); err != nil {
		return err
	}
	return nil
}

func (h *Handler) Ping(ctx context.Context, services []gocb.ServiceType) (*gocb.PingReport, error) {
	report, err := h.state.bucket.Ping(services)
	if err != nil {
		return nil, err
	}
	return report, nil
}
