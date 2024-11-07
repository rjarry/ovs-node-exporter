// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2024 Robin Jarry

package appctl

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rjarry/ovs-exporter/config"
	"github.com/rjarry/ovs-exporter/lib"
	"github.com/rjarry/ovs-exporter/log"
)

var pmdStatsMetrics = map[string]lib.Metric{
	"packets received": {
		Name:        "ovs_pmd_packets_received",
		Description: "packets received",
		ValueType:   prometheus.CounterValue,
		Labels:      []string{"cpu", "numa"},
	},
	"packet recirculations": {
		Name:        "ovs_pmd_packet_recirculations",
		Description: "packet recirculations",
		ValueType:   prometheus.CounterValue,
		Labels:      []string{"cpu", "numa"},
	},
	"avg. datapath passes per packet": {
		Name:        "ovs_pmd_avg_dp_passes_per_packet",
		Description: "avg. datapath passes per packet",
		ValueType:   prometheus.GaugeValue,
		Labels:      []string{"cpu", "numa"},
	},
	"phwol hits": {
		Name:        "ovs_pmd_phwol_hits",
		Description: "phwol hits",
		ValueType:   prometheus.CounterValue,
		Labels:      []string{"cpu", "numa"},
	},
	"mfex opt hits": {
		Name:        "ovs_pmd_mfex_opt_hits",
		Description: "mfex opt hits",
		ValueType:   prometheus.CounterValue,
		Labels:      []string{"cpu", "numa"},
	},
	"simple match hits": {
		Name:        "ovs_pmd_simple_match_hits",
		Description: "simple match hits",
		ValueType:   prometheus.CounterValue,
		Labels:      []string{"cpu", "numa"},
	},
	"emc hits": {
		Name:        "ovs_pmd_emc_hits",
		Description: "emc hits",
		ValueType:   prometheus.CounterValue,
		Labels:      []string{"cpu", "numa"},
	},
	"smc hits": {
		Name:        "ovs_pmd_smc_hits",
		Description: "smc hits",
		ValueType:   prometheus.CounterValue,
		Labels:      []string{"cpu", "numa"},
	},
	"megaflow hits": {
		Name:        "ovs_pmd_megaflow_hits",
		Description: "megaflow hits",
		ValueType:   prometheus.CounterValue,
		Labels:      []string{"cpu", "numa"},
	},
	"avg. subtable lookups per megaflow hit": {
		Name:        "ovs_pmd_avg_subtable_lookups_per_megaflow_hit",
		Description: "avg. subtable lookups per megaflow hit",
		ValueType:   prometheus.GaugeValue,
		Labels:      []string{"cpu", "numa"},
	},
	"miss with success upcall": {
		Name:        "ovs_pmd_miss_with_success_upcall",
		Description: "miss with success upcall",
		ValueType:   prometheus.CounterValue,
		Labels:      []string{"cpu", "numa"},
	},
	"miss with failed upcall": {
		Name:        "ovs_pmd_miss_with_failed_upcall",
		Description: "miss with failed upcall",
		ValueType:   prometheus.CounterValue,
		Labels:      []string{"cpu", "numa"},
	},
	"avg. packets per output batch": {
		Name:        "ovs_avg_packets_per_output_batch",
		Description: "avg. packets per output batch",
		ValueType:   prometheus.GaugeValue,
		Labels:      []string{"cpu", "numa"},
	},
	"idle cycles": {
		Name:        "ovs_pmd_idle_cycles",
		Description: "idle cycles",
		ValueType:   prometheus.CounterValue,
		Labels:      []string{"cpu", "numa"},
	},
	"processing cycles": {
		Name:        "ovs_pmd_processing_cycles",
		Description: "processing cycles",
		ValueType:   prometheus.CounterValue,
		Labels:      []string{"cpu", "numa"},
	},
	"avg cycles per packet": {
		Name:        "ovs_pmd_avg_cycles_per_packet",
		Description: "avg cycles per packet",
		ValueType:   prometheus.GaugeValue,
		Labels:      []string{"cpu", "numa"},
	},
	"avg processing cycles per packet": {
		Name:        "ovs_pmd_avg_processing_cycles_per_packet",
		Description: "avg processing cycles per packet",
		ValueType:   prometheus.GaugeValue,
		Labels:      []string{"cpu", "numa"},
	},
}

func makeMetric(cpu, numa, name, value string) prometheus.Metric {
	m, ok := pmdStatsMetrics[name]
	if !ok {
		return nil
	}

	val, err := strconv.ParseFloat(strings.Fields(value)[0], 64)
	if err != nil {
		log.Errf("%s: %s: %s", name, value, err)
		return nil
	}

	return prometheus.MustNewConstMetric(m.Desc(), m.ValueType, val, cpu, numa)
}

type PmdRxqCollector struct{}

func (PmdRxqCollector) MetricSet() config.MetricSet {
	return config.METRICS_PERF
}

func (c *PmdRxqCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range pmdStatsMetrics {
		log.Debugf("%T: enabling metric %s", c, m.Name)
		ch <- m.Desc()
	}
}

var (
	// "pmd thread numa_id 0 core_id 39:"
	pmdThreadRe = regexp.MustCompile(`(?m)^pmd thread numa_id (\d+) core_id (\d+):$`)
	// " idle cycles: 28600864377118 (100.00%)"
	statRe = regexp.MustCompile(`(?m)^\s+(.+): (.+)$`)
)

func (c *PmdRxqCollector) Collect(ch chan<- prometheus.Metric) {
	buf := call("dpif-netdev/pmd-stats-show")
	if buf == "" {
		return
	}

	cpu := ""
	numa := ""

	scanner := bufio.NewScanner(strings.NewReader(buf))
	for scanner.Scan() {
		line := scanner.Text()

		if cpu != "" && numa != "" {
			match := statRe.FindStringSubmatch(line)
			if match != nil {
				metric := makeMetric(cpu, numa, match[1], match[2])
				if metric != nil {
					ch <- metric
				}
				continue
			}
			cpu = ""
			numa = ""
		}

		match := pmdThreadRe.FindStringSubmatch(line)
		if match != nil {
			numa = match[1]
			cpu = match[2]
		}
	}
}

func init() {
	register(new(PmdRxqCollector))
}
