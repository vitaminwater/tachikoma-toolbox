package tachikoma

import (
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var MetricsRegistry = prometheus.NewRegistry()
var HealthRegistery = prometheus.NewRegistry()

func init() {
	HealthRegistery.MustRegister(prometheus.NewProcessCollector(os.Getpid(), ""))
	HealthRegistery.MustRegister(prometheus.NewGoCollector())
}

func StartMetricsHandle() {
	log.Info("Start metrics server")
	s := http.NewServeMux()
	s.Handle("/metrics", promhttp.HandlerFor(MetricsRegistry, promhttp.HandlerOpts{}))
	http.ListenAndServe(":8080", s)
}

func StartHealthHandle() {
	log.Info("Start health server")
	s := http.NewServeMux()
	s.Handle("/metrics", promhttp.HandlerFor(HealthRegistery, promhttp.HandlerOpts{}))
	http.ListenAndServe(":8081", s)
}

func Start() {
	go StartMetricsHandle()
	go StartHealthHandle()
	select {}
}

func MustRegister(cs ...prometheus.Collector) {
	MetricsRegistry.MustRegister(cs...)
}

func Labels(source, base, counter string) []string {
	return []string{source, strings.ToUpper(base), strings.ToUpper(counter)}
}
