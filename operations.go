package odatas

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/xid"
	"reflect"
)

func Write(i interface{},typ string) error {
	fields := make(map[string]interface{})
	fields["ID"]= typ + "::" + xid.New().String()
	if reflect.ValueOf(i).Kind() == reflect.Struct {
		v := reflect.ValueOf(i)
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Kind() == reflect.Struct {
				//TODO
				continue
			}
			k := v.Type().Field(i).Name
			val := v.Field(i)
			fields[k] = fmt.Sprintf("%v",val)
		}

		json, err := json.Marshal(fields)
		if err != nil {
			return err
		}
		fmt.Println(string(json))
	} else {
		return errors.New("not a struct")
	}
	fmt.Println("dass")
	fmt.Println(fields["OrdId"])
	fmt.Println(fields["CustomerId"])
	fmt.Println(len(fields))

	//handler.bucket.Insert(map["ID"],json)
	return nil
}



func Read() {

}
