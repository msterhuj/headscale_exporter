package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/msterhuj/headscale_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

// Parameter

var (
	conf   = collector.CollectorConfig{}
	logger *slog.Logger
)

func main() {
	// parse args
	kingpin.Flag(
		"headscale.endpoint",
		"Endpoint of the headscale API",
	).Default("http://127.0.0.1:8080/api/v1").StringVar(&conf.Address)

	kingpin.Flag(
		"headscale.token",
		"Token for the headscale API",
	).Required().StringVar(&conf.Token)

	metricsPath := kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
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
	exporter := collector.NewExporter(conf, logger)
	if !exporter.ValidateToken() {
		os.Exit(1)
	}
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
