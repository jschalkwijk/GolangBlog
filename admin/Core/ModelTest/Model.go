package ModelTest

import (
	"reflect"
	"fmt"
	"strconv"
)


func GetAll(model interface{})  {
	myType := reflect.TypeOf(model)
	fmt.Println(myType)
	slice := reflect.MakeSlice(reflect.SliceOf(myType), 0, 0)

	val := reflect.ValueOf(myType).Elem()
	fmt.Println(val)

	m := make(map[string]string)



	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		fmt.Println(typeField)
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			m[typeField.Name] = strconv.FormatInt(val.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			m[typeField.Name] = strconv.FormatUint(val.Uint(), 10)
		case reflect.String:
			m[typeField.Name] = val.String()
		}
	}

	fmt.Println(m)
	// Create a pointer to a slice value and set it to the slice
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)
	fmt.Println(slice)
	fmt.Println(x.Interface())
}

//func GetAll(m interface{}) (map[string]reflect.Type) {
//	typ := reflect.TypeOf(m)
//	// if a pointer to a struct is passed, get the type of the dereferenced object
//	if typ.Kind() == reflect.Ptr{
//		typ = typ.Elem()
//	}
//
//	// create an attribute data structure as a map of types keyed by a string.
//	attrs := make(map[string]reflect.Type)
//	// Only structs are supported so return an empty result if the passed object
//	// isn't a struct
//	if typ.Kind() != reflect.Struct {
//		fmt.Printf("%v type can't have attributes inspected\n", typ.Kind())
//		return attrs
//	}
//
//	// loop through the struct's fields and set the map
//	for i := 0; i < typ.NumField(); i++ {
//		p := typ.Field(i)
//		if !p.Anonymous {
//			attrs[p.Name] = p.Type
//		}
//	}
//
//	return attrs
//}

