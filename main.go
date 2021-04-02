package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/buneyev/network-access-checker/checker"
	"github.com/buneyev/network-access-checker/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config.filepath", "config.json", "Path to the config file")
	flag.Parse()

	conf := models.NewConfig(configPath)
	resourceChecker := checker.NewChecker(conf, time.NewTicker(10*time.Second))
	go resourceChecker.Run()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("ListenAndServe on " + conf.MetricsPort)
	http.ListenAndServe(":"+conf.MetricsPort, nil)
}
