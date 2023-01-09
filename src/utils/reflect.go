package utils

import (
	"reflect"
)

func GetFields(s interface{}) []reflect.StructField {
	// Check if s is a pointer
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		panic("GetFields: s is a pointer")
	}

	var fields []reflect.StructField
	for i := 0; i < v.NumField(); i++ {
		fields = append(fields, v.Type().Field(i))
	}
	return fields
}

func GetField(s interface{}, name string) reflect.StructField {
	// Check if s is a pointer
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		panic("GetField: s is a pointer")
	}

	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Name == name {
			return v.Type().Field(i)
		}
	}
	return reflect.StructField{}
}

func SetField(field reflect.Value, value interface{}) {
	field.Set(reflect.ValueOf(value))
}
