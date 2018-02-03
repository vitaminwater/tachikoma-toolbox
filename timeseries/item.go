package timeseries

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

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
var REQUIRED_TAGS = []string{
	"tsname", "tshelper", "tstype",
}

/**
 * Allows to back-and-forth between the struct field names and their values
 */

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

type LabelFn func(Field, interface{}, []string) []string

type Field struct {
	Name     string
	Tag      reflect.StructTag
	Selector Selector
	Metrics  []Metric
	LabelFn  LabelFn
}

func (p Field) process(labels []string, d interface{}) {
	v := p.Selector(d)
	for _, m := range p.Metrics {
		ls := labels
		if p.LabelFn != nil {
			ls = p.LabelFn(p, d, ls)
		}
		m.Index(ls, v)
	}
}

var MetricsByName = map[string]func() Metric{
	"word_count": NewWordCount,
	"gauge":      NewGaugeMetric,
	"summary":    NewSummaryMetric,
}

func NewField(tag reflect.Tag, selector Selector) Field {
	if types, ok = tag.Lookup("tstype"); ok == false {
		switch k {
		case reflect.String:
			types = "word_count"
		case reflect.Float64:
			types = "gauge,summary"
		default:
			continue
		}
	}

	metrics := make([]Metric, 0)
	ts := strings.Split(types, ",")
	for _, t := range ts {
		mn, ok := MetricsByName[t]
		if ok == false {
			continue
		}
		metrics = append(metrics, MetricsByName[mn])
	}

	f := Field{
		Name:     tag.Lookup("tsname"),
		Tag:      tag,
		Selector: selector,
		Metrics:  metrics,
	}
}

/**
 * Source
 */

type Source struct {
	Name   string
	Fields []Field
}

type SourceItem struct {
	Source Source
	Origin string
	Labels []string
	Data   Data
}

func (si SourceItem) Index() {
	for _, f := range si.Source.Fields {
		f.process(si.Label, si.Data.Data)
	}
}

func (s Source) NewSourceItem(origin string, d interface{}) (SourceItem, error) {
	si := SourceItem{
		Source: s,
		Origin: origin,
		Data:   NewData(d),
	}
	return si
}

func (s Source) Field(name string) (Field, error) {
	for _, f := range s.Fields {
		if f.Name == name {
			return f, nil
		}
	}
	return Field{}, fmt.Errorf("Unknown field %s", name)
}

func NewSource(name string, d interface{}) Source {
	if d == nil {
		return errors.New("d is nil.")
	}

	fields := make([]Field)
	traverse := func(s Selector) error {
		d := s(d)
		d := reflect.ValueOf(d)
		t := reflect.TypeOf(d)

		if t.Kind() == reflect.Ptr {
			d = d.Elem()
			t = t.Elem()
		}

		if t.Kind() != reflect.Struct {
			return nil, errors.New("data should be struct or pointer to struct")
		}

		for i := 0; i < t.NumField(); i++ {
			var name, helper, types string
			var ok bool
			k := t.FieldByIndex(i).Type.Kind()
			if k == reflect.Struct {
				traverse(IndexSelector(i, s))
			}
			tag := t.Field(i).Tag
			for _, r := range REQUIRED_TAGS {
				if _, ok = tag.Lookup(r); ok == false {
					continue
				}
			}
			fields = append(fields, NewField(tag, IndexSelector(i, s)))
		}
	}
	if err := traverse(InitialSelector()); err != nil {
		Fatal(err)
	}

	i := Source{
		Name:   name,
		Fields: fields,
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
