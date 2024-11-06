package ovsdb

import (
	"github.com/ovn-org/libovsdb/client"
	"github.com/prometheus/client_golang/prometheus"
)

type OvsMetric struct {
	Desc  *prometheus.Desc
	Value func(ovs client.Client) prometheus.Metric
}
