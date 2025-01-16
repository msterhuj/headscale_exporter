package collector

import (
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "headscale"
)

type CollectorConfig struct {
	Address string
}

var (
	up = typedDesc{
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Headscale UP", // help
			nil,            //[]string{"host"}, // label dynamique
			nil,            // label static
		),
		prometheus.GaugeValue,
	}
)

type typedDesc struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
}

func (d *typedDesc) mustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(d.desc, d.valueType, value, labels...)
}

type Exporter struct {
	address string

	logger *slog.Logger
}

func NewExporter(config CollectorConfig, logger *slog.Logger) Exporter {
	return Exporter{
		address: config.Address,
		logger:  logger,
	}
}

func (e Exporter) Describe(ch chan<- *prometheus.Desc) {}

func (e Exporter) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	e.logger.Debug("Scrape starting")
	defer func() {
		e.logger.Debug("Scrape completed", "seconds", time.Since(start).Seconds())
	}()

	ch <- up.mustNewConstMetric(1)
}
