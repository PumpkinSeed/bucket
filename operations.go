package bucket

import (
	"context"
	"reflect"

	"github.com/couchbase/gocb"
	"github.com/rs/xid"
)

type writerF func(string, string, interface{}, uint32) (gocb.Cas, error)
type readerF func(string, string, interface{}, uint32) (gocb.Cas, error)

func (h *Handler) Insert(ctx context.Context, typ, id string, q interface{}, ttl uint32) (string, error) {
	if id == "" {
		id = xid.New().String()
	}
	id, err := h.write(ctx, typ, id, q, func(typ, id string, ptr interface{}, ttl uint32) (gocb.Cas, error) {
		documentID := typ + "::" + id
		return h.state.bucket.Insert(documentID, ptr, ttl)
	}, ttl)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *Handler) write(ctx context.Context, typ, id string, q interface{}, f writerF, ttl uint32) (string, error) {
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
			refTag, hasRefTag := rtQField.Tag.Lookup("referenced")

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
					if _, err := h.write(ctx, refTag, id, rvQField.Interface(), f, ttl); err != nil {
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
	_, err := f(typ, id, fields, ttl)
	return id, err
}

func (h *Handler) Get(ctx context.Context, typ, id string, ptr interface{}) error {
	if err := h.read(ctx, typ, id, ptr, 0, func(typ, id string, ptr interface{}, ttl uint32) (gocb.Cas, error) {
		documentID := typ + "::" + id
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
					refTag, hasRefTag := rtQField.Tag.Lookup("referenced")
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

func (h *Handler) Upsert(ctx context.Context, typ, id string, q interface{}, ttl uint32) (string, error) {
	if id == "" {
		id = xid.New().String()
	}
	id, err := h.write(ctx, typ, id, q, func(typ, id string, q interface{}, ttl uint32) (gocb.Cas, error) {
		documentID := typ + "::" + id
		return h.state.bucket.Upsert(documentID, q, ttl)
	}, ttl)
	if err != nil {
		return "", err
	}
	return id, nil

}

func (h *Handler) Touch(ctx context.Context, typ, id string, ptr interface{}, ttl uint32) error {
	types := []string{typ}
	e := getDocumentTypes(ptr, types, id)
	if e != nil {
		return e
	}

	for _, typ := range types {
		_, err := h.state.bucket.Touch(typ+"::"+id, 0, ttl)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) GetAndTouch(ctx context.Context, typ, id string, ptr interface{}, ttl uint32) error {
	if err := h.read(ctx, typ, id, ptr, ttl, func(typ, id string, ptr interface{}, ttl uint32) (gocb.Cas, error) {
		documentID := typ + "::" + id
		return h.state.bucket.GetAndTouch(documentID, ttl, ptr)
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
