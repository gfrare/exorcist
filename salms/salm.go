package salms

import (
	"crypto/sha512"
	"encoding/hex"
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
func InitAndExecuteMetrics(page string) {
	ritualsList := rituals.ListRituals(page)

	removableRituals := markRemovableRituals(ritualsList)
	for metric, ritual := range removableRituals {
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
		go executeMetric(gauge, metric, ritual, slaughterChannel)
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

func markRemovableRituals(ritualList map[string]rituals.Ritual) map[string]rituals.Ritual {
	if currentRituals == nil || len(currentRituals) == 0 {
		return make(map[string]rituals.Ritual)
	}

	return markDifferentRituals(currentRituals, ritualList)
}

func markNewRituals(ritualList map[string]rituals.Ritual) map[string]rituals.Ritual {
	if currentRituals == nil || len(currentRituals) == 0 {
		return ritualList
	}

	return markDifferentRituals(ritualList, currentRituals)
}

func markDifferentRituals(firstMap map[string]rituals.Ritual, secondMap map[string]rituals.Ritual) map[string]rituals.Ritual {
	returnMap := make(map[string]rituals.Ritual)

	for k, v := range firstMap {
		ritual, ok := secondMap[k]
		if ok == true {
			if ritual != v {
				returnMap[k] = v
			}
		} else {
			returnMap[k] = v
		}
	}
	return returnMap
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

func executeMetric(gauge prometheus.Gauge, metric string, ritual rituals.Ritual, slaughterChannel chan (bool)) {
	gauge.Set(0)
	run := true
	waitTime := time.Duration(ritual.Timer) * time.Second

	for run {
		output, err := exec.Command("sh", "-c", ritual.Command).Output()
		if err != nil {
			log.Fatal(err)
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(string(output[:])), 64)
		if err != nil {
			log.Fatal(err)
		}
		gauge.Set(value)
		log.Printf("Ritual \"%s\", value: %f", metric, value)

		select {
		case <-slaughterChannel:
			run = false
		case <-time.After(waitTime):
		}
	}
}
