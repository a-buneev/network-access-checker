package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/buneyev/network-access-checker/checker"
	"github.com/buneyev/network-access-checker/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config.filepath", "config.json", "Path to the config file")
	flag.Parse()

	conf := models.NewConfig(configPath)
	resourceChecker := checker.NewChecker(conf)
	go resourceChecker.Run()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("ListenAndServe on " + conf.MetricsPort)
	http.ListenAndServe(":"+conf.MetricsPort, nil)
}
