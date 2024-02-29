package prome

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetricer struct {
	RequestTotalCounter *prometheus.CounterVec
	RequestSummary      *prometheus.SummaryVec
	RequestHistogram    *prometheus.HistogramVec
}

func NewPrometheusMetricer() (*PrometheusMetricer, error) {
	requestTotalCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: fmt.Sprintf("%srequests_total", DefaultMetricPrefix),
		Help: "Total requests processed, partitioned by endpoint and status",
	}, []string{
		fmt.Sprintf("%s%s", DefaultMetricLabelPrefix, MetricLabelService),
		fmt.Sprintf("%s%s", DefaultMetricLabelPrefix, MetricLabelEndpoint),
		fmt.Sprintf("%s%s", DefaultMetricLabelPrefix, MetricLabelStatus),
	})

	timeCounterSummary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: fmt.Sprintf("%slatency_microseconds", DefaultMetricPrefix),
			Help: "Request latencies in microseconds, partitioned by endpoint",
		},
		[]string{
			fmt.Sprintf("%s%s", DefaultMetricLabelPrefix, MetricLabelService),
			fmt.Sprintf("%s%s", DefaultMetricLabelPrefix, MetricLabelEndpoint),
		},
	)

	timeCounterHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: fmt.Sprintf("%srequest_duration_seconds", DefaultMetricPrefix),
			Help: "Request time in seconds, partitioned by endpoint",
		},
		[]string{
			fmt.Sprintf("%s%s", DefaultMetricLabelPrefix, MetricLabelService),
			fmt.Sprintf("%s%s", DefaultMetricLabelPrefix, MetricLabelEndpoint),
		},
	)

	for _, collector := range []prometheus.Collector{
		requestTotalCounter,
		timeCounterSummary,
		timeCounterHistogram,
	} {
		if err := prometheus.DefaultRegisterer.Register(collector); err != nil {
			// if already registered, skip fatal
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				return nil, err
			}
		}
	}
	return &PrometheusMetricer{
		RequestTotalCounter: requestTotalCounter,
		RequestSummary:      timeCounterSummary,
		RequestHistogram:    timeCounterHistogram,
	}, nil
}
