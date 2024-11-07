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
	"github.com/rjarry/ovs-exporter/appctl"
	"github.com/rjarry/ovs-exporter/config"
	"github.com/rjarry/ovs-exporter/lib"
	"github.com/rjarry/ovs-exporter/log"
	"github.com/rjarry/ovs-exporter/ovsdb"
)

func main() {
	if err := config.Parse(); err != nil {
		// logging not initialized yet, directly write to stderr
		fmt.Fprintf(os.Stderr, "error: failed to parse config: %s\n", err)
		os.Exit(1)
	}
	if err := log.InitLogging(config.LogLevel); err != nil {
		// logging not initialized yet, directly write to stderr
		fmt.Fprintf(os.Stderr, "error: failed to init log: %s\n", err)
		os.Exit(1)
	}

	log.Debugf("initializing collectors")

	var collectors []lib.Collector
	collectors = append(collectors, appctl.Collectors()...)
	collectors = append(collectors, ovsdb.Collectors()...)
	registry := prometheus.NewRegistry()

	for _, c := range collectors {
		if config.MetricSets.Has(c.MetricSet()) {
			log.Infof("registering %T", c)

			if err := registry.Register(c); err != nil {
				log.Critf("collector: %s", err)
				os.Exit(1)
			}
		} else {
			log.Infof("%T not registered, metric set not enabled", c)
		}
	}

	log.Noticef("listening on http://[::]%s", config.HttpListen)

	handler := promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			ErrorLog:            log.PrometheusLogger(),
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: 10,
			Timeout:             2 * time.Second,
			EnableOpenMetrics:   true,
		},
	)
	if err := http.ListenAndServe(config.HttpListen, handler); err != nil {
		log.Critf("listen: %s", err)
		os.Exit(1)
	}
}
