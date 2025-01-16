package collector

import (
	"encoding/json"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	headscale_policy_acl = typedDesc{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "policy_acl"),
			"Number of policies", nil, nil, // label static
		),
		valueType: prometheus.GaugeValue,
	}
	headscale_policy_group = typedDesc{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "policy_group"),
			"Number of groups", nil, nil, // label static
		),
		valueType: prometheus.GaugeValue,
	}
	headscale_policy_host = typedDesc{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "policy_host"),
			"Number of hosts", nil, nil, // label static
		),
		valueType: prometheus.GaugeValue,
	}
	headscale_policy_tag_owner = typedDesc{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "policy_tag_owner"),
			"Number of tag owners", nil, nil, // label static
		),
		valueType: prometheus.GaugeValue,
	}
	headscale_policy_group_members = typedDesc{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "policy_group_members"),
			"Number of members in each group",
			[]string{"group"}, // dynamic label
			nil,               // static label
		),
		valueType: prometheus.GaugeValue,
	}
)

type PolicyResponce struct {
	Policies  string `json:"policy"`
	UpdatedAt string `json:"updatedAt"`
}

type Policy struct {
	Acls      []Acl               `json:"acls"`
	Groups    map[string][]string `json:"groups"`
	Hosts     map[string]string   `json:"hosts"`
	TagOwners map[string][]string `json:"tagOwners"`
}

type Acl struct {
	Action      string   `json:"action"`
	Destination []string `json:"dst"`
	Source      []string `json:"src"`
	Protocol    string   `json:"proto"`
}

func (e Exporter) gatherPolicy(ch chan<- prometheus.Metric) (Policy, error) {
	start := time.Now()
	defer func() {
		e.logger.Debug("Gathering policy completed", "seconds", time.Since(start).Seconds())
	}()
	var responseObject PolicyResponce
	var policy Policy
	responseData, err := e.queryPath("/policy")
	if err != nil {
		e.logger.Error("Error gathering policy", "error", err)
		return policy, err
	}
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		e.logger.Error("Error parsing policy response", "error", err)
		return policy, err
	}

	err = json.Unmarshal([]byte(responseObject.Policies), &policy)
	if err != nil {
		e.logger.Error("Error parsing policy", "error", err)
		return policy, err
	}

	total_policies := len(policy.Acls)
	ch <- headscale_policy_acl.mustNewConstMetric(float64(total_policies))

	total_groups := len(policy.Groups)
	ch <- headscale_policy_group.mustNewConstMetric(float64(total_groups))

	total_hosts := len(policy.Hosts)
	ch <- headscale_policy_host.mustNewConstMetric(float64(total_hosts))

	total_tag_owners := len(policy.TagOwners)
	ch <- headscale_policy_tag_owner.mustNewConstMetric(float64(total_tag_owners))

	for group, members := range policy.Groups {
		ch <- headscale_policy_group_members.mustNewConstMetric(float64(len(members)), group)
	}

	return policy, nil
}
