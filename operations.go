package main

import (
	"fmt"
	"reflect"
)

type order struct {
	ordId      int
	customerId int
}

func Write(s interface{}) {
	if reflect.ValueOf(s).Kind() == reflect.Struct {
		v := reflect.ValueOf(s)
		for i:= 0; i < v.NumField(); i++ {
			fmt.Println("field type",v.Type())
			fmt.Println(v.Field(i))
		}
	}
}

func main() {

}

func Read() {

}
