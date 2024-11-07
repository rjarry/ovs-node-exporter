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

type PmdRxqCollector struct{}

var pmdStatsMetrics = map[string]lib.Metric{
	"idle cycles": {
		Set:         config.METRICS_PERF,
		Name:        "ovs_pmd_cpu_idle_cycles",
		Description: "Number of CPU cycles spent doing empty polls.",
		Labels:      []string{"cpu", "numa"},
		ValueType:   prometheus.CounterValue,
	},
	"processing cycles": {
		Set:         config.METRICS_PERF,
		Name:        "ovs_pmd_cpu_busy_cycles",
		Description: "Number of CPU cycles spent processing packets.",
		Labels:      []string{"cpu", "numa"},
		ValueType:   prometheus.CounterValue,
	},
}

func makeMetric(cpu, numa, name, value string) prometheus.Metric {
	m, ok := pmdStatsMetrics[name]

	if !ok || !config.MetricSets.Has(m.Set) {
		return nil
	}

	var val float64

	switch name {
	case "idle cycles", "processing cycles":
		val, _ = strconv.ParseFloat(strings.Fields(value)[0], 64)
	}

	return prometheus.MustNewConstMetric(m.Desc(), m.ValueType, val, cpu, numa)
}

func (c *PmdRxqCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range pmdStatsMetrics {
		if config.MetricSets.Has(m.Set) {
			log.Debugf("%T: enabling metric %s", c, m.Name)
			ch <- m.Desc()
		}
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
