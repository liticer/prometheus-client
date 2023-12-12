// Copyright 2022 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build go1.17
// +build go1.17

// A minimal example of how to include Prometheus instrumentation.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()

	// Create a new registry.
	reg := prometheus.NewRegistry()

	// Register metrics from GoCollector collecting statistics from the Go Runtime.
	// This enabled default, recommended metrics with the additional, recommended metric for
	// goroutine scheduling latencies histogram that is currently bit too expensive for default option.
	//
	// See the related GopherConUK talk to learn more: https://www.youtube.com/watch?v=18dyI_8VFa0
	reg.MustRegister(
		collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(
				collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/sched/latencies:seconds")},
			),
		),
	)

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			OpenMetricsOptions: promhttp.OpenMetricsOptions{
				// Opt into OpenMetrics to support exemplars.
				Enable: true,
			},
		},
	))
	fmt.Println("Hello world from new Go Collector!")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
