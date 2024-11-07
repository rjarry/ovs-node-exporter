// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2024 Robin Jarry

package ovsdb

import (
	"context"
	"time"

	"github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/libovsdb/model"
	"github.com/ovn-org/libovsdb/ovsdb"
	"github.com/rjarry/ovs-exporter/config"
	"github.com/rjarry/ovs-exporter/lib"
	"github.com/rjarry/ovs-exporter/log"
)

var (
	// initialized via register() in modules init()
	collectors []lib.Collector
	// initialized once in Collectors()
	schema model.ClientDBModel
)

func register(c lib.Collector) {
	collectors = append(collectors, c)
}

func query(table string) []ovsdb.Row {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	log.Debugf("connecting to ovsdb: %s", config.OvsdbEndpoint)

	db, err := client.NewOVSDBClient(
		schema,
		client.WithEndpoint(config.OvsdbEndpoint),
		client.WithLogger(log.OvsdbLogger()),
	)
	if err != nil {
		log.Errf("NewOVSDBClient: %s", err)
		return nil
	}
	if err = db.Connect(ctx); err != nil {
		log.Errf("db.Connect: %s", err)
		return nil
	}
	defer db.Disconnect()

	results, err := db.Transact(ctx, ovsdb.Operation{
		Op:    ovsdb.OperationSelect,
		Table: table,
	})
	if err != nil {
		log.Errf("transact: %s", err)
		return nil
	}

	for _, result := range results {
		return result.Rows
	}

	return nil
}

func Collectors() []lib.Collector {
	schema, _ = model.NewClientDBModel("Open_vSwitch", nil)
	return collectors
}
