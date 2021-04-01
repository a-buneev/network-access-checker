package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/buneyev/network-access-checker/checker"
	"github.com/buneyev/network-access-checker/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestMain(t *testing.T) {
	resourceChecker := checker.NewChecker(models.NewConfig("config.json"), time.NewTicker(3*time.Second))
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
	fmt.Println(rr.Body.String())
	if rr.Body.String() == "" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), "")
	}
}
