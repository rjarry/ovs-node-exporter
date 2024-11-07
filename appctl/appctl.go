// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2024 Robin Jarry

package appctl

import (
	"fmt"
	"io"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"github.com/rjarry/ovs-exporter/config"
	"github.com/rjarry/ovs-exporter/lib"
	"github.com/rjarry/ovs-exporter/log"
)

func call(method string, args ...string) string {
	var sockpath string
	var pid int
	var err error
	var buf []byte
	var f *os.File
	var reply string

	if f, err = os.Open(config.AppctlPidfile); err != nil {
		log.Errf("os.Open: %s", err)
		return ""
	}
	if buf, err = io.ReadAll(f); err != nil {
		log.Errf("io.ReadAll: %s", err)
		return ""
	}
	if pid, err = strconv.Atoi(strings.TrimSpace(string(buf))); err == nil {
		sockpath = strings.TrimSuffix(config.AppctlPidfile, ".pid")
		sockpath = fmt.Sprintf("%s.%d.ctl", sockpath, pid)
	}

	conn, err := net.Dial("unix", sockpath)
	if err != nil {
		log.Errf("net.Dial: %s", err)
		return ""
	}

	client := rpc.NewClientWithCodec(NewClientCodec(conn))
	defer func() {
		err := client.Close()
		if err != nil {
			log.Warningf("close: %s", err)
		}
	}()

	if args == nil {
		args = make([]string, 0)
	}

	log.Debugf("calling: %s %s", method, args)
	if err := client.Call(method, args, &reply); err != nil {
		log.Errf("call(%s): %s", method, err)
		return ""
	}

	return reply
}

var collectors []lib.Collector

func register(c lib.Collector) {
	collectors = append(collectors, c)
}

func Collectors() []lib.Collector {
	return collectors
}
