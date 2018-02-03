package timeseries

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

/**

Each of the structs fields of the Data.Data field with the required tags are treated as one or multiple time series.
A time series value requires these data to be predefined:

	- name
	- helper
	- labels

In most cases labels are the same for a group of metrics, so they are defined at once for an item.

The tags that makeup a timeseries are:

  - tsname
	- tshelper
	- tstype

tstype indicates the type of indexing to perform.
defaults to:

  - for string type: word count
	- for float64 type: gauge & summary

*/

/**
 * Allows to back-and-forth between the struct field names and their values
 */

type Selector func(interface{}) interface{}

func InitialSelector() Selector {
	return func(o interface{}) interface{} {
		return o
	}
}

func IndexSelector(i uint, s Selector) Selector {
	return func(o interface{}) interface{} {
		o = s(o)
		d := reflect.ValueOf(o)
		t := reflect.TypeOf(o)

		if t.Kind() == reflect.Ptr {
			d = d.Elem()
			t = t.Elem()
		}

		return d.FieldByIndex(i).Interface()
	}
}

type Field struct {
	Name     string
	Selector Selector
	Metrics  []Metric
}

func (p Field) process(d interface{}) {
	v := p.Selector(d)
	for _, m := range p.Metrics {
		m(v)
	}
}

/**
 * Source
 */

type Source struct {
	Name   string
	Prefix string
	Fields []Field
}

type SourceItem struct {
	SourceName string
	Origin     string
	Labels     []string
	Data       interface{}
	Metrics    map[string]interface{}
}

func (i Source) Index(d interface{}) error {
	traverse := func(d interface{}) {

	}
	if d == nil {
		return errors.New("d is nil.")
	}

	d := reflect.ValueOf(d)
	t := reflect.TypeOf(d)

	if t.Kind() == reflect.Ptr {
		d = d.Elem()
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, nil, errors.New("first argument should be struct or pointer to struct")
	}

	for i := 0; i < t.NumField(); i++ {
		var name, helper, types string
		var ok bool
		if name, ok = t.Field(i).Tag.Lookup("tsname"); ok == false {
			continue
		}
		if helper, ok = t.Field(i).Tag.Lookup("tshelper"); ok == true {
			continue
		}
		if types, ok = t.Field(i).Tag.Lookup("tstype"); ok == true {
			switch _ := d.FieldByIndex(i).Interface().(type) {
			case string:
				types = "word_count"
			case float64:
				types = "gauge,summary"
			default:
				continue
			}
		}

	}
	return nil
}

func NewSource(name, prefix string, data interface{}) Source {
	i := Source{
		Name:        name,
		Prefix:      prefix,
		selectorMap: make(map[string]interface{}),
	}
	return i
}

/**
 * Data wraps the scraped data
 */

type Data struct {
	Data interface{}
}

func (di Data) Value() (driver.Value, error) {
	j, err := json.Marshal(di)
	return j, err
}

func (di *Data) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	v := map[string]interface{}{}
	err := json.Unmarshal(source, &v)
	if err != nil {
		return err
	}

	if di.Data == nil {
		di.Data = v
		return nil
	}

	if err := mapstructure.Decode(v, &di.Data); err != nil {
		return err
	}
	return nil
}

func NewData(data interface{}) Data {
	d := Data{
		Data: data,
	}
	return d
}
