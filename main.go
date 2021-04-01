package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/buneyev/network-access-checker/checker"
	"github.com/buneyev/network-access-checker/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	conf := models.NewConfig("config.json")
	resourceChecker := checker.NewChecker(conf, time.NewTicker(10*time.Second))
	go resourceChecker.Run()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("ListenAndServe on " + conf.MetricsPort)
	http.ListenAndServe(":"+conf.MetricsPort, nil)
}
