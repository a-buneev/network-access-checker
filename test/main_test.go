package test

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
	config := models.NewConfig("config.json")
	resourceChecker := checker.NewChecker(config)
	go resourceChecker.Run()
	go config.Reload(4)
	time.Sleep(10 * time.Second)
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

	assertGauge := map[string]string{
		"www.213y4y.com:443": "monitoring_network_access_status{resource_addr=\"www.213y4y.com\",resource_name=\"for testing\",resource_port=\"443\"} 0",
		"www.google.com:443": "monitoring_network_access_status{resource_addr=\"www.google.com\",resource_name=\"Google\",resource_port=\"443\"} 1",
		"www.google.com:80":  "monitoring_network_access_status{resource_addr=\"www.google.com\",resource_name=\"Google\",resource_port=\"80\"} 1",
		"www.yandex.com:443": "monitoring_network_access_status{resource_addr=\"www.yandex.com\",resource_name=\"Yandex\",resource_port=\"443\"} 1",
	}

	//fmt.Println(rr.Body.String())
	for addr, result := range assertGauge {
		if !strings.Contains(rr.Body.String(), result) {
			t.Errorf("handler returned unexpected body: %v not found", addr)
		}
	}

	assertDuration := map[string]int{
		"monitoring_network_req_duration_bucket": 48,
		"monitoring_network_req_duration_sum":    4,
		"monitoring_network_req_duration_count":  4,
	}

	for metricName, result := range assertDuration {
		if strings.Count(rr.Body.String(), metricName) != result {
			t.Errorf("handler returned unexpected body: %v == %v", result, strings.Count(rr.Body.String(), metricName))
		}
	}

}
