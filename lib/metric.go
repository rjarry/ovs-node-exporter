package lib

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rjarry/ovs-exporter/config"
)

type Collector interface {
	prometheus.Collector

	MetricSet() config.MetricSet
}

type Metric struct {
	Name        string
	Description string
	Labels      []string
	ConstLabels prometheus.Labels
	ValueType   prometheus.ValueType
	desc        *prometheus.Desc
}

func (m *Metric) Desc() *prometheus.Desc {
	if m.desc == nil {
		m.desc = prometheus.NewDesc(m.Name, m.Description, m.Labels, m.ConstLabels)
	}
	return m.desc
}
