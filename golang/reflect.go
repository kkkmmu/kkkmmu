package main

import (
	"fmt"
	"reflect"
)

type MT int

func (mt *MT) Method1(name string) error {
	fmt.Println("Call method1 with ", name)
	return nil
}

func (mt *MT) Method2(name string, id string) error {
	fmt.Println("Call method2 with ", name, id)
	return nil
}

func (mt *MT) Call(name string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
		fmt.Println(inputs[i])
	}

	reflect.ValueOf(mt).MethodByName(name).Call(inputs)
}

func main() {
	var s *MT
	val := reflect.TypeOf(s)

	for i := 0; i < val.NumMethod(); i++ {
		method := val.Method(i)
		fmt.Println(method.Name)
	}

	s.Call("Method1", "t1111111")
	s.Call("Method2", "t1111111", "t22222222")
}
