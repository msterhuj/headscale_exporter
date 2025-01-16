package collector

import (
	"encoding/json"

	"github.com/prometheus/client_golang/prometheus"
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

func (e Exporter) gatherApiKeys(ch chan<- prometheus.Metric) {
	var responseObject ApiKeysResponce
	responseData, err := e.queryPath("/apikey")
	if err != nil {
		e.logger.Error("Error gathering api keys", "error", err)
		return
	}
	json.Unmarshal(responseData, &responseObject)
	total_api_keys := len(responseObject.ApiKeys)
	ch <- headscale_api_keys.mustNewConstMetric(float64(total_api_keys))
}
