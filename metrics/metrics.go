// Package metrics provides metrics interfaces sourced from github.com/go-kit/kit/metrics
package metrics

// Counter describes a metric that accumulates values monotonically.
type Counter interface {
	With(labelValues ...string) Counter
	Add(delta float64)
}

// Gauge describes a metric that takes specific values over time.
type Gauge interface {
	With(labelValues ...string) Gauge
	Set(value float64)
	Add(delta float64)
}

// Histogram describes a metric that takes repeated observations of the same
// kind of thing, and produces a statistical summary of those observations,
// typically expressed as quantiles or buckets.
type Histogram interface {
	With(labelValues ...string) Histogram
	Observe(value float64)
}
