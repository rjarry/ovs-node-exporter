// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2024 Robin Jarry

package ovsdb

import (
	"context"
	"time"

	"github.com/ovn-org/libovsdb/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rjarry/ovs-exporter/config"
)

var RejectedCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "srht",
	Subsystem: "lists",
	Name:      "conn_rejected",
	Help:      "Total number of rejected connections or messages.",
})

var DroppedCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "srht",
	Subsystem: "lists",
	Name:      "emails_dropped",
	Help:      "Total number of silently dropped messages.",
})

var EmailsCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "srht",
	Subsystem: "lists",
	Name:      "emails_received",
	Help:      "Total number of emails received.",
})

func Collectors(conf *config.Config) []prometheus.Collector {
	c, err := client.NewOVSDBClient(schema, client.WithEndpoint(endpoint))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = c.Connect(ctx)
	if err != nil {
		return err
	}

	return nil
}
