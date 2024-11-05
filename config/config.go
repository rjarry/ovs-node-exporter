// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2024 Robin Jarry

package config

import (
	"log/syslog"
	"os"

	"github.com/rjarry/ovs-exporter/log"
	"github.com/vaughan0/go-ini"
)

type Config struct {
	HttpEndpoint string
	LogLevel     syslog.Priority
}

var conf = Config{
	HttpEndpoint: ":1981",
	LogLevel:     syslog.LOG_NOTICE,
}

func ParseConfig() (*Config, error) {
	path, configInEnv := os.LookupEnv("OVS_EXPORTER_CONFIG")
	if !configInEnv {
		path = "/etc/ovs-exporter.conf"
	}

	file, err := ini.LoadFile(path)
	if err != nil && configInEnv {
		return nil, err
	}

	value, ok := os.LookupEnv("OVS_EXPORTER_HTTP_ENDPOINT")
	if !ok {
		value, ok = file.Get("main", "http-endpoint")
	}
	if ok {
		conf.HttpEndpoint = value
	}
	value, ok = os.LookupEnv("OVS_EXPORTER_LOG_LEVEL")
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

	return &conf, nil
}
