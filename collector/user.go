package collector

import (
	"encoding/json"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	headscale_user = typedDesc{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "headscale_user"),
			"Number of users", nil, nil, // label static
		),
		valueType: prometheus.GaugeValue,
	}
)

type UserResponce struct {
	Users []User `json:"users"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	// todo convert to date
}

func (e Exporter) getUser() (UserResponce, error) {
	var responseObject UserResponce
	responseData, err := e.queryPath("/user")
	if err != nil {
		return responseObject, err
	}
	json.Unmarshal(responseData, &responseObject)
	return responseObject, nil
}

func (e Exporter) gatherUsers(ch chan<- prometheus.Metric) {
	start := time.Now()
	defer func() {
		e.logger.Debug("Gathering users completed", "seconds", time.Since(start).Seconds())
	}()
	responseObject, err := e.getUser()
	if err != nil {
		e.logger.Error("Error gathering users", "error", err)
		return
	}
	total_users := len(responseObject.Users)
	ch <- headscale_user.mustNewConstMetric(float64(total_users))
}
