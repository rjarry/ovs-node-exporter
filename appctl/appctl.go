// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2024 Robin Jarry

package appctl

import (
	client "github.com/influxdata/influxdb1-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rjarry/ovs-exporter/config"
)

func Collectors(conf *config.Config) []prometheus.Collector {
	client.Client()

	return nil
}
