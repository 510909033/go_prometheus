package demo

import (
	"math/rand"
	"time"

	//      "github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus"
)

var suffix = "3"

func Pcounter() {
	var opts = prometheus.CounterOpts{
		Name: "test_inc" + suffix,
		Help: "inc数-" + suffix,
	}
	incCounter := prometheus.NewCounter(opts)
	prometheus.MustRegister(incCounter)
	go func() {
		for {
			incCounter.Add(float64(1))
			time.Sleep(time.Second * 3)
		}
	}()
}

func Pguage() {
	var optsGauge = prometheus.GaugeOpts{
		Name: "test_gauge_set",
		Help: "gauge_set数",
	}
	gauge := prometheus.NewGauge(optsGauge)
	prometheus.MustRegister(gauge)
	go func() {
		for {
			gauge.Set(float64(rand.Int63n(10000)))
			time.Sleep(time.Second * 3)
		}

	}()

}

func Phistogram() {
	histogramOpts := prometheus.HistogramOpts{
		Name:    "test_histogram",
		Help:    "test_histogram数",
		Buckets: []float64{500, 1000, 2000, 4000, 8000},
	}
	histogram := prometheus.NewHistogram(histogramOpts)
	prometheus.MustRegister(histogram)

	go func() {
		for {
			histogram.Observe(rand.Float64() * 3128)
			time.Sleep(time.Second * 3)
		}
	}()

}
