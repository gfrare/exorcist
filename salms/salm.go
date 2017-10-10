package salms

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gfrare/exorcist/rituals"
	"github.com/prometheus/client_golang/prometheus"
)

// InitAndExecuteMetrics public function
func InitAndExecuteMetrics() {
	rituals := rituals.ListRituals()
	for metric, ritual := range rituals {
		gauge := initMetric(metric)
		go executeMetric(gauge, ritual)
	}
}

func initMetric(metric string) prometheus.Gauge {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metric,
		Help: metric,
	})
	prometheus.MustRegister(gauge)
	return gauge
}

func executeMetric(gauge prometheus.Gauge, ritual rituals.Ritual) {
	gauge.Set(0)

	for {
		output, err := exec.Command("sh", "-c", ritual.Command).Output()
		if err != nil {
			log.Fatal(err)
		}

		i, err := strconv.ParseFloat(strings.TrimSpace(string(output[:])), 64)
		if err != nil {
			log.Fatal(err)
		}
		gauge.Set(i)
		fmt.Println(i)

		time.Sleep(time.Duration(ritual.Timer) * time.Second)
	}
}
