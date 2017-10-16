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
	Pages map[string]Page `json:"pages"`
}

// Page struct
type Page struct {
	Rituals map[string]Ritual `json:"rituals"`
}

// Ritual struct
type Ritual struct {
	Command string `json:"command"`
	Timer   uint8  `json:"timer"`
}

// AddRitual public function
func AddRitual(page string, metric string, ritual Ritual) {
	grimoire := readGrimoire()
	fmt.Println(grimoire)
	if _, ok := grimoire.Pages[page]; !ok {
		grimoire.Pages[page] = Page{}
	}
	if grimoire.Pages[page].Rituals == nil {
		currentPage := grimoire.Pages[page]
		currentPage.Rituals = make(map[string]Ritual)
		grimoire.Pages[page] = currentPage
	}
	grimoire.Pages[page].Rituals[metric] = ritual
	writeGrimoire(grimoire)
}

// RemoveRitual public function
func RemoveRitual(page string, metric string) {
	grimoire := readGrimoire()

	if pg, existsPage := grimoire.Pages[page]; existsPage {
		if _, existsRitual := pg.Rituals[metric]; existsRitual {
			delete(pg.Rituals, metric)
			writeGrimoire(grimoire)
		}
	}
}

// ListRituals public functions
func ListRituals(page string) map[string]Ritual {
	grimoire := readGrimoire()
	return grimoire.Pages[page].Rituals
}

// GetGrimoire public functions
func GetGrimoire() Grimoire {
	return readGrimoire()
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

	if grimoire.Pages == nil {
		grimoire.Pages = make(map[string]Page)
	}

	return grimoire
}

func writeGrimoire(grimoire Grimoire) {
	b, err := json.MarshalIndent(grimoire, "", "  ")
	log.Printf("Grimoire: %s", string(b[:]))
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(ConfigurationFile, b, 0666); err != nil {
		log.Fatal(err)
	}
}
