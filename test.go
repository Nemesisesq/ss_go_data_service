package main

//import (
//	"gopkg.in/oleiade/reflections.v1"
//	"fmt"
//	"reflect"
//)
//
//func main() {
//	fmt.Println("Hello, playground")
//
//	type MyStruct struct {
//		FirstField  string
//		SecondField int
//		ThirdField  string
//	}
//
//	s := MyStruct{
//		FirstField:  "first value",
//		SecondField: 2,
//		ThirdField:  "third value",
//	}
//
//	fieldsToExtract := []string{"FirstField", "ThirdField"}
//
//	for _, fieldName := range fieldsToExtract {
//		value, _ := reflections.GetField(s, fieldName)
//
//		fmt.Println(value)
//	}
//
//	data := []string{"one", "two", "three"}
//	test(data)
//	moredata := []int{1, 2, 3}
//	test(moredata)
//
//}
//
//func test(t interface{}) {
//	switch reflect.TypeOf(t).Kind() {
//	case reflect.Slice:
//		s := reflect.ValueOf(t)
//
//		for i := 0; i < s.Len(); i++ {
//			fmt.Println(s.Index(i))
//		}
//	}
//}