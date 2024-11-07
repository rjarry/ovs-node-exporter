package ovsdb

import (
	"context"

	"github.com/ovn-org/libovsdb/ovsdb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rjarry/ovs-exporter/config"
	"github.com/rjarry/ovs-exporter/lib"
	"github.com/rjarry/ovs-exporter/log"
)

var version = lib.Metric{
	Set:         config.METRICS_BASE,
	Name:        "ovs_build_info",
	Description: "Version and library from which OVS binaries were built.",
	Labels:      []string{"ovs_version", "dpdk_version", "db_version"},
	Type:        prometheus.GaugeValue,
}

type OpenvSwitchCollector struct{}

func (c *OpenvSwitchCollector) Describe(ch chan<- *prometheus.Desc) {
	if config.MetricSets.Has(version.Set) {
		log.Debugf("%T: enabling metric %s", c, version.Name)
		ch <- version.Desc()
	}
}

func (c *OpenvSwitchCollector) Collect(ch chan<- prometheus.Metric) {
	if !connect() {
		return
	}

	results, err := DB.Transact(context.Background(), ovsdb.Operation{
		Op:    ovsdb.OperationSelect,
		Table: "Open_vSwitch",
	})
	if err != nil {
		log.Errf("transact: %s", err)
		return
	}
	for _, res := range results {
		for _, row := range res.Rows {
			if config.MetricSets.Has(version.Set) {
				ch <- prometheus.MustNewConstMetric(
					version.Desc(),
					prometheus.GaugeValue, 1,
					row["ovs_version"].(string),
					row["dpdk_version"].(string),
					row["db_version"].(string),
				)
			}
		}
	}
}

func init() {
	register(new(OpenvSwitchCollector))
}
