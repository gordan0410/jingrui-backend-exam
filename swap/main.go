package main

import (
	"fmt"
	"reflect"
)

func swap[T any](a, b T) {
	aAddress := reflect.ValueOf(a)
	if k := aAddress.Kind(); k != reflect.Ptr && k != reflect.Interface {
		if k == reflect.Invalid {
			panic("a is nil")
		}
		panic("a is not a pointer/interface")
	}
	aValue := aAddress.Elem()
	bAddress := reflect.ValueOf(b)
	bValue := bAddress.Elem()

	if aAddress.IsNil() || bAddress.IsNil() {
		panic("a or b is nil")
	}

	if !aValue.CanSet() || !bValue.CanSet() {
		panic("a or b cannot set")
	}

	temp := reflect.New(aValue.Type())
	temp.Elem().Set(aValue)
	aValue.Set(bValue)
	bValue.Set(temp.Elem())
}

func main() {
	a := 10
	b := 20

	fmt.Printf("a = %d, &a = %p\n", a, &a)
	fmt.Printf("b = %d, &b = %p\n", b, &b)

	swap(&a, &b)

	fmt.Printf("a = %d, &a = %p\n", a, &a)
	fmt.Printf("b = %d, &b = %p\n", b, &b)
}
