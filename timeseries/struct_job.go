package timeseries

import (
	"errors"
	"reflect"
	"strings"

	tachikoma "github.com/vitaminwater/tachikoma-toolbox"
)

/**

Each of the structs fields of the Data.Data field with the required tags are treated as one or multiple time series.
A time series value requires these data to be predefined:

	- name
	- helper
	- labels

The tags that makeup a timeseries are:

  - tsname
	- tshelper
	- tstype

tstype indicates the type of indexing to perform.
defaults to:

  - for string type: word count
	- for float64 type: gauge & summary

*/

type Timeserie struct {
	Name string
	Help string
	Type string
}

var REQUIRED_TAGS = []string{
	"tsname", "tshelper", "tstype",
}

var DefaultMetricsByName = map[string]func(string, string, []string) Metric{
	"word_count": NewWordCount,
	"gauge":      NewGaugeMetric,
	"summary":    NewSummaryMetric,
}

type TimeseriesJobGenerator func(timeserie Timeserie, labels Labels, f reflect.StructField, selector tachikoma.Selector) []tachikoma.Job

func GenerateTimeseries(defaultLabels Labels, g TimeseriesJobGenerator) tachikoma.StructJobGenerator {
	labels := defaultLabels.Clone()

	return func(f reflect.StructField, selector tachikoma.Selector) []tachikoma.Job {
		tag := f.Tag
		jobs := make([]tachikoma.Job, 0)

		var tsname, tshelp, tstypes, tslabel string
		var ok bool
		if tsname, ok = tag.Lookup("tsname"); ok == false {
			return jobs
		}
		if tshelp, ok = tag.Lookup("tshelp"); ok == false {
			return jobs
		}

		if tslabel, ok = tag.Lookup("tslabel"); ok == true {
			tslabel = strings.Replace(tslabel, " ", "", -1)
			pairs := strings.Split(tslabel, ",")
			for _, p := range pairs {
				kv := strings.Split(p, ":")
				if len(kv) != 2 {
					tachikoma.Fatal(errors.New("Additional labels in tag should be in the form key:value"))
				}
				labels[kv[0]] = StaticLabel(kv[1])
			}
		}

		if tstypes, ok = tag.Lookup("tstypes"); ok == false {
			k := f.Type.Kind()
			switch k {
			case reflect.String:
				tstypes = "word_count"
			case reflect.Float64:
				tstypes = "gauge,summary"
			default:
				tachikoma.Fatal(errors.New("No types specified and unable to set default based on Kind"))
			}
		}

		ts := strings.Split(tstypes, ",")
		for _, t := range ts {
			timeserie := Timeserie{
				Name: tsname,
				Help: tshelp,
				Type: t,
			}
			jobs = append(jobs, g(timeserie, labels, f, selector)...)
		}

		return jobs
	}
}
