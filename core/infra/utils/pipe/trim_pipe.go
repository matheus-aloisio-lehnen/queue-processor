package pipe

import (
	"reflect"
	"strings"
)

func TrimPipe(dto any) {
	v := reflect.ValueOf(dto).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.String && f.CanSet() {
			f.SetString(strings.TrimSpace(f.String()))
		}
	}
}
