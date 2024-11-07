# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2024 Robin Jarry

version = $(shell git describe --long --abbrev=12 --tags --dirty 2>/dev/null || echo 0.1)
src = $(shell find * -type f -name '*.go') go.mod go.sum
go_ldflags :=
go_ldflags += -X main.version=$(version)

.PHONY: all
all: ovs-node-exporter

ovs-node-exporter: $(src)
	go build -trimpath -ldflags='$(go_ldflags)' -o $@

.PHONY: debug
debug: ovs-node-exporter.debug

ovs-node-exporter.debug: $(src)
	go build -ldflags='$(go_ldflags)' -gcflags=all="-N -l" -o $@

.PHONY: format
format:
	gofmt -w .

.PHONY: lint
lint:
	@gofmt -d . | grep ^ \
	&& echo The above files need to be formatted. Please run make format. && exit 1 \
	|| echo All files formated.

.PHONY: run
run: ovs-node-exporter
	OVS_NODE_EXPORTER=etc/dev.conf ./$<
