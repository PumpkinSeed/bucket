package odatas

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

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
	})
	if err != nil {
		log.Fatal(err)
	}
	return h
}

func getDocumentTypes(ptr interface{}, typs []string, id string) error {
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
			err := getDocumentTypes(structField.Addr().Interface(), typs, id)
			if err != nil {
				return err
			}
			continue
		}
		if inputFieldName == "" {
			inputFieldName = typeField.Name
			if structFieldKind == reflect.Struct {
				err := getDocumentTypes(structField.Addr().Interface(), typs, id)
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
