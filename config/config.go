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

type Config struct {
	HttpListen    string
	OvsdbEndpoint string
	LogLevel      syslog.Priority
	MetricSets    MetricSet
}

var conf = Config{
	HttpListen:    ":1981",
	OvsdbEndpoint: "unix:/run/openvswitch/ovsdb.sock",
	LogLevel:      syslog.LOG_NOTICE,
	MetricSets:    METRICS_BASE | METRICS_ERRORS,
}

func ParseConfig() (*Config, error) {
	path, configInEnv := os.LookupEnv("OVS_NODE_EXPORTER_CONFIG")
	if !configInEnv {
		path = "/etc/ovs-exporter.conf"
	}

	file, err := ini.LoadFile(path)
	if err != nil {
		if configInEnv {
			return nil, err
		}
		file = make(ini.File)
	}

	// [main].http-listen
	value, ok := os.LookupEnv("OVS_NODE_EXPORTER_HTTP_LISTEN")
	if !ok {
		value, ok = file.Get("main", "http-listen")
	}
	if ok {
		conf.HttpListen = value
	}

	// [main].ovsdb-endpoint
	value, ok = os.LookupEnv("OVS_NODE_EXPORTER_OVSDB_ENDPOINT")
	if !ok {
		value, ok = file.Get("main", "ovsdb-endpoint")
	}
	if ok {
		conf.OvsdbEndpoint = value
	}

	// [main].log-level
	value, ok = os.LookupEnv("OVS_NODE_EXPORTER_LOG_LEVEL")
	if !ok {
		value, ok = file.Get("main", "log-level")
	}
	if ok {
		prio, err := log.ParseLogLevel(value)
		if err != nil {
			return nil, err
		}
		conf.LogLevel = prio
	}

	// [metrics].categories
	value, ok = os.LookupEnv("OVS_NODE_EXPORTER_METRICS_SETS")
	if !ok {
		value, ok = file.Get("metrics", "sets")
	}
	if ok {
		conf.MetricSets = 0
		for _, word := range strings.Fields(value) {
			switch word {
			case "base":
				conf.MetricSets |= METRICS_BASE
			case "errors":
				conf.MetricSets |= METRICS_ERRORS
			case "counters":
				conf.MetricSets |= METRICS_COUNTERS
			case "perf":
				conf.MetricSets |= METRICS_PERF
			default:
				return nil, fmt.Errorf("invalid metric set: %q", word)
			}
		}
		if conf.MetricSets == 0 {
			return nil, fmt.Errorf("no metric sets enabled")
		}
	}

	return &conf, nil
}
