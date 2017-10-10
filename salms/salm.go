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

var currentRituals map[string]rituals.Ritual
var currentGauges map[string]prometheus.Gauge

func init() {
	currentGauges = make(map[string]prometheus.Gauge)
}

// InitAndExecuteMetrics public function
func InitAndExecuteMetrics() {
	ritualsList := rituals.ListRituals()

	removableRituals := markRemovableRituals(ritualsList)
	for _, metric := range removableRituals {
		unregisterMetric(metric)
	}

	newRituals := markNewRituals(ritualsList)
	for metric, ritual := range newRituals {
		gauge := registerMetric(metric)
		go executeMetric(gauge, ritual)
	}

	currentRituals = ritualsList

	// var ritualsList map[string]rituals.Ritual

	// if currentRituals == nil {
	// 	ritualsList = rituals.ListRituals()
	// } else {
	// 	// TODO: compare
	// 	// ritualsList = rituals.ListRituals()
	// }
	// for metric, ritual := range ritualsList {
	// 	gauge := initMetric(metric)
	// 	go executeMetric(gauge, ritual)
	// }
	// currentRituals = ritualsList
}

func markRemovableRituals(ritualList map[string]rituals.Ritual) []string {
	if currentRituals == nil || len(currentRituals) == 0 {
		return make([]string, 0)
	}
	removableRituals := make([]string, 0, 100)

	for k, v := range currentRituals {
		ritual, ok := ritualList[k]
		if ok == true {
			if ritual != v {
				removableRituals = append(removableRituals, k)
			}
		} else {
			removableRituals = append(removableRituals, k)
		}
	}
	fmt.Println("REMOVABLE-RITUALS", removableRituals)
	return removableRituals
}

func markNewRituals(ritualList map[string]rituals.Ritual) map[string]rituals.Ritual {
	if currentRituals == nil || len(currentRituals) == 0 {
		return ritualList
	}

	additionalRituals := make(map[string]rituals.Ritual)
	for k, v := range ritualList {
		ritual, ok := currentRituals[k]
		if ok == true {
			if ritual != v {
				additionalRituals[k] = v
			}
		} else {
			additionalRituals[k] = v
		}
	}
	fmt.Println("NEW-RITUALS", additionalRituals)
	return additionalRituals
}

func registerMetric(metric string) prometheus.Gauge {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metric,
		Help: metric,
	})
	prometheus.MustRegister(gauge)
	currentGauges[metric] = gauge
	return gauge
}

func unregisterMetric(metric string) {
	gauge := currentGauges[metric]
	prometheus.Unregister(gauge)
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
