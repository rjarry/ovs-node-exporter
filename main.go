// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2024 Robin Jarry

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rjarry/ovs-exporter/config"
	"github.com/rjarry/ovs-exporter/log"
	"github.com/rjarry/ovs-exporter/ovsdb"
)

func main() {
	conf, err := config.ParseConfig()
	if err != nil {
		// logging not initialized yet, directly write to stderr
		fmt.Fprintf(os.Stderr, "error: failed to parse config: %s\n", err)
		os.Exit(1)
	}
	err = log.InitLogging(conf.LogLevel)
	if err != nil {
		// logging not initialized yet, directly write to stderr
		fmt.Fprintf(os.Stderr, "error: failed to init log: %s\n", err)
		os.Exit(1)
	}

	log.Debugf("initializing collectors")

	var collectors []prometheus.Collector
	collectors = append(collectors, ovsdb.Collector(conf))

	for _, c := range collectors {
		log.Debugf("registering %T", c)
		err := prometheus.Register(c)
		if err != nil {
			log.Critf("collector: %s", err)
			os.Exit(1)
		}
	}

	log.Noticef("listening on http://[::]%s", conf.HttpListen)

	handler := promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			ErrorLog:            log.PrometheusLogger(),
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: 10,
			Timeout:             2 * time.Second,
			EnableOpenMetrics:   true,
		},
	)
	err = http.ListenAndServe(conf.HttpListen, handler)
	if err != nil {
		log.Critf("listen: %s", err)
		os.Exit(1)
	}
}
