// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2024 Robin Jarry

package config

import (
	"fmt"
	"log/syslog"
	"os"
	"strings"

	"github.com/rjarry/ovs-exporter/log"
	"github.com/vaughan0/go-ini"
)

type MetricSet uint32

const (
	METRICS_BASE MetricSet = 1 << iota
	METRICS_ERRORS
	METRICS_COUNTERS
	METRICS_PERF
)

func (s MetricSet) Has(o MetricSet) bool {
	return s&o == o
}

var (
	// [main].http-listen or OVS_NODE_EXPORTER_HTTP_LISTEN
	HttpListen string = ":1981"

	// [main].ovsdb-endpoint or OVS_NODE_EXPORTER_OVSDB_ENDPOINT
	OvsdbEndpoint string = "unix:/run/openvswitch/db.sock"

	// [main].log-level or OVS_NODE_EXPORTER_LOG_LEVEL
	LogLevel syslog.Priority = syslog.LOG_NOTICE

	// [metrics].sets or OVS_NODE_EXPORTER_METRICS_SETS
	MetricSets MetricSet = METRICS_BASE | METRICS_ERRORS
)

func Parse() error {
	path, configInEnv := os.LookupEnv("OVS_NODE_EXPORTER_CONFIG")
	if !configInEnv {
		path = "/etc/ovs-node-exporter.conf"
	}

	file, err := ini.LoadFile(path)
	if err != nil {
		if configInEnv {
			return err
		}
		file = make(ini.File)
	}

	// [main].http-listen
	value, ok := os.LookupEnv("OVS_NODE_EXPORTER_HTTP_LISTEN")
	if !ok {
		value, ok = file.Get("main", "http-listen")
	}
	if ok {
		HttpListen = value
	}

	// [main].ovsdb-endpoint
	value, ok = os.LookupEnv("OVS_NODE_EXPORTER_OVSDB_ENDPOINT")
	if !ok {
		value, ok = file.Get("main", "ovsdb-endpoint")
	}
	if ok {
		OvsdbEndpoint = value
	}

	// [main].log-level
	value, ok = os.LookupEnv("OVS_NODE_EXPORTER_LOG_LEVEL")
	if !ok {
		value, ok = file.Get("main", "log-level")
	}
	if ok {
		prio, err := log.ParseLogLevel(value)
		if err != nil {
			return err
		}
		LogLevel = prio
	}

	// [metrics].categories
	value, ok = os.LookupEnv("OVS_NODE_EXPORTER_METRICS_SETS")
	if !ok {
		value, ok = file.Get("metrics", "sets")
	}
	if ok {
		MetricSets = 0
		for _, word := range strings.Fields(value) {
			switch word {
			case "base":
				MetricSets |= METRICS_BASE
			case "errors":
				MetricSets |= METRICS_ERRORS
			case "counters":
				MetricSets |= METRICS_COUNTERS
			case "perf":
				MetricSets |= METRICS_PERF
			default:
				return fmt.Errorf("invalid metric set: %q", word)
			}
		}
		if MetricSets == 0 {
			return fmt.Errorf("no metric sets enabled")
		}
	}

	return nil
}
