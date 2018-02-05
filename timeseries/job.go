package timeseries

import (
	"fmt"

	tachikoma "github.com/vitaminwater/tachikoma-toolbox"
)

type Job struct {
	Name     string
	Selector tachikoma.Selector
	Labels   Labels
	Metric   Metric
}

func (j Job) GetName() string {
	return j.Name
}

func (j Job) Run(i interface{}) error {
	return j.Metric.Index(j.Labels.Values(i), j.Selector(i))
}

func NewJob(timeserie Timeserie, labels Labels, selector tachikoma.Selector) Job {
	mn, ok := DefaultMetricsByName[timeserie.Type]
	if ok == false {
		tachikoma.Fatal(fmt.Errorf("Unknown metric name %s", timeserie.Type))
	}

	j := Job{
		Name:     timeserie.Name,
		Labels:   labels,
		Selector: selector,
		Metric:   mn(timeserie.Name, timeserie.Help, labels.Keys()),
	}
	return j
}
