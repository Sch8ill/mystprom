package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sch8ill/mystprom/config"
)

// custom metrics registry to discard default go metrics
var registry = prometheus.NewRegistry()

func Listen() error {
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(config.MetricsAddress, nil); err != nil {
		return err
	}

	return nil
}
