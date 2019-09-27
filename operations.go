package odatas

import (
	"context"
	"fmt"
	"strings"

	"github.com/couchbase/gocb"

	"github.com/rs/xid"

	"reflect"
)

func (h *Handler) Write(ctx context.Context, typ string, q interface{}) (string, error) {
	id, err := h.write(ctx, typ, "", q)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *Handler) write(ctx context.Context, typ, id string, q interface{}) (string, error) {
	if !h.state.inspect(typ) {
		err := h.state.setType(typ, typ)
		if err != nil {
			return "", err
		}
	}
	fields := make(map[string]interface{})
	if id == "" {
		id = xid.New().String()
	}
	documentID := typ + "::" + id

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

			if rvQField.Kind() == reflect.Ptr {
				rvQField = reflect.Indirect(rvQField)
			}
			if rvQField.Kind() == reflect.Struct {
				if tag, ok := rtQField.Tag.Lookup(tagJson); ok {
					if strings.Contains(tag, ",omitempty") {
						tag = strings.Replace(tag, ",omitempty", "", -1)
					}
					if _, err := h.write(ctx, tag, id, rvQField.Interface()); err != nil {
						return id, err
					}
				}
			} else {
				if tag, ok := rtQField.Tag.Lookup(tagJson); ok {
					fields[tag] = rvQField.Interface()
				}
			}
		}
	}

	_, err := h.state.bucket.Insert(documentID, fields, 0)
	return id, err
}

func (h *Handler) Read(ctx context.Context, typ, id string, ptr interface{}) error {
	documentID := typ + "::" + id

	_, err := h.state.bucket.Get(documentID, ptr)
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
					rvQField.Set(reflect.New(rvQField.Type().Elem()))
					if tag, ok := rtQField.Tag.Lookup(tagJson); ok {
						if strings.Contains(tag, ",omitempty") {
							tag = strings.Replace(tag, ",omitempty", "", -1)
						}
						if err = h.Read(ctx, tag, id, rvQField.Interface()); err != nil {
							return err
						}
					}
				}
			}
		}
	} else {
		// err should be pointer
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

func (h *Handler) Upsert(ctx context.Context, typ, id string, ptr interface{}, ttl uint32) error {
	if !h.state.inspect(typ) {
		err := h.state.setType(typ, typ)
		if err != nil {
			return err
		}
	}
	fields := make(map[string]interface{})
	if id == "" {
		id = xid.New().String()
	}
	documentID := typ + "::" + id

	rvQ := reflect.ValueOf(ptr)
	rtQ := rvQ.Type()

	if rtQ.Kind() == reflect.Ptr {
		rvQ = reflect.Indirect(rvQ)
		rtQ = rvQ.Type()
	}

	if rtQ.Kind() == reflect.Struct {
		for i := 0; i < rvQ.NumField(); i++ {
			rvQField := rvQ.Field(i)
			rtQField := rtQ.Field(i)

			if rvQField.Kind() == reflect.Ptr {
				rvQField = reflect.Indirect(rvQField)
			}
			if rvQField.Kind() == reflect.Struct {
				if tag, ok := rtQField.Tag.Lookup(tagJson); ok {
					if strings.Contains(tag, ",omitempty") {
						tag = strings.Replace(tag, ",omitempty", "", -1)
					}
					if err := h.Upsert(ctx, tag, id, rvQField.Interface(), ttl); err != nil {
						return err
					}
				}
			} else {
				if tag, ok := rtQField.Tag.Lookup(tagJson); ok {
					fields[tag] = rvQField.Interface()
				}
			}
		}
	}

	_, err := h.state.bucket.Upsert(documentID, fields, ttl)
	return err
}

func (h *Handler) Touch(typ, id string, ptr interface{}, ttl int) error {
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

func GetAndTouch() error {
	return nil
}

func (h *Handler) Ping(ctx context.Context, services []gocb.ServiceType) (*gocb.PingReport, error) {
	report, err := h.state.bucket.Ping(services)
	if err != nil {
		return nil, err
	}
	return report, nil
}
