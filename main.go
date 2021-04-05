package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	resourceChecker := checker.NewChecker(conf)
	go resourceChecker.Run()

	server := newMetricsServer(conf)
	fmt.Println("ListenAndServe on " + conf.MetricsPort)
	go server.ListenAndServe()

	gracefulShutdown(resourceChecker, server)
	fmt.Println("Success end")
}

func gracefulShutdown(checker *checker.Checker, server *http.Server) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	signal.Notify(sigint, syscall.SIGTERM)
	<-sigint
	fmt.Println("Start graceful shutdown")
	checker.Stop()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}
	fmt.Println("Checker stoped")
}

func newMetricsServer(conf *models.Config) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return &http.Server{
		Addr:    ":" + conf.MetricsPort,
		Handler: mux,
	}
}
