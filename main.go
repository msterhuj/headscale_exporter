package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

const (
	namespace = "headscale"
)

// Parameter
type CollectorConfig struct {
	Address string
}

var (
	conf   = CollectorConfig{}
	logger *slog.Logger
)

// define metrics structure and list
type typedDesc struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
}

func (d *typedDesc) mustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(d.desc, d.valueType, value, labels...)
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

// Exporter with collector
type Exporter struct {
	address string

	logger *slog.Logger
}

func NewExporter(config CollectorConfig, logger *slog.Logger) Exporter {
	return Exporter{
		address: conf.Address,
		logger:  logger,
	}
}

func (e Exporter) Describe(ch chan<- *prometheus.Desc) {}

func (e Exporter) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	logger.Debug("Scrape starting")
	defer func() {
		logger.Debug("Scrape completed", "seconds", time.Since(start).Seconds())
	}()

	ch <- up.mustNewConstMetric(1)
}

func main() {
	// parse args
	kingpin.Flag(
		"headscale.endpoint",
		"Endpoint of the headscale API",
	).Default("127.0.0.1:8080").StringVar(&conf.Address)

	metricsPath := kingpin.Flag(
		"web.telemetry-path",
		"Path under which to exporse metrics.",
	).Default("/metrics").String()

	toolkitFlags := kingpinflag.AddFlags(kingpin.CommandLine, ":9191")

	promslogConfig := &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version.Print("headscale_exporter"))
	kingpin.Parse()

	logger = promslog.New(promslogConfig)
	logger.Info("Starting headscale_exporter", "version", version.Info())

	// register exporter metrics
	exporter := NewExporter(conf, logger)
	prometheus.MustRegister(exporter)

	// handle data
	http.Handle(*metricsPath, promhttp.Handler())
	if *metricsPath != "/" && *metricsPath != "" {
		landingConfig := web.LandingConfig{
			Name:        "Headscale Exporter",
			Description: "Prometheus Exporter for Headscale",
			Version:     version.Info(),
			Links: []web.LandingLinks{
				{
					Address: *metricsPath,
					Text:    "Metrics",
				}, {
					Address: "https://headscale.net/",
					Text:    "HeadScale",
				},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			logger.Error("error creating landing page", "err", err)
			os.Exit(1)
		}
		http.Handle("/", landingPage)
	}
	server := &http.Server{}
	if err := web.ListenAndServe(server, toolkitFlags, logger); err != nil {
		logger.Error("HTTP listener stopped", "error", err)
		os.Exit(1)
	}
}