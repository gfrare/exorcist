package salms

import (
	"crypto/sha512"
	"encoding/hex"
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
var slaughterChannels map[string]chan bool

func init() {
	currentGauges = make(map[string]prometheus.Gauge)
	slaughterChannels = make(map[string]chan bool)
}

// InitAndExecuteMetrics public function
func InitAndExecuteMetrics() {
	ritualsList := rituals.ListRituals()

	removableRituals := markRemovableRituals(ritualsList)
	for _, metric := range removableRituals {
		ritual := currentRituals[metric]
		sin := generateOriginalSin(metric, ritual)
		slaughterChannel := slaughterChannels[sin]
		slaughterChannel <- true
		close(slaughterChannel)
		delete(slaughterChannels, sin)
		unregisterMetric(metric)
	}

	newRituals := markNewRituals(ritualsList)
	for metric, ritual := range newRituals {
		sin := generateOriginalSin(metric, ritual)
		gauge := registerMetric(metric)
		slaughterChannel := make(chan bool)
		slaughterChannels[sin] = slaughterChannel
		go executeMetric(gauge, ritual, slaughterChannel)
	}

	currentRituals = ritualsList
}

func generateOriginalSin(metric string, ritual rituals.Ritual) string {
	seed := metric + ritual.Command + strconv.Itoa(int(ritual.Timer))
	hasher := sha512.New512_256()
	hasher.Write([]byte(seed))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash
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
	delete(currentGauges, metric)
	prometheus.Unregister(gauge)
}

func executeMetric(gauge prometheus.Gauge, ritual rituals.Ritual, slaughterChannel chan (bool)) {
	gauge.Set(0)
	run := true
	waitTime := time.Duration(ritual.Timer) * time.Second

	for run {
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

		select {
		case <-slaughterChannel:
			run = false
		case <-time.After(waitTime):
		}
	}
}
