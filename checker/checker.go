package checker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/buneyev/network-access-checker/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Checker struct {
	config      *models.Config
	ctx         context.Context
	cancel      context.CancelFunc
	gaugeOpts   *prometheus.GaugeVec
	reqDuration *prometheus.HistogramVec
	closeFlag   chan bool
}

func NewChecker(config *models.Config) *Checker {
	return &Checker{
		config: config,
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
	chckr.ctx, chckr.cancel = context.WithCancel(context.Background())
	chckr.closeFlag = make(chan bool)
	var wg sync.WaitGroup
	for _, resource := range chckr.config.ResourceList {
		wg.Add(1)
		go chckr.worker(resource, &wg)
	}
	wg.Wait()
	chckr.closeFlag <- true
}

func (chckr *Checker) Reloader(reload chan bool) {
	for {
		<-reload
		fmt.Println("Start reload checker")
		chckr.Stop()
		close(chckr.closeFlag)
		go chckr.Run()
		fmt.Println("End reload checker")
	}
}

func (chckr *Checker) Stop() {
	chckr.cancel()
	<-chckr.closeFlag
}

func (chckr *Checker) worker(resource models.Resource, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Duration(chckr.config.CheckPeriodSeconds) * time.Second)
	for {
		select {
		case <-chckr.ctx.Done():
			wg.Done()
			return
		case <-ticker.C:
			if err := chckr.checkConnection(resource.Name, resource.Host, resource.Ports, time.Duration(chckr.config.CheckConnectionTimeout)*time.Second); err != nil {
				log.Printf("Connecting error: %v, host: %v, port: %v", err.Error(), resource.Host, resource.Ports)
			}
		}
	}
}
