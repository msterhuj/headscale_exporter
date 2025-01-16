package collector

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "headscale"
)

type CollectorConfig struct {
	Address string
	Token   string
}

// define all metrics here

type typedDesc struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
}

func (d *typedDesc) mustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(d.desc, d.valueType, value, labels...)
}

type Exporter struct {
	address string
	headers http.Header
	logger  *slog.Logger
}

func NewExporter(config CollectorConfig, logger *slog.Logger) Exporter {
	return Exporter{
		address: config.Address,
		headers: http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer " + config.Token},
			"User-Agent":    {"headscale_exporter"},
		},
		logger: logger,
	}
}

func (e Exporter) Describe(ch chan<- *prometheus.Desc) {}

func (e Exporter) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	e.logger.Debug("Scrape starting")
	defer func() {
		e.logger.Debug("Scrape completed", "seconds", time.Since(start).Seconds())
	}()
	e.gatherApiKeys(ch)
	e.gatherUsers(ch)
}

func (e Exporter) queryPath(path string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", e.address+path, nil)
	req.Header = e.headers
	if err != nil {
		e.logger.Error("Failed to build new request", "error", err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		e.logger.Error("Failed request on", "path", path, "error", err)
		return nil, err
	}
	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		e.logger.Error("Failed to read response", "error", err)
		return nil, err
	}
	return responseData, nil
}

func (e Exporter) ValidateToken() bool {
	e.logger.Debug("Validating token...")
	_, err := e.queryPath("/apikey")
	if err != nil {
		e.logger.Error("Invalid token", "error", err)
		return false
	}
	e.logger.Debug("Token is valid")
	return true
}
