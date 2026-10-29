package telemetry

import "github.com/VictoriaMetrics/metrics"

var RequestCounter = metrics.NewCounter(`http_request_total{handler="root"}`)
