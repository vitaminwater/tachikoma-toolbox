package timeseries

import (
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

func (m GaugeMetric) Index(ls []string, d interface{}) {
	f, ok := d.(float64)
	if ok != true {
		return
	}
	m.g.WithLabelValues(labels...).Set(f)
}

func (m GaugeMetric) Register(registry *prometheus.Registry) {
	registry.MustRegister(m.g)
}

func NewGaugeMetric() GaugeMetric {
	m := GaugeMetric{
		g: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: fmt.Sprintf("%s_gauge", name),
				Help: help,
			},
			[]string{"source", "base", "counter"},
		),
	}
	return m
}

/**
 * Summary
 */

type SummaryMetric struct {
}

func (m SummaryMetric) Index(ls []string, d interface{}) {
	f, ok := d.(float64)
	if ok != true {
		return
	}

	m.s.WithLabelValues(labels...).Observe(f)
}

func (m SummaryMetric) Register(registry *prometheus.Registry) {
	registry.MustRegister(m.g)
}

func NewSummaryMetric() SummaryMetric {
	m := SummaryMetric{
		g: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: fmt.Sprintf("%s_summary", name),
				Help: help,
			},
			[]string{"source", "base", "counter"},
		),
	}
	return m
}

/**
 * WordCount
 */

type WordCount struct {
}

func (wc WordCount) Index(ls []string, d interface{}) {

}
