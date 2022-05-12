package gp

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

var instance *GP

func GetGP() *GP {
	if instance == nil {
		instance = &GP{metrics: map[string]*Metric{}}
	}
	return instance
}

type GP struct {
	metrics map[string]*Metric
}

func (gp *GP) Metric(name string) (*Metric, bool) {
	if m, exist := gp.metrics[name]; exist {
		return m, exist
	}
	return nil, false
}

func (gp *GP) RegisterMetrics(ms []*Metric) error {
	var err error
	for _, m := range ms {
		if err = gp.RegisterMetric(m); err != nil {
			return err
		}
	}
	return nil
}

func (gp *GP) RegisterMetric(m *Metric) error {
	if _, ok := gp.metrics[m.Name]; ok {
		return errors.Errorf("metric '%s' is existed", m.Name)
	}
	if m.Name == "" {
		return errors.Errorf("metric name cannot be empty.")
	}
	gp.metrics[m.Name] = m

	// init Metric.vec
	var vec prometheus.Collector
	switch m.Type {
	case Counter:
		vec = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: m.Name,
				Help: m.Description,
			},
			m.Labels,
		)
	case Gauge:
		vec = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: m.Name,
				Help: m.Description,
			},
			m.Labels,
		)
	case Histogram:
		vec = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    m.Name,
				Help:    m.Description,
				Buckets: m.Buckets,
			},
			m.Labels,
		)
	case Summary:
		vec = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       m.Name,
				Help:       m.Description,
				Objectives: m.Objectives,
			},
			m.Labels,
		)
	}

	if err := prometheus.Register(vec); err != nil {
		return errors.Errorf("%s could not be registered in Prometheus", m.Name)
	}

	m.vec = vec

	return nil
}
