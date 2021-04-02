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
	Timeout      time.Duration
}

func NewChecker(config *models.Config) *Checker {
	return &Checker{
		ResourceList: config.ResourceList,
		Ticker:       time.NewTicker(time.Duration(config.CheckPeriodSeconds) * time.Second),
		Timeout:      time.Duration(config.CheckConnectionTimeout * int(time.Second)),
	}
}

func (checker *Checker) Run() {
	gaugeOpts := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "monitoring",
			Subsystem: "network",
			Name:      "access_checker",
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
				if err := checkConnection(resource.Host, resource.Ports, checker.Timeout); err != nil {
					gaugeOpts.WithLabelValues(resource.Name, resource.Host).Set(0)
				} else {
					gaugeOpts.WithLabelValues(resource.Name, resource.Host).Set(1)
				}
			}
		}
	}
}
