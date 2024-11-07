// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2024 Robin Jarry

package ovsdb

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rjarry/ovs-exporter/config"
	"github.com/rjarry/ovs-exporter/lib"
	"github.com/rjarry/ovs-exporter/log"
)

var vswitchMetrics = []lib.Metric{
	{
		Name:        "ovs_build_info",
		Description: "Version and library from which OVS binaries were built.",
		Labels:      []string{"ovs_version", "dpdk_version", "db_version"},
		ValueType:   prometheus.GaugeValue,
	},
}

type OpenvSwitchCollector struct{}

func (OpenvSwitchCollector) MetricSet() config.MetricSet {
	return config.METRICS_BASE
}

func (c *OpenvSwitchCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range vswitchMetrics {
		log.Debugf("%T: enabling metric %s", c, m.Name)
		ch <- m.Desc()
	}
}

func (c *OpenvSwitchCollector) Collect(ch chan<- prometheus.Metric) {
	rows := query("Open_vSwitch")
	if rows == nil {
		return
	}
	row := rows[0]

	for _, m := range vswitchMetrics {
		labels := make([]string, 0, len(m.Labels))
		for _, name := range m.Labels {
			labels = append(labels, row[name].(string))
		}
		ch <- prometheus.MustNewConstMetric(m.Desc(), m.ValueType, 1, labels...)
	}
}

func init() {
	register(new(OpenvSwitchCollector))
}
