package timeseries

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type Metric interface {
	Index([]string, interface{})
	Register(*prometheus.Registry)
}

/**
 * Gauge
 */

type GaugeMetric struct {
	g *prometheus.GaugeVec
}

func (m GaugeMetric) Index(labels []string, d interface{}) error {
	f, ok := d.(float64)
	if ok != true {
		return errors.New("GaugeMetric requires float64")
	}
	m.g.WithLabelValues(labels...).Set(f)
}

func NewGaugeMetric(labels []string) GaugeMetric {
	m := GaugeMetric{
		g: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: fmt.Sprintf("%s_gauge", name),
				Help: help,
			},
			labels,
		),
	}
	registry.MustRegister(m.g)
	return m
}

/**
 * Summary
 */

type SummaryMetric struct {
}

func (m SummaryMetric) Index(labels []string, d interface{}) error {
	f, ok := d.(float64)
	if ok != true {
		return errors.New("SummaryMetric requires float64")
	}

	m.s.WithLabelValues(labels...).Observe(f)
}

func NewSummaryMetric(labels []string, registry *prometheus.Registry) SummaryMetric {
	m := SummaryMetric{
		g: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: fmt.Sprintf("%s_summary", name),
				Help: help,
			},
			labels,
		),
	}
	registry.MustRegister(m.g)
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

func NewWordCount(labels []string) WordCount {
	wc := WordCount{}
	//registry.MustRegister(m.g)
	return wc
}
