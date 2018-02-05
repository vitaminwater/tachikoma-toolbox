package timeseries

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

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
	f, ok := d.(float64)
	if ok != true {
		return errors.New("GaugeMetric requires float64")
	}
	m.g.WithLabelValues(labels...).Set(f)
	return nil
}

func NewGaugeMetric(name, help string, labels []string) Metric {
	m := GaugeMetric{
		g: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: fmt.Sprintf("%s_gauge", name),
				Help: help,
			},
			labels,
		),
	}
	mustRegister(m.g)
	return m
}

/**
 * Summary
 */

type SummaryMetric struct {
	g *prometheus.SummaryVec
}

func (m SummaryMetric) Index(labels []string, d interface{}) error {
	f, ok := d.(float64)
	if ok != true {
		return errors.New("SummaryMetric requires float64")
	}

	m.g.WithLabelValues(labels...).Observe(f)
	return nil
}

func NewSummaryMetric(name, help string, labels []string) Metric {
	m := SummaryMetric{
		g: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: fmt.Sprintf("%s_summary", name),
				Help: help,
			},
			labels,
		),
	}
	mustRegister(m.g)
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
	//registry.MustRegister(m.g)
	return wc
}
