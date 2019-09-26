package odatas

import (
	"errors"
	"fmt"
	"strings"

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
				if rvQField.Kind() == reflect.Struct {
					if tag, ok := rtQField.Tag.Lookup(tagJson); ok {
						if _, err := h.write(rvQField.Interface(), removeOmitempty(tag), id); err != nil {
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

	_, err := h.state.bucket.Insert(documentID, fields, 0)
	return id, err
}

func removeOmitempty(tag string) string {
	if strings.Contains(tag, ",omitempty") {
		tag = strings.Replace(tag, ",omitempty", "", -1)
	}
	return tag
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
					fmt.Println(reflect.Indirect(rvQField.Elem()))
					if reflect.Indirect(rvQField).Kind() != reflect.Struct {
						continue
					}
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
	e := h.remove(ptr, id, typs)
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

func (h *Handler) remove(ptr interface{}, id string, typs []string) error {
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
			err := h.remove(structField.Addr().Interface(), id, typs)
			if err != nil {
				return err
			}
			continue
		}

		if inputFieldName == "" {
			inputFieldName = typeField.Name

			if structFieldKind == reflect.Struct {
				err := h.remove(structField.Addr().Interface(), id, typs)
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
