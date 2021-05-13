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
	"github.com/buneyev/network-access-checker/server"
)

func main() {
	var (
		configPath                    string
		autoReloadConfig              bool
		autoReloadConfigPeriodSeconds int
		serverPort                    string
	)

	flag.StringVar(&configPath, "config.filepath", "config.json", "Path to the config file")
	flag.BoolVar(&autoReloadConfig, "config.autoreload", false, "Autoreload config file")
	flag.IntVar(&autoReloadConfigPeriodSeconds, "config.autoreloadperiod", 4, "Period of autoreload config file")
	flag.StringVar(&serverPort, "server.port", "2112", "Port for ListenAndServe")
	flag.Parse()

	conf := models.NewConfig(configPath)
	resourceChecker := checker.NewChecker(conf)
	go resourceChecker.Run()

	if autoReloadConfig {
		reload := make(chan bool)
		go conf.Reload(autoReloadConfigPeriodSeconds, reload)
		go resourceChecker.Reloader(reload)
	}

	srv := server.NewMetricsServer(serverPort)
	fmt.Println("ListenAndServe on :" + serverPort)
	go func() {
		log.Println(srv.ListenAndServe())
	}()

	gracefulShutdown(resourceChecker, srv)
	fmt.Println("Success end")
}

func gracefulShutdown(checker *checker.Checker, srv *http.Server) {
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

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}
	fmt.Println("Checker stoped")
}
