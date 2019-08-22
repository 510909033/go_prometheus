package main

import (
	"math/rand"
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
	var opts = prometheus.CounterOpts{
		Name: "test_inc",
		Help: "inc数",
	}
	incCounter := prometheus.NewCounter(opts)
	prometheus.MustRegister(incCounter)

	var optsGauge = prometheus.GaugeOpts{
		Name: "test_gauge_set",
		Help: "gauge_set数",
	}

	gauge := prometheus.NewGauge(optsGauge)
	prometheus.MustRegister(gauge)

	histogramOpts := prometheus.HistogramOpts{
		Name:    "test_histogram",
		Help:    "test_histogram数",
		Buckets: []float64{500, 1000, 2000, 4000, 8000},
	}
	histogram := prometheus.NewHistogram(histogramOpts)
	prometheus.MustRegister(histogram)

	for {
		incCounter.Inc()
		gauge.Set(float64(rand.Int63n(100000)))
		histogram.Observe(rand.Float64() * 3128)
		time.Sleep(time.Second)
	}
}
