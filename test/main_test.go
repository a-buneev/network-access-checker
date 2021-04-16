package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/buneyev/network-access-checker/checker"
	"github.com/buneyev/network-access-checker/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
{
    "resourceList":[
        {
            "name":"Google",
            "host":"www.google.com",
            "ports":["443","80"]
        },
        {
            "name":"Yandex",
            "host":"www.yandex.com",
            "ports":["443"]
        },
        {
            "name":"for testing",
            "host":"www.213y4y.com",
            "ports":["443","80"]
        }
    ],
    "checkPeriodSeconds":3,
    "checkConnectionTimeout":2
}
*/

func TestMain(t *testing.T) {
	config := models.NewConfig("config.json")
	testConfig := models.Config{
		CheckerConfig: models.CheckerConfig{
			ResourceList: []models.Resource{
				{
					Name:  "Google",
					Host:  "www.google.com",
					Ports: []string{"443", "80"},
				},
				{
					Name:  "Yandex",
					Host:  "www.yandex.com",
					Ports: []string{"443"},
				},
				{
					Name:  "for testing",
					Host:  "www.213y4y.com",
					Ports: []string{"443", "80"},
				},
			},
			CheckPeriodSeconds:     3,
			CheckConnectionTimeout: 2,
		},
	}
	if fmt.Sprintf("%v", config.CheckerConfig) != fmt.Sprintf("%v", testConfig.CheckerConfig) {
		t.Fatal("Config parse error!")
	}

	resourceChecker := checker.NewChecker(config)
	go resourceChecker.Run()
	reload := make(chan bool)
	go config.Reload(4, reload)
	go resourceChecker.Reloader(reload)
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
