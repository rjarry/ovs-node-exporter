// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2024 Robin Jarry

package ovsdb

import (
	"context"
	"time"

	"github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/libovsdb/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rjarry/ovs-exporter/config"
	"github.com/rjarry/ovs-exporter/log"
)

var (
	// initialized via register() in modules init()
	collectors []prometheus.Collector
	// initialized once in Collectors()
	schema model.ClientDBModel
	// nil at startup, initialized with connect()
	DB client.Client
)

func register(c prometheus.Collector) {
	collectors = append(collectors, c)
}

func connect() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if DB != nil {
		err := DB.Echo(ctx)
		if err == nil {
			return true
		}
		log.Warningf("db.Echo: %s", err)
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
	}

	log.Noticef("connecting to ovsdb: %s", config.OvsdbEndpoint)

	db, err := client.NewOVSDBClient(
		schema,
		client.WithEndpoint(config.OvsdbEndpoint),
		client.WithLogger(log.OvsdbLogger()),
	)
	if err != nil {
		log.Errf("NewOVSDBClient: %s", err)
		return false
	}
	if err = db.Connect(ctx); err != nil {
		log.Errf("db.Connect: %s", err)
		return false
	}
	if err = db.Echo(ctx); err != nil {
		log.Errf("db.Echo: %s", err)
		return false
	}

	DB = db
	return true
}

func Collectors() []prometheus.Collector {
	schema, _ = model.NewClientDBModel("Open_vSwitch", nil)
	return collectors
}
