package tachikoma

import "reflect"

type Selector func(interface{}) interface{}

func InitialSelector() Selector {
	return func(o interface{}) interface{} {
		return o
	}
}

func IndexSelector(i int, s Selector) Selector {
	return func(o interface{}) interface{} {
		o = s(o)
		d := reflect.ValueOf(o)
		t := reflect.TypeOf(o)

		if t.Kind() == reflect.Ptr {
			d = d.Elem()
			t = t.Elem()
		}

		return d.Field(i).Interface()
	}
}
