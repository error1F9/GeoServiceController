package metrics

import (
	"GeoService/internal/modules/address/entity"
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var PrometheusMetrics = NewMetrics(prometheus.DefaultRegisterer)

type Metrics struct {
	endPointCounter *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	cacheDuration   *prometheus.HistogramVec
	apiDuration     *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		endPointCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "total_reqs",
				Help: "Total number of requests",
			},
			[]string{"endpoint"},
		),
		requestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "request_duration_seconds",
			Help: "Time for request",
		},
			[]string{"endpoint"},
		),
		cacheDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "request_to_cache",
			Help: "Time for request to cache",
		},
			[]string{"method"},
		),
		apiDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "request_to_api",
			Help: "Time for request to api",
		},
			[]string{"method"},
		),
	}

	reg.MustRegister(m.endPointCounter)
	reg.MustRegister(m.requestDuration)
	reg.MustRegister(m.cacheDuration)
	reg.MustRegister(m.apiDuration)
	return m
}

func (m *Metrics) MethodRequestDuration(methodName string, fn func(ctx context.Context, params ...string) ([]*entity.Address, error, bool)) func(ctx context.Context, params ...string) ([]*entity.Address, error) {
	return func(ctx context.Context, params ...string) ([]*entity.Address, error) {
		start := time.Now()

		result, err, fromCache := fn(ctx, params...)

		duration := time.Since(start).Seconds()
		if !fromCache {
			m.apiDuration.WithLabelValues(methodName).Observe(duration)
		} else {
			m.cacheDuration.WithLabelValues(methodName).Observe(duration)
		}
		return result, err
	}
}
