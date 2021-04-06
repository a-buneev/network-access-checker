package checker

import (
	"log"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func (chckr *Checker) checkConnection(name, host string, ports []string, timeout time.Duration) error {
	for _, port := range ports {
		timer := prometheus.NewTimer(chckr.reqDuration.WithLabelValues(name, host, port))
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			chckr.gaugeOpts.WithLabelValues(name, host, port).Set(0)
			log.Printf("Connecting error: %v, host: %v, port: %v", err.Error(), host, port)
			return err
		}
		timer.ObserveDuration()
		if conn != nil {
			defer conn.Close()
		}
		chckr.gaugeOpts.WithLabelValues(name, host, port).Set(1)
	}
	return nil
}
