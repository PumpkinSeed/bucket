package bucket

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// NullTimeout is the library's built in NullTimeout type for Opts
type NullTimeout struct {
	valid bool
	Value time.Duration
}

// NullTimeoutMillisec creates a NullTimeout with the given millisec
func NullTimeoutMillisec(dur uint64) NullTimeout {
	return NullTimeout{
		valid: true,
		Value: time.Duration(dur) * time.Millisecond,
	}
}

// NullTimeoutSec creates a NullTimeout with the given sec
func NullTimeoutSec(dur uint64) NullTimeout {
	return NullTimeout{
		valid: true,
		Value: time.Duration(dur) * time.Second,
	}
}

// NullTimeoutFrom creates a NullTimeout with the given time.Duration
func NullTimeoutFrom(dur time.Duration) NullTimeout {
	return NullTimeout{
		valid: true,
		Value: dur,
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func setupBasicAuth(req *http.Request) {
	req.Header.Add("Authorization", "Basic "+basicAuth("Administrator", "password"))
}

func defaultHandler() *Handler {
	h, err := New(&Configuration{
		Username:       "Administrator",
		Password:       "password",
		BucketName:     bucketName,
		BucketPassword: "",
		Separator:      "::",
	})
	if err != nil {
		log.Fatal(err)
	}
	return h
}

func removeOmitempty(tag string) string {
	if strings.Contains(tag, ",omitempty") {
		tag = strings.Replace(tag, ",omitempty", "", -1)
	}
	return tag
}

func getDocumentTypes(value interface{}) ([]string, error) {
	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		return nil, ErrInvalidGetDocumentTypesParam
	}

	var typs []string
	var val = reflect.New(reflect.ValueOf(value).Type().Elem())
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("value argument must be a struct")
	}
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if val, ok := typeField.Tag.Lookup(tagReferenced); ok {
			typs = append(typs, val)
			if structField.IsNil() && structField.CanSet() {
				structField.Set(reflect.New(structField.Type().Elem()))
			}
			moreTypes, err := getDocumentTypes(structField.Interface())
			if err != nil {
				return nil, err
			}
			typs = append(typs, moreTypes...)
		}
	}
	return typs, nil
}

func getStructAddressableSubfields(value reflect.Value) map[string]interface{} {
	if value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
	}

	var result = make(map[string]interface{})
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		if tag, ok := typ.Field(i).Tag.Lookup(tagReferenced); ok && tag != "" {
			result[tag] = value.Field(i).Addr().Interface()
		}
	}

	return result
}
