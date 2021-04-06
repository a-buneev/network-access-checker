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
	gaugeOpts   *prometheus.GaugeVec
	reqDuration *prometheus.HistogramVec
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
		gaugeOpts: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "monitoring",
				Subsystem: "network",
				Name:      "access_status",
				Help:      "Check network access to resources",
			}, []string{"resource_name", "resource_addr", "resource_port"}),
		reqDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "monitoring",
			Subsystem: "network",
			Name:      "req_duration",
			Help:      "Duration of requests",
		}, []string{"resource_name", "resource_addr", "resource_port"}),
	}
}

func (chckr *Checker) Run() {
	var wg sync.WaitGroup
	for _, resource := range chckr.ResourceList {
		wg.Add(1)
		go chckr.worker(resource, &wg)
	}
	wg.Wait()
}

func (chckr *Checker) Stop() {
	chckr.cancel()
}

func (chckr *Checker) worker(resource models.Resource, wg *sync.WaitGroup) {
	ticker := time.NewTicker(chckr.CheckPeriod)
	for {
		select {
		case <-chckr.ctx.Done():
			wg.Done()
			return
		case <-ticker.C:
			chckr.checkConnection(resource.Name, resource.Host, resource.Ports, chckr.Timeout)
		}
	}
}
