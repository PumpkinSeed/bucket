package odatas

import (
	"errors"
	"fmt"
	"github.com/couchbase/gocb"
	"github.com/rs/xid"
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
				err := write(b, a.Tag.Get("json"), id)
				if err != nil {
					return err
				}
			} else {
				k := v.Type().Field(i).Tag.Get("json")
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

func Read(id, documentType string, result interface{}) {

}
