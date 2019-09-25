package odatas

import (
	"errors"
	"fmt"
	"github.com/rs/xid"
	"gopkg.in/couchbase/gocb.v1"
	"strconv"
	"strings"

	"reflect"
)

var (
	placeholderBucket *gocb.Bucket
)

func placeholderInit() {
	if placeholderBucket == nil {
		cb, err := gocb.Connect("couchbase://localhost")
		if err != nil {
			panic(err)
		}

		err = cb.Authenticate(gocb.PasswordAuthenticator{
			Username: "Administrator",
			Password: "password",
		})
		if err != nil {
			panic(err)
		}

		placeholderBucket, err = cb.OpenBucket("company", "")
		if err != nil {
			panic(err)
		}
	}
}

func Insert(q interface{}, typ string) error {
	err := write(q, typ, "")
	if err != nil {
		return err
	}
	return nil
}

func write(q interface{}, typ, id string) error {
	fields := make(map[string]interface{})
	if id == "" {
		id = xid.New().String()
	}
	documentID := typ + "::" + id

	//var jso []byte
	if reflect.ValueOf(q).Kind() == reflect.Struct {
		v := reflect.ValueOf(q)
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Kind() == reflect.Struct {
				a := v.Type().Field(i)
				b := v.Field(i).Interface()
				fieldName := strings.Split(a.Tag.Get("json"),",")[0]
				err := write(b, fieldName, id)
				if err != nil {
					return err
				}
			} else {
				k := strings.Split(v.Type().Field(i).Tag.Get("json"),",")[0]
				val := v.Field(i)
				fields[k] = fmt.Sprintf("%v", val)
			}
		}

		//jso, err := json.Marshal(fields)
		//if err != nil {
		//	return err
		//}
		//fmt.Println(string(jso))
	} else {
		return errors.New("not a struct")
	}
	_, err := placeholderBucket.Insert(documentID, fields, 0)
	return err
}


func read(id,t string, ptr interface{}) error {
	documentID := t + "::" + id
	var doc interface{}

	_, err := placeholderBucket.Get(documentID, &doc)
	if err != nil {
		return err
	}

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
		//inputFieldName := typeField.Tag.Get("json")
		inputFieldName := strings.Split(typeField.Tag.Get("json"),",")[0]
		//inputFieldName := typeField.Name
		if structFieldKind == reflect.Struct {
			err := read(id,inputFieldName, structField.Addr().Interface())
			if err != nil {
				return err
			}
			continue
		}

		if inputFieldName == "" {
			inputFieldName = typeField.Name

			if structFieldKind == reflect.Struct {
				err := read( id,inputFieldName, structField.Addr().Interface())
				if err != nil {
					return err
				}
				continue
			}
		}

		var inputValue string
		 for _,key := range reflect.ValueOf(doc).MapKeys() {
			 val := reflect.Indirect(key).Interface()
		 	if inputFieldName == val {
				inputValue = fmt.Sprintf("%v",reflect.ValueOf(doc).MapIndex(key).Interface())
			}
		 }

			if err := setWithProperType(typeField.Type.Kind(), inputValue, structField); err != nil {
				return err
			}

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
		case "*int64":
			return setPtrIntField(val, 64, structField)
		case "*string":
			structField.Set(reflect.ValueOf(&val))
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
