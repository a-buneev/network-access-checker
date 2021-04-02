package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/buneyev/network-access-checker/checker"
	"github.com/buneyev/network-access-checker/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestMain(t *testing.T) {
	resourceChecker := checker.NewChecker(models.NewConfig("config.json"))
	go resourceChecker.Run()
	time.Sleep(6 * time.Second)
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), "monitoring_network_network_access_checker{resourceAddr=\"www.213y4y.com\",resourceName=\"for testing\"} 0") {
		t.Errorf("handler returned unexpected body: www.213y4y.com not found")
	}
	if !strings.Contains(rr.Body.String(), "monitoring_network_network_access_checker{resourceAddr=\"www.google.com\",resourceName=\"Google\"} 1") {
		t.Errorf("handler returned unexpected body: www.google.com not found")
	}
	if !strings.Contains(rr.Body.String(), "monitoring_network_network_access_checker{resourceAddr=\"www.yandex.com\",resourceName=\"Yandex\"} 1") {
		t.Errorf("handler returned unexpected body: www.yandex.com not found")
	}
}
