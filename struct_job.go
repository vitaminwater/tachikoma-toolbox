package tachikoma

import (
	"errors"
	"reflect"
)

type StructJobGenerator func(reflect.StructField, Selector) []Job

func JobsFromStruct(d interface{}, g StructJobGenerator) ([]Job, error) {
	jobs := make([]Job, 0)
	if d == nil {
		return jobs, errors.New("d is nil.")
	}

	var traverse func(s Selector) error
	traverse = func(s Selector) error {
		d := s(d)
		v := reflect.ValueOf(d)
		t := reflect.TypeOf(d)

		if t.Kind() == reflect.Ptr {
			v = v.Elem()
			t = t.Elem()
		}

		if t.Kind() != reflect.Struct {
			return errors.New("data should be struct or pointer to struct")
		}

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			k := f.Type.Kind()
			if k == reflect.Struct {
				if err := traverse(IndexSelector(i, s)); err != nil {
					return err
				}
				continue
			} else if k == reflect.Array {
				continue
			}
			j := g(f, IndexSelector(i, s))
			jobs = append(jobs, j...)
		}
		return nil
	}
	if err := traverse(InitialSelector()); err != nil {
		Fatal(err)
	}

	return jobs, nil
}
