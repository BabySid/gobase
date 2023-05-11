package gobase

import "reflect"

func GetShortNameOfPointer(i any) string {
	True(reflect.TypeOf(i).Kind() == reflect.Pointer)
	return reflect.TypeOf(i).Elem().Name()
}

func GetFullNameOfPointer(i any) string {
	True(reflect.TypeOf(i).Kind() == reflect.Pointer)
	return reflect.TypeOf(i).Elem().String()
}
