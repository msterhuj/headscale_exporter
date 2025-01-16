package collector

import (
	"encoding/json"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	headscale_node = typedDesc{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "headscale_node"),
			"Number of nodes", nil, nil, // label static
		),
		valueType: prometheus.GaugeValue,
	}
	headscale_node_online = typedDesc{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "headscale_node_online"),
			"Number of online nodes", nil, nil, // label static
		),
		valueType: prometheus.GaugeValue,
	}
)

type NodeResponce struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Id          string   `json:"id"`
	MachineKey  string   `json:"machineKey"`
	NodeKey     string   `json:"nodeKey"`
	DiscoKey    string   `json:"discoKey"`
	IpAddresses []string `json:"ipAddresses"`
	Name        string   `json:"name"`
	User        User     `json:"user"`
	LastSeen    string   `json:"lastSeen"`
	Expiry      string   `json:"expiry"`
	// todo add preAuthKey
	CreatedAt      string   `json:"createdAt"`
	RegisterMethod string   `json:"registerMethod"`
	ForcedTags     []string `json:"forcedTags"`
	InvalidTags    []string `json:"invalidTags"`
	ValidTags      []string `json:"validTags"`
	GivenName      string   `json:"givenName"`
	Online         bool     `json:"online"`
}

func (e Exporter) gatherNodes(ch chan<- prometheus.Metric) (NodeResponce, error) {
	start := time.Now()
	defer func() {
		e.logger.Debug("Gathering nodes completed", "seconds", time.Since(start).Seconds())
	}()
	var responseObject NodeResponce
	responseData, err := e.queryPath("/node")
	if err != nil {
		e.logger.Error("Error gathering nodes", "error", err)
		return responseObject, err
	}
	e.logger.Debug("Parsing nodes response", "response", string(responseData))
	json.Unmarshal(responseData, &responseObject)

	// Count nodes
	total_nodes := len(responseObject.Nodes)
	ch <- headscale_node.mustNewConstMetric(float64(total_nodes))

	// Count online nodes
	total_online_nodes := 0
	for _, node := range responseObject.Nodes {
		if node.Online {
			total_online_nodes++
		}
	}
	ch <- headscale_node_online.mustNewConstMetric(float64(total_online_nodes))

	return responseObject, nil
}
