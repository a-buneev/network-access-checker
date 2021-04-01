package checker

import (
	"log"
	"net"
	"time"
)

func checkConnection(host string, ports []string) error {
	for _, port := range ports {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 3*time.Second)
		if err != nil {
			log.Printf("Connecting error: %v, host: %v, port: %v", err.Error(), host, port)
			return err
		}
		if conn != nil {
			defer conn.Close()
		}
	}
	return nil
}
