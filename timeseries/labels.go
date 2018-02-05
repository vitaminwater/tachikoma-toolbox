package timeseries

import "sort"

type LabelFn func(string, interface{}) string
type Labels map[string]LabelFn

func (ls Labels) Keys() []string {
	keys := make([]string, len(ls))

	i := 0
	for k := range ls {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	return keys
}

func (ls Labels) Values(d interface{}) []string {
	values := make([]string, len(ls))
	ks := ls.Keys()
	for i, k := range ks {
		values[i] = ls[k](k, d)
	}
	return values
}

func (ls Labels) AsMap(i interface{}) map[string]string {
	m := make(map[string]string)
	for k, fn := range ls {
		m[k] = fn(k, i)
	}
	return m
}

func (ls Labels) Clone() Labels {
	labels := make(Labels)
	for k, v := range ls {
		labels[k] = v
	}
	return labels
}

func StaticLabel(val string) LabelFn {
	return func(string, interface{}) string {
		return val
	}
}
