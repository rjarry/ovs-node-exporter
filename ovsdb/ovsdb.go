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
	"github.com/rjarry/ovs-exporter/schema/ovs"
)

type OvsdbCollector struct {
	endpoint string
	schema   model.ClientDBModel
	db       client.Client
}

func (c *OvsdbCollector) connect() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if c.db != nil {
		err := c.db.Echo(ctx)
		if err == nil {
			return true
		}
		log.Warningf("db.Echo: %s", err)
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
	}

	log.Noticef("connecting to ovsdb: %s", c.endpoint)

	db, err := client.NewOVSDBClient(
		c.schema,
		client.WithEndpoint(c.endpoint),
		client.WithLogger(log.Logger()),
	)
	if err != nil {
		log.Errf("NewOVSDBClient: %s", err)
		return false
	}
	if err = db.Connect(ctx); err != nil {
		log.Errf("db.Connect: %s", err)
		return false
	}
	if err = c.db.Echo(ctx); err != nil {
		log.Errf("db.Echo: %s", err)
		return false
	}

	c.db = db
	return true
}

func (c *OvsdbCollector) Describe(d chan<- *prometheus.Desc) {
}

func (c *OvsdbCollector) Collect(m chan<- prometheus.Metric) {
	if !c.connect() {
		return
	}
}

func Collector(conf *config.Config) prometheus.Collector {
	schema, err := ovs.FullDatabaseModel()
	if err != nil {
		panic(err)
	}
	return &OvsdbCollector{
		endpoint: conf.OvsdbEndpoint,
		schema:   schema,
	}
}
