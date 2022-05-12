package gp

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricType int

const (
	None MetricType = iota
	Counter
	Gauge
	Histogram
	Summary
)

type Metric struct {
	Type        MetricType
	Name        string
	Description string
	Labels      []string
	Buckets     []float64
	Objectives  map[float64]float64

	vec prometheus.Collector
}

func (m *Metric) SetValue(labelValues []string, value float64) error {
	if m.Type == None {
		return errors.Errorf("metric '%s' not existed.", m.Name)
	}

	if m.Type != Gauge {
		return errors.Errorf("metric '%s' not Gauge type", m.Name)
	}
	m.vec.(*prometheus.GaugeVec).WithLabelValues(labelValues...).Set(value)
	return nil
}

func (m *Metric) IncBy(labelValues []string, value float64) error {
	if m.Type == None {
		return errors.Errorf("metric '%s' not existed.", m.Name)
	}

	if m.Type != Counter {
		return errors.Errorf("metric '%s' not Counter type", m.Name)
	}
	m.vec.(*prometheus.CounterVec).WithLabelValues(labelValues...).Add(value)
	return nil
}

func (m *Metric) Observe(labelValues []string, value float64) error {
	if m.Type == 0 {
		return errors.Errorf("metric '%s' not existed.", m.Name)
	}
	if m.Type != Histogram && m.Type != Summary {
		return errors.Errorf("metric '%s' not Histogram or Summary type", m.Name)
	}
	switch m.Type {
	case Histogram:
		m.vec.(*prometheus.HistogramVec).WithLabelValues(labelValues...).Observe(value)
		break
	case Summary:
		m.vec.(*prometheus.SummaryVec).WithLabelValues(labelValues...).Observe(value)
		break
	}
	return nil
}
