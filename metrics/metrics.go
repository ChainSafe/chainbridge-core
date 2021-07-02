package metrics

import (
	"fmt"
	"net/http"
)

// ChainMetrics is a public struct that includes data related to transfers occuring over the chainbridge
type ChainMetrics struct {
	TotalAmountTransferred int
	TotalNumberOfTransfers int
}

// New is a public function to initialise a new instance of ChainMetrics
func New() *ChainMetrics {
	return &ChainMetrics{}
}

// MetricsHandler is a public method to provide a formatted summary of chain metrics for the metrics http server
func (c *ChainMetrics) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Metrics:\n------\nTotal Amount Transferred: %v\nTotal Number Of Transfers: %v\n", c.TotalAmountTransferred, c.TotalNumberOfTransfers)
}
