package bucket

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type NullTimeout struct {
	valid bool
	Value time.Duration
}

func NullTimeoutMillisec(dur uint64) NullTimeout {
	return NullTimeout{
		valid: true,
		Value: time.Duration(dur) * time.Millisecond,
	}
}

func NullTimeoutSec(dur uint64) NullTimeout {
	return NullTimeout{
		valid: true,
		Value: time.Duration(dur) * time.Second,
	}
}

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

func getDocumentTypes(ptr interface{}, typs []string) error {
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
			err := getDocumentTypes(structField.Addr().Interface(), typs)
			if err != nil {
				return err
			}
			continue
		}
		if inputFieldName == "" {
			inputFieldName = typeField.Name
			if structFieldKind == reflect.Struct {
				err := getDocumentTypes(structField.Addr().Interface(), typs)
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
