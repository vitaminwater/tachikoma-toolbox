package timeseries

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/prometheus/client_golang/prometheus"
	tachikoma "github.com/vitaminwater/tachikoma-toolbox"
)

type metricRepository map[string]interface{}

func (r metricRepository) gaugeVec(name, help string, labels []string) *prometheus.GaugeVec {
	name = fmt.Sprintf("g_%s", name)
	if m, ok := r[name]; ok == false {
		r[name] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: name,
				Help: help,
			},
			labels,
		)
		mustRegister(r[name].(prometheus.Collector))
		return r[name].(*prometheus.GaugeVec)
	} else {
		return m.(*prometheus.GaugeVec)
	}
}

func (r metricRepository) summaryVec(name, help string, labels []string) *prometheus.SummaryVec {
	name = fmt.Sprintf("s_%s", name)
	if m, ok := r[name]; ok == false {
		r[name] = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: name,
				Help: help,
			},
			labels,
		)
		mustRegister(r[name].(prometheus.Collector))
		return r[name].(*prometheus.SummaryVec)
	} else {
		return m.(*prometheus.SummaryVec)
	}
}

var repo metricRepository = metricRepository{}

type Metric interface {
	Index([]string, interface{}) error
}

/**
 * Gauge
 */

type GaugeMetric struct {
	g *prometheus.GaugeVec
}

func (m GaugeMetric) Index(labels []string, d interface{}) error {
	if reflect.TypeOf(d).Kind() != reflect.Float64 {
		tachikoma.Fatal(errors.New("GaugeMetric requires a float64 value"))
	}

	m.g.WithLabelValues(labels...).Set(reflect.ValueOf(d).Float())
	return nil
}

func NewGaugeMetric(name, help string, labels []string) Metric {
	m := GaugeMetric{
		g: repo.gaugeVec(name, help, labels),
	}
	return m
}

/**
 * Summary
 */

type SummaryMetric struct {
	g *prometheus.SummaryVec
}

func (m SummaryMetric) Index(labels []string, d interface{}) error {
	if reflect.TypeOf(d).Kind() != reflect.Float64 {
		tachikoma.Fatal(errors.New("SummaryMetric requires a float64 value"))
	}

	m.g.WithLabelValues(labels...).Observe(reflect.ValueOf(d).Float())
	return nil
}

func NewSummaryMetric(name, help string, labels []string) Metric {
	m := SummaryMetric{
		g: repo.summaryVec(name, help, labels),
	}
	return m
}

/**
 * WordCount
 */

type WordCount struct {
}

func (wc WordCount) Index(labels []string, d interface{}) error {
	return nil
}

func NewWordCount(name, help string, labels []string) Metric {
	wc := WordCount{}
	return wc
}
