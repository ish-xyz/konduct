package exporter

import "fmt"

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
	for _, x := range rep.Results {
		fmt.Println(x)
	}
	return fmt.Errorf("invalid exporter mode selected")
}

// import (
// 	"fmt"

// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/push"
// )

// func ExamplePusher_Push() {
// 	completionTime := prometheus.NewGauge(prometheus.GaugeOpts{
// 		Name: "db_backup_last_completion_timestamp_seconds",
// 		Help: "The timestamp of the last successful completion of a DB backup.",
// 	})
// 	completionTime.SetToCurrentTime()
// 	if err := push.New("http://pushgateway:9091", "db_backup").
// 		Collector(completionTime).
// 		Grouping("db", "customers").
// 		Push(); err != nil {
// 		fmt.Println("Could not push completion time to Pushgateway:", err)
// 	}
// }
//
//
// kubetest_result{name="test-case", failed=int, succeded=int, source=kubeapi/filesystem}
// kubetest_status 0/1
