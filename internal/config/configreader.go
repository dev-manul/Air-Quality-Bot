package config

import (
	"github.com/jinzhu/configor"
	"log"
	"reflect"
)

var ConfigDebug = false
var ConfigVerbose = false

func Provide[T any](config T, files ...string) func() (T, error) {
	t := reflect.TypeOf(config)
	v := reflect.New(t)
	initializeStruct(t, v.Elem())
	c, _ := v.Interface().(*T)
	config = *c

	return func() (T, error) {
		err := configor.
			New(&configor.Config{Debug: ConfigDebug, ErrorOnUnmatchedKeys: true, Verbose: ConfigVerbose}).
			Load(&config, files...)
		if err != nil {
			log.Panic(err)
		}
		return config, err
	}
}

func initializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			initializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			initializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		default:
		}
	}
}
