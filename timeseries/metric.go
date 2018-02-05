package timeseries

import (
	"fmt"
	"reflect"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
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
		logrus.Info("gaugeVec ", name)
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
		logrus.Info("summaryVec ", name)
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
