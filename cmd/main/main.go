package main

import (
	"time"
	//      "github.com/prometheus/client_golang/prometheus"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	AddPrometheus()
	select {}
}

func AddPrometheus() {
	go func() {
		server := http.NewServeMux()
		server.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9101", server)
	}()

	inc()

}

func inc() {
	incCounter := prometheus.NewCounter("test_inc", "inc数")
	for {
		incCounter.Inc()
		time.Sleep(time.Second)
	}
}