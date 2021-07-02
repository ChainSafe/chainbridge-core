// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// ChainMetrics is a public struct that includes data related to transfers occuring over the chainbridge
type ChainMetrics struct {
	// Total amount of tokens that have been transferred
	AmountTransferred prometheus.Counter
	// Total number of transfers that have occurred
	NumberOfTransfers prometheus.Counter
}

// NewChainMetrics is a public function to initialise a new instance of ChainMetrics
func NewChainMetrics() *ChainMetrics {
	metrics := &ChainMetrics{
		AmountTransferred: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "chainbridge_total_amount_transferred",
			Help: "Number of tokens transferred across bridge",
		}),
		NumberOfTransfers: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "chainbridge_total_number_of_transfers",
			Help: "Number of transfers occurred across bridge",
		}),
	}

	prometheus.MustRegister(metrics.AmountTransferred)
	prometheus.MustRegister(metrics.NumberOfTransfers)

	return metrics
}
