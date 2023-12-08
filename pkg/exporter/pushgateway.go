package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	MODE_PUSHGATEWAY = "pushgateway"
	MODE_PROMETHEUS  = "prometheus"
)

func NewPushgatewayExporter(addr string) (Exporter, error) {
	if addr == "" {
		return nil, fmt.Errorf("empty mode ")
	}
	return &PushgatewayExporter{
		//Debug:   debug,
		Address: addr,
	}, nil
}

func (e *PushgatewayExporter) Export(rep *Report) error {

	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kubetest_result",
			Help: "The result of a single e2e test case",
		},
		[]string{"name"},
	)

	for _, testresult := range rep.Results {

		value := float64(0)
		if testresult.Status {
			value = 1
		}
		metric.WithLabelValues(testresult.Name).Set(value)
	}

	err := push.New(e.Address, "kubetest").Gatherer(prometheus.DefaultGatherer).Push()

	return err
}
