package collector

import (
	"encoding/json"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	headscale_api_keys = typedDesc{
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "headscale_api_keys"),
			"Number of API keys",
			nil, //[]string{"host"}, // label dynamique
			nil, // label static
		),
		prometheus.GaugeValue,
	}
)

type ApiKeysResponce struct {
	ApiKeys []ApiKey `json:"apiKeys"`
}

type ApiKey struct {
	Id         string `json:"id"`
	Prefix     string `json:"prefix"`
	Expiration string `json:"expiration"`
	CreatedAt  string `json:"createdAt"`
	LastSeen   string `json:"lastSeen"`
	// todo convert to date
}

func (e Exporter) gatherApiKeys(ch chan<- prometheus.Metric) (ApiKeysResponce, error) {
	start := time.Now()
	defer func() {
		e.logger.Debug("Gathering api keys completed", "seconds", time.Since(start).Seconds())
	}()
	var responseObject ApiKeysResponce
	e.logger.Debug("Gathering api keys")
	responseData, err := e.queryPath("/apikey")
	if err != nil {
		e.logger.Error("Error gathering api keys", "error", err)
		return responseObject, err
	}
	e.logger.Debug("Parsing api keys response", "response", string(responseData))
	json.Unmarshal(responseData, &responseObject)
	total_api_keys := len(responseObject.ApiKeys)
	ch <- headscale_api_keys.mustNewConstMetric(float64(total_api_keys))
	return responseObject, nil
}
