package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/sirupsen/logrus"
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

	logrus.Infoln("exporting results to pushgateway")

	pusher := push.New(e.Address, "kubetest")

	for _, testresult := range rep.Results {

		metric := prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "kubetest_result",
				Help: "The result of a single e2e test case",
				ConstLabels: map[string]string{
					"name":    testresult.Name,
					"message": testresult.Message,
				},
			},
		)
		if testresult.Status {
			metric.Set(1)
		}
		pusher = pusher.Collector(metric)
	}

	err := pusher.Push()

	return err
}

func (e *PushgatewayExporter) IsVerbose() bool {
	return false
}
