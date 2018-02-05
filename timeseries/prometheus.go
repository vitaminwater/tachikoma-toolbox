package timeseries

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var metricsRegistry = prometheus.NewRegistry()
var healthRegistery = prometheus.NewRegistry()

func init() {
	healthRegistery.MustRegister(prometheus.NewProcessCollector(os.Getpid(), ""))
	healthRegistery.MustRegister(prometheus.NewGoCollector())
}

func startMetricsHandle() {
	log.Info("Start metrics server")
	s := http.NewServeMux()
	s.Handle("/metrics", promhttp.HandlerFor(metricsRegistry, promhttp.HandlerOpts{}))
	http.ListenAndServe(":8080", s)
}

func startHealthHandle() {
	log.Info("Start health server")
	s := http.NewServeMux()
	s.Handle("/metrics", promhttp.HandlerFor(healthRegistery, promhttp.HandlerOpts{}))
	http.ListenAndServe(":8081", s)
}

func mustRegister(cs ...prometheus.Collector) {
	metricsRegistry.MustRegister(cs...)
}
