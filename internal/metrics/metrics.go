package metrics

import (
	"time"

	"github.com/trustwallet/blockatlas/db"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	nameSpace = "blockatlas"

	workerBlockParsing = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: "worker",
			Name:      "block_parsing",
			Help:      "Last parsed block",
		},
		[]string{
			"coin",
			"priority",
		},
	)
)

func setupUpdateTrackerMetrics(db *db.Instance) {
	go func() {
		for {
			trackers, err := db.GetLastParsedBlockNumbers()
			if err != nil {
				continue
			}
			for _, tracker := range trackers {
				labels := prometheus.Labels{"coin": tracker.Coin, "priority": tracker.Priority}
				workerBlockParsing.With(labels).Set(float64(tracker.Height))
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func Setup(db *db.Instance) {
	prometheus.DefaultRegisterer.Unregister(prometheus.NewGoCollector())
	prometheus.DefaultRegisterer.Unregister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	prometheus.MustRegister(workerBlockParsing)

	setupUpdateTrackerMetrics(db)
}
