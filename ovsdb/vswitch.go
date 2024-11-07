package ovsdb

import (
	"context"

	"github.com/ovn-org/libovsdb/ovsdb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rjarry/ovs-exporter/config"
	"github.com/rjarry/ovs-exporter/lib"
	"github.com/rjarry/ovs-exporter/log"
)

var metrics = []lib.Metric{
	{
		Set:         config.METRICS_BASE,
		Name:        "ovs_build_info",
		Description: "Version and library from which OVS binaries were built.",
		Labels:      []string{"ovs_version", "dpdk_version", "db_version"},
		ValueType:   prometheus.GaugeValue,
	},
}

type OpenvSwitchCollector struct{}

func (c *OpenvSwitchCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range metrics {
		if config.MetricSets.Has(m.Set) {
			log.Debugf("%T: enabling metric %s", c, m.Name)
			ch <- m.Desc()
		}
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
	row := results[0].Rows[0]

	for _, m := range metrics {
		if config.MetricSets.Has(m.Set) {
			labels := make([]string, 0, len(m.Labels))
			for _, name := range m.Labels {
				labels = append(labels, row[name].(string))
			}
			ch <- prometheus.MustNewConstMetric(m.Desc(), m.ValueType, 1, labels...)
		}
	}
}

func init() {
	register(new(OpenvSwitchCollector))
}
