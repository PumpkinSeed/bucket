package odatas

import (
	"errors"
	"fmt"
	"strings"

	"github.com/couchbase/gocb"

	"github.com/rs/xid"

	"reflect"
)

func (h *Handler) Write(q interface{}, typ string) (string, error) {
	documentID, err := h.write(q, typ, "")
	if err != nil {
		return "", err
	}
	return documentID, nil
}

func (h *Handler) write(q interface{}, typ, id string) (string, error) {
	fields := make(map[string]interface{})
	if id == "" {
		id = xid.New().String()
	}
	documentID := typ + "::" + id

	rvQ := reflect.ValueOf(q)
	rtQ := rvQ.Type()

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
					if _, err := h.write(rvQField.Interface(), tag, id); err != nil {
						return id, err
					}
				}
			} else {
				if tag, ok := rtQField.Tag.Lookup(tagJson); ok {
					fields[tag] = rvQField.Interface()
				}
			}
		}
	} else {
		return id, errors.New("not a struct")
	}

	_, err := h.state.bucket.Insert(documentID, fields, 0)
	return id, err
}

func (h *Handler) Read(document, id string, ptr interface{}) error {
	documentID := document + "::" + id

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
						if err = h.Read(tag, id, rvQField.Interface()); err != nil {
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

func (h *Handler) Remove(id, t string, ptr interface{}) error {
	typs := []string{t}
	if err := h.getDocumentTypes(ptr, id, typs); err != nil {
		return err
	}

	for _, typ := range typs {
		_, err := h.state.bucket.Remove(typ+"::"+id, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) getDocumentTypes(ptr interface{}, id string, typs []string) error {
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()
	if typ.Kind() != reflect.Struct {
		return errors.New("second argument must be a struct")
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
			err := h.getDocumentTypes(structField.Addr().Interface(), id, typs)
			if err != nil {
				return err
			}
			continue
		}

		if inputFieldName == "" {
			inputFieldName = typeField.Name

			if structFieldKind == reflect.Struct {
				err := h.getDocumentTypes(structField.Addr().Interface(), id, typs)
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

func (h *Handler) Upsert(id, t string, ptr interface{}, ttl int) error {
	fields := make(map[string]interface{})
	documentID := t + "::" + id

	if reflect.ValueOf(ptr).Kind() == reflect.Struct {
		v := reflect.ValueOf(ptr)
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.Kind() == reflect.Struct {
				a := v.Type().Field(i)
				b := field.Interface()
				fieldName := strings.Split(a.Tag.Get("json"), ",")[0]
				if err := h.Upsert(id, fieldName, b, ttl); err != nil {
					return err
				}
			} else {
				k := strings.Split(v.Type().Field(i).Tag.Get("json"), ",")[0]
				fields[k] = fmt.Sprintf("%v", field)
			}
		}
	} else {
		return errors.New("not a struct")
	}

	_, err := h.state.bucket.Upsert(documentID, fields, uint32(ttl))
	return err
}

func (h *Handler) Touch(id, t string, ptr interface{}, ttl int) error {
	types := []string{t}
	e := h.getDocumentTypes(ptr, id, types)
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

func (h *Handler) Ping() (*gocb.PingReport, error) {
	report, err := h.state.bucket.Ping(nil)
	if err != nil {
		return nil, err
	}
	return report, nil
}
