package rituals

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// Grimoire struct
type Grimoire struct {
	Rituals map[string]Ritual `json:"rituals"`
}

// Ritual struct
type Ritual struct {
	Command string `json:"command"`
	Timer   uint8  `json:"timer"`
}

// AddRitual public function
func AddRitual(metric string, ritual Ritual) {
	grimoire := readGrimoire()
	grimoire.Rituals[metric] = ritual
	writeGrimoire(grimoire)
}

// RemoveRitual public function
func RemoveRitual(metric string) {
	grimoire := readGrimoire()
	delete(grimoire.Rituals, metric)
	writeGrimoire(grimoire)
}

// ListRituals public functions
func ListRituals() map[string]Ritual {
	return readGrimoire().Rituals
}

func readGrimoire() Grimoire {
	b, err := ioutil.ReadFile("grimoire.json")
	if err != nil {
		log.Fatal("FATAL1\n", err)
	}
	var grimoire Grimoire
	if len(b) == 0 {
		grimoire = Grimoire{}
	} else if err := json.Unmarshal(b, &grimoire); err != nil {
		log.Fatal("FATAL2\n", err)
	}

	if grimoire.Rituals == nil {
		grimoire.Rituals = make(map[string]Ritual)
	}

	return grimoire
}

func writeGrimoire(grimoire Grimoire) {
	b, err := json.MarshalIndent(grimoire, "", "  ")
	fmt.Println(string(b[:]))
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("grimoire.json", b, 0777); err != nil {
		log.Fatal(err)
	}
}
