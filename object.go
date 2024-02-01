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

func CopyField(dst any, src any, fields ...string) {
	dObj := reflect.ValueOf(dst).Elem()
	sObj := reflect.ValueOf(src).Elem()

	for _, item := range fields {
		dField := dObj.FieldByName(item)
		sField := sObj.FieldByName(item)

		TrueF(dField.IsValid() && sField.IsValid(), "invalid field: %s", item)

		dField.Set(sField.FieldByName(item))
	}
}
