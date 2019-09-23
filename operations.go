package odatas

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/xid"
	"reflect"
)

func Insert(q interface{},typ string) error {
	err := write(q,typ,"")
	if err != nil {
		return err
	}
	return nil
}

func write(q interface{},typ,id string) error {
	fields := make(map[string]interface{})
	if id == "" {
		id = xid.New().String()
	}
	fields["ID"]= typ + "::" + id

	if reflect.ValueOf(q).Kind() == reflect.Struct {
		v := reflect.ValueOf(q)
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Kind() == reflect.Struct {
				a := v.Type().Field(i)
				b := v.Field(i).Interface()
				err := write(b,a.Name,id)
				if err != nil {
					return err
				}
			} else {
				k := v.Type().Field(i).Name
				val := v.Field(i)
				fields[k] = fmt.Sprintf("%v", val)
			}
		}

		json, err := json.Marshal(fields)
		if err != nil {
			return err
		}
		fmt.Println(string(json))
	} else {
		return errors.New("not a struct")
	}
	//handler.bucket.Insert(map["ID"],json)
	return nil
}

func Read() {

}

