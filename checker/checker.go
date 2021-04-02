package checker

import (
	"time"

	"github.com/buneyev/network-access-checker/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Checker struct {
	ResourceList []models.Resource
	Ticker       *time.Ticker
}

func NewChecker(config *models.Config, ticker *time.Ticker) *Checker {
	return &Checker{
		ResourceList: config.ResourceList,
		Ticker:       ticker,
	}
}

func (checker *Checker) Run() {
	gaugeOpts := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "monitoring",
			Subsystem: "network",
			Name:      "network_access_checker",
			Help:      "Check network access to resources",
		},
		[]string{"resourceName", "resourceAddr"},
	)
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case <-checker.Ticker.C:
			for _, resource := range checker.ResourceList {
				if err := checkConnection(resource.Host, resource.Ports); err != nil {
					gaugeOpts.WithLabelValues(resource.Name, resource.Host).Set(0)
				} else {
					gaugeOpts.WithLabelValues(resource.Name, resource.Host).Set(1)
				}
			}
		}
	}
}
