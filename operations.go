package odatas

import (
	"errors"
	"fmt"
	"strconv"
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

	_, err := th.state.bucket.Insert(documentID, fields, 0)
	return id, err
}

func (h *Handler) Reed(document, id string, ptr interface{}) error {
	documentID := document + "::" + id

	_, err := th.state.bucket.Get(documentID, ptr)
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
						if err = h.Reed(tag, id, rvQField.Interface()); err != nil {
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

//func (h *Handler) Read(id, t string, ptr interface{}) error {
//	documentID := t + "::" + id
//	var doc interface{}
//
//	_, err := h.bucket.Get(documentID, &doc)
//	if err != nil {
//		return err
//	}
//
//	typ := reflect.TypeOf(ptr).Elem()
//	val := reflect.ValueOf(ptr).Elem()
//
//	if typ.Kind() != reflect.Struct {
//		return errors.New("second argument must be a struct")
//	}
//
//	for i := 0; i < typ.NumField(); i++ {
//		typeField := typ.Field(i)
//		structField := val.Field(i)
//
//		if !structField.CanSet() {
//			fmt.Println("field ", i, "cannot be set")
//			continue
//		}
//
//		structFieldKind := structField.Kind()
//		inputFieldName := strings.Split(typeField.Tag.Get("json"), ",")[0]
//		if structFieldKind == reflect.Struct {
//			err := h.Read(id, inputFieldName, structField.Addr().Interface())
//			if err != nil {
//				return err
//			}
//			continue
//		}
//
//		if inputFieldName == "" {
//			inputFieldName = typeField.Name
//
//			if structFieldKind == reflect.Struct {
//				err := h.Read(id, inputFieldName, structField.Addr().Interface())
//				if err != nil {
//					return err
//				}
//				continue
//			}
//		}
//
//		var inputValue string
//		for _, key := range reflect.ValueOf(doc).MapKeys() {
//			val := reflect.Indirect(key).Interface()
//			if inputFieldName == val {
//				fmt.Println(key)
//				fmt.Println(val)
//				inputValue = fmt.Sprintf("%v", reflect.ValueOf(doc).MapIndex(key).Interface())
//			}
//		}
//
//		if err := setWithProperType(typeField.Type.Kind(), inputValue, structField); err != nil {
//			return err
//		}
//
//	}
//	return nil
//}

func (h *Handler) Remove(id, t string, ptr interface{}) error {
	typs := []string{t}
	e := h.remove(ptr, id, typs)
	if e != nil {
		return e
	}

	for _, typ := range typs {
		_, err := th.state.bucket.Remove(typ+"::"+id, 0)
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

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	case reflect.Ptr:
		switch structField.Type().String() {
		case "*int":
			return setPtrIntField(val, 0, structField)
		case "*int8":
			return setPtrIntField(val, 8, structField)
		case "*int16":
			return setPtrIntField(val, 16, structField)
		case "*int32":
			return setPtrIntField(val, 32, structField)
		case "*int64":
			return setPtrIntField(val, 64, structField)
		case "*uint":
			return setPtrUintField(val, 0, structField)
		case "*uint8":
			return setPtrUintField(val, 8, structField)
		case "*uint16":
			return setPtrUintField(val, 16, structField)
		case "*uint32":
			return setPtrUintField(val, 32, structField)
		case "*uint64":
			return setPtrUintField(val, 64, structField)
		case "*string":
			structField.Set(reflect.ValueOf(&val))
		case "*bool":
			return setPtrBoolField(val, structField)
		case "*float32":
			return setPtrFloatField(val, 32, structField)
		case "*float64":
			return setPtrFloatField(val, 64, structField)
		default:
			return errors.New("unknown type")
		}

	default:
		return errors.New("unknown type")
	}
	return nil
}

func setPtrIntField(value string, bitSize int, field reflect.Value) error {
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.Set(reflect.ValueOf(&intVal))
	}
	return err
}

func setPtrUintField(value string, bitSize int, field reflect.Value) error {
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.Set(reflect.ValueOf(&uintVal))
	}
	return err
}

func setPtrFloatField(value string, bitsize int, field reflect.Value) error {
	floatVal, err := strconv.ParseFloat(value, bitsize)
	if err == nil {
		field.Set(reflect.ValueOf(&floatVal))
	}
	return err
}

func setPtrBoolField(value string, field reflect.Value) error {
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		field.Set(reflect.ValueOf(&boolVal))
	}
	return err
}

func setIntField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(value string, field reflect.Value) error {
	if value == "" {
		value = "false"
	}
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0.0"
	}
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}
