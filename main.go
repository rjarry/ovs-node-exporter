package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rjarry/ovs-exporter/appctl"
	"github.com/rjarry/ovs-exporter/log"
	"github.com/rjarry/ovs-exporter/ovsdb"
)

func die(message string, args ...any) {
	log.Critf(message, args...)
	os.Exit(1)
}

func main() {
	config, err := ParseConfig()
	if err != nil {
		// logging not initialized yet, directly write to stderr
		fmt.Fprintf(os.Stderr, "error: failed to parse config: %s\n", err)
		os.Exit(1)
	}
	err = log.InitLogging(config.LogLevel)
	if err != nil {
		// logging not initialized yet, directly write to stderr
		fmt.Fprintf(os.Stderr, "error: failed to init log: %s\n", err)
		os.Exit(1)
	}

	log.Debugf("initializing collectors")

	var collectors []prometheus.Collector

	collectors = append(collectors, ovsdb.Collectors()...)

	for _, c := range appctl.Collectors() {
		log.Debugf("registering %v", c)
		err := prometheus.Register(c)
		if err != nil {
			die("collector: %s", err)
		}
	}

	log.Infof("listening on http://[::]%s", config.HttpEndpoint)

	err = http.ListenAndServe(config.HttpEndpoint, promhttp.Handler())
	if err != nil {
		die("listen: %s", err)
	}
}
