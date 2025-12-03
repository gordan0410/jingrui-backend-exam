package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func TrimAllStrings(a any) {
	aAddress := reflect.ValueOf(a)
	if k := aAddress.Kind(); k != reflect.Ptr {
		panic("a is not a pointer")
	}
	if v := aAddress.IsNil(); v {
		panic("a is nil")
	}
	// pointer pointer 處理
	for aAddress.Elem().Kind() == reflect.Ptr {
		aAddress = aAddress.Elem()
	}

	aValue := aAddress.Elem()
	if aValue.Kind() != reflect.Struct {
		panic("a is not pointing to a struct")
	}

	for aValue.IsValid() {
		var tempValue reflect.Value
		fieldIdx := aValue.NumField()
		for i := 0; i < fieldIdx; i++ {
			if aValue.Field(i).Kind() == reflect.String {
				trimmed := strings.TrimSpace(aValue.Field(i).String())
				aValue.Field(i).SetString(trimmed)
			}
			if aValue.Field(i).Kind() == reflect.Ptr && aValue.Field(i).Elem().Kind() == reflect.Struct {
				tempValue = aValue.Field(i).Elem()
			}
		}
		aValue = tempValue
	}
}

func main() {
	type Person struct {
		Name string
		Age  int
		Next *Person
	}

	a := &Person{
		Name: " name ",
		Age:  20,
		Next: &Person{
			Name: " name2 ",
			Age:  21,
			Next: &Person{
				Name: " name3 ",
				Age:  22,
			},
		},
	}

	// 不確定為什麼要放pointer, pointer進去
	TrimAllStrings(&a)

	m, _ := json.Marshal(a)

	fmt.Println(string(m))
	//
	//a.Next = a
	//
	//TrimAllStrings(&a)

	// 不太知道為何 但應該是name3?
	//fmt.Println(a.Next.Next.Name == "name")
	fmt.Println(a.Next.Next.Name == "name3")
}
