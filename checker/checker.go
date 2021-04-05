package checker

import (
	"context"
	"sync"
	"time"

	"github.com/buneyev/network-access-checker/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Checker struct {
	ResourceList []models.Resource
	//Ticker       *time.Ticker
	CheckPeriod time.Duration
	Timeout     time.Duration
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewChecker(config *models.Config) *Checker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Checker{
		ResourceList: config.ResourceList,
		//Ticker:       time.NewTicker(time.Duration(config.CheckPeriodSeconds) * time.Second),
		CheckPeriod: time.Duration(config.CheckPeriodSeconds) * time.Second,
		Timeout:     time.Duration(config.CheckConnectionTimeout * int(time.Second)),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (checker *Checker) Run() {
	gaugeOpts := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "monitoring",
			Subsystem: "network",
			Name:      "access_status",
			Help:      "Check network access to resources",
		},
		[]string{"resource_name", "resource_addr"},
	)
	var wg sync.WaitGroup
	for _, resource := range checker.ResourceList {
		wg.Add(1)
		go checker.worker(resource, gaugeOpts, &wg)
	}
	wg.Wait()
}

func (checker *Checker) Stop() {
	checker.cancel()
}

func (checker *Checker) worker(resource models.Resource, gaugeOpts *prometheus.GaugeVec, wg *sync.WaitGroup) {
	ticker := time.NewTicker(checker.CheckPeriod)
	for {
		select {
		case <-checker.ctx.Done():
			wg.Done()
			return
		case <-ticker.C:
			if err := checkConnection(resource.Host, resource.Ports, checker.Timeout); err != nil {
				gaugeOpts.WithLabelValues(resource.Name, resource.Host).Set(0)
			} else {
				gaugeOpts.WithLabelValues(resource.Name, resource.Host).Set(1)
			}
		}
	}
}
