package lib

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rjarry/ovs-exporter/config"
)

type Metric struct {
	Set         config.MetricSet
	Name        string
	Description string
	Labels      []string
	ConstLabels prometheus.Labels
	Type        prometheus.ValueType
	desc        *prometheus.Desc
}

func (m *Metric) Desc() *prometheus.Desc {
	if m.desc == nil {
		m.desc = prometheus.NewDesc(m.Name, m.Description, m.Labels, m.ConstLabels)
	}
	return m.desc
}
