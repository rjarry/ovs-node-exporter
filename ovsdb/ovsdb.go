package ovsdb

import "github.com/prometheus/client_golang/prometheus"

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

func Collectors() []prometheus.Collector {
	return nil
}
