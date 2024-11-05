package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	RejectedCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "srht",
		Subsystem: "lists",
		Name:      "conn_rejected",
		Help:      "Total number of rejected connections or messages.",
	})

	DroppedCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "srht",
		Subsystem: "lists",
		Name:      "emails_dropped",
		Help:      "Total number of silently dropped messages.",
	})

	EmailsCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "srht",
		Subsystem: "lists",
		Name:      "emails_received",
		Help:      "Total number of emails received.",
	})

	ErrorsCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "srht",
		Subsystem: "lists",
		Name:      "email_errors",
		Help:      "Total number of erroneous emails received.",
	})

	BounceCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "srht",
		Subsystem: "lists",
		Name:      "email_bounced",
		Help:      "Total number of bounced emails.",
	})

	ForwardsCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "srht",
		Subsystem: "lists",
		Name:      "forwards_processed",
		Help:      "Total number of emails forwarded to subscribers.",
	})

	CommandsCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "srht",
		Subsystem: "lists",
		Name:      "commands_processed",
		Help:      "Total number of commands processed, e.g. +subscribe.",
	})

	prometheus.MustRegister(
		RejectedCounter,
		DroppedCounter,
		EmailsCounter,
		ErrorsCounter,
		BounceCounter,
		ForwardsCounter,
		CommandsCounter,
	)
	log.Println("exposing prometheus metrics over http://[::]:1234")
	err := http.ListenAndServe(":1234", promhttp.Handler())
	if err != nil {
		log.Fatal(err)
	}
}
