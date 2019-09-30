package odatas

import (
	"context"
	"fmt"
	"strings"

	"github.com/couchbase/gocb"

	"github.com/rs/xid"

	"reflect"
)

func (h *Handler) Insert(ctx context.Context, typ string, q interface{}) (string, error) {
	id, err := h.write(ctx, typ, xid.New().String(), q, func(typ, id string, ptr interface{}, ttl int) (gocb.Cas, error) {
		documentID := typ + "::" + id
		return h.state.bucket.Insert(documentID, ptr, 0)
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *Handler) write(ctx context.Context, typ, id string, q interface{}, f func(string, string, interface{}, int) (gocb.Cas, error)) (string, error) {
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

			if rvQField.Kind() == reflect.Ptr && rvQField.IsNil() {
				if tag, ok := rtQField.Tag.Lookup(tagJson); ok {
					fields[removeOmitempty(tag)] = nil
				}
			} else {
				if rvQField.Kind() == reflect.Ptr {
					rvQField = reflect.Indirect(rvQField)
				}
				_, hasReferenceTag := rtQField.Tag.Lookup("referenced")

				if rvQField.Kind() == reflect.Struct && hasReferenceTag {
					if tag, ok := rtQField.Tag.Lookup(tagJson); ok {
						if _, err := h.write(ctx, removeOmitempty(tag), id, rvQField.Interface(), f); err != nil {
							return id, err
						}
					}
				} else {
					if tag, ok := rtQField.Tag.Lookup(tagJson); ok {
						fields[removeOmitempty(tag)] = rvQField.Interface()
					}
				}
			}
		}
	}
	_, err := f(typ, id, q, -1)
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

func (h *Handler) read(ctx context.Context, typ, id string, ptr interface{}, ttl int, f func(string, string, interface{}, int) (gocb.Cas, error)) error {

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
					if _, hasReferencedTag := rtQField.Tag.Lookup("referenced"); !hasReferencedTag || rvQField.Type().Elem().Kind() != reflect.Struct {
						continue
					}
					rvQField.Set(reflect.New(rvQField.Type().Elem()))
					if tag, ok := rtQField.Tag.Lookup(tagJson); ok {
						if strings.Contains(tag, ",omitempty") {
							tag = strings.Replace(tag, ",omitempty", "", -1)
						}
						if err = h.Get(ctx, tag, id, rvQField.Interface()); err != nil {
							return err
						}
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
	e := h.remove(ctx, typs, ptr, id)
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

func (h *Handler) remove(ctx context.Context, typs []string, ptr interface{}, id string) error {
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()
	if typ.Kind() != reflect.Struct {
		return ErrFirstParameterNotStruct
	}
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)

		if !structField.CanSet() {
			fmt.Println("field ", i, "cannot be set")
			continue
		}

		structFieldKind := structField.Kind()
		inputFieldName := strings.Split(typeField.Tag.Get("json"), ",")[0]
		if structFieldKind == reflect.Struct {
			err := h.remove(ctx, typs, structField.Addr().Interface(), id)
			if err != nil {
				return err
			}
			continue
		}

		if inputFieldName == "" {
			inputFieldName = typeField.Name

			if structFieldKind == reflect.Struct {
				err := h.remove(ctx, typs, structField.Addr().Interface(), id)
				if err != nil {
					return err
				}
				continue
			}
		}
		typs = append(typs, inputFieldName)
	}
	return nil
}

func (h *Handler) Upsert(ctx context.Context, typ, id string, q interface{}, ttl uint32) (string, error) {
	if id == "" {
		id = xid.New().String()
	}
	id, err := h.write(ctx, typ, id, q, func(typ, id string, q interface{}, ttl int) (gocb.Cas, error) {
		documentID := typ + "::" + id
		return h.state.bucket.Upsert(documentID, q, uint32(ttl))
	})
	if err != nil {
		return "", err
	}
	return id, nil

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
