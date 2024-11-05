package main

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
	path, ok := os.LookupEnv("OVS_EXPORTER_CONFIG")
	if !ok {
		path = "/etc/ovs-exporter.conf"
	}

	file, err := ini.LoadFile(path)
	if err != nil {
		return nil, err
	}

	if v, ok := file.Get("main", "http-endpoint"); ok {
		conf.HttpEndpoint = v
	}
	if v, ok := file.Get("main", "log-level"); ok {
		prio, err := log.ParseLogLevel(v)
		if err != nil {
			return nil, err
		}
		conf.LogLevel = prio
	}

	return &conf, nil
}
