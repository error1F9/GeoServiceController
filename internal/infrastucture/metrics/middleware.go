package metrics

import (
	"net/http"
	"time"
)

func (m *Metrics) DurationAndCounterMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			duration := time.Since(start).Seconds()
			m.requestDuration.WithLabelValues(r.URL.Path).Observe(duration)
			m.endPointCounter.WithLabelValues(r.URL.Path).Inc()
		}()
		next.ServeHTTP(w, r)
	})
}
