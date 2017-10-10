package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/cobra"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Grimoire is the rapresentation of json file struct
type Grimoire struct {
	Rituals map[string]Ritual `json:"rituals"`
}

// Ritual jkj
type Ritual struct {
	Command string `json:"command"`
	Timer   uint8  `json:"timer"`
}

func main() {
	var name string

	cmdExorcism := &cobra.Command{
		Use:   "exorcism",
		Short: "Begin an exorcism",
		Long:  "Begin an exorcism",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("exorcism", args)
			initServer(args)
		},
	}

	cmdInvoke := &cobra.Command{
		Use:   "invoke",
		Short: "Invoke a daemon",
		Long:  "Invoke a daemon",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("invocation", args)
			fmt.Println("flag", name)
			command := strings.Join(args, " ")
			invoke(name, command, 2)
		},
	}

	cmdRecite := &cobra.Command{
		Use:   "recite",
		Short: "Recite a salm",
		Long:  "Recite a salm",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("recite", args)
		},
	}

	cmdBanish := &cobra.Command{
		Use:   "banish",
		Short: "Banish a daemon",
		Long:  "Banish a daemon",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("banish", args)
		},
	}

	cmdInvoke.Flags().StringVarP(&name, "name", "n", "default", "give a name to invocation")

	rootCmd := &cobra.Command{Use: "exorcist"}
	rootCmd.AddCommand(cmdExorcism)
	rootCmd.AddCommand(cmdInvoke)
	rootCmd.AddCommand(cmdRecite)
	rootCmd.AddCommand(cmdBanish)
	rootCmd.Execute()

	// Expose the registered metrics via HTTP.
	// http.Handle("/metrics", promhttp.Handler())
	// log.Fatal(http.ListenAndServe(*addr, nil))
}

func initServer(args []string) {
	// print del file al cambiamento
}

func invoke(metric string, command string, timer uint8) {
	ritual := Ritual{command, timer}
	addRitual(metric, ritual)
}

func addRitual(metric string, ritual Ritual) {
	grimoire := readGrimoire()
	grimoire.Rituals[metric] = ritual
	writeGrimoire(grimoire)
}

func removeRitual(metric string) {
	grimoire := readGrimoire()
	delete(grimoire.Rituals, metric)
	writeGrimoire(grimoire)
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

// func runCommand(args []string) {
// 	str := strings.Join(args, " ")

// 	go func() {
// 		for {
// 			out, err := exec.Command("sh", "-c", str).Output()
// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			i, err := strconv.ParseFloat(strings.TrimSpace(string(out[:])), 64)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			fmt.Println(i)

// 			time.Sleep(2 * time.Second)
// 		}
// 	}()
// }
