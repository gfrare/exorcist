package rituals

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// ConfigurationFile is the name of the file that contains the program configuration
const ConfigurationFile = "grimoire.json"

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
	if _, exists := grimoire.Rituals[metric]; exists {
		delete(grimoire.Rituals, metric)
		writeGrimoire(grimoire)
	}
}

// ListRituals public functions
func ListRituals() map[string]Ritual {
	return readGrimoire().Rituals
}

func readGrimoire() Grimoire {
	// Check if the configuration file exists. If it doesn't create it
	if _, err := os.Stat(ConfigurationFile); os.IsNotExist(err) {
		file, err := os.Create(ConfigurationFile)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	b, err := ioutil.ReadFile(ConfigurationFile)
	if err != nil {
		log.Fatal(err)
	}
	var grimoire Grimoire
	if len(b) == 0 {
		grimoire = Grimoire{}
	} else if err := json.Unmarshal(b, &grimoire); err != nil {
		log.Fatal(err)
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
	if err := ioutil.WriteFile(ConfigurationFile, b, 0666); err != nil {
		log.Fatal(err)
	}
}
