package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gfrare/exorcist/god_eye"
	"github.com/gfrare/exorcist/rituals"
	"github.com/olekukonko/tablewriter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

func main() {
	var name string
	var command string
	var timer *uint8
	var port string

	cmdSummon := &cobra.Command{
		Use:   "summon",
		Short: "Summon the exorcist",
		Long:  "Summon the exorcist",
		Run: func(cmd *cobra.Command, args []string) {
			go godEye.Watch()

			initServer(port)
		},
	}

	cmdInvoke := &cobra.Command{
		Use:   "invoke",
		Short: "Invoke a daemon",
		Long:  "Invoke a daemon",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Invoking ritual \"%s\" with command \"%s\" and timer %d", name, command, *timer)
			invoke(name, command, *timer)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if name == "" {
				log.Fatal("Metric name is mandatory")
			}
			if command == "" {
				log.Fatal("Metric command is mandatory")
			}
		},
	}

	cmdRecite := &cobra.Command{
		Use:   "recite",
		Short: "Recite a salm",
		Long:  "Recite a salm",
		Run: func(cmd *cobra.Command, args []string) {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Metric", "Command", "Timer"})
			table.SetBorder(false)
			table.SetHeaderColor(tablewriter.Color(tablewriter.FgRedColor),
				tablewriter.Color(tablewriter.FgGreenColor),
				tablewriter.Color(tablewriter.FgBlueColor))
			table.SetColumnColor(tablewriter.Color(tablewriter.FgRedColor),
				tablewriter.Color(tablewriter.FgGreenColor),
				tablewriter.Color(tablewriter.FgBlueColor))
			table.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_DEFAULT,
				tablewriter.ALIGN_RIGHT})

			for metric, ritual := range rituals.ListRituals() {
				row := []string{metric, ritual.Command, strconv.Itoa(int(ritual.Timer))}
				table.Append(row)
			}
			table.Render()
		},
	}

	cmdBanish := &cobra.Command{
		Use:   "banish",
		Short: "Banish a daemon",
		Long:  "Banish a daemon",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Banishing ritual \"%s\"", name)
			rituals.RemoveRitual(name)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if name == "" {
				log.Fatal("Metric name is mandatory")
			}
		},
	}

	cmdInvoke.Flags().StringVarP(&name, "name", "n", "", "give a name to invocation")
	cmdInvoke.MarkPersistentFlagRequired("name")
	cmdInvoke.Flags().StringVarP(&command, "command", "c", "", "command to execute")
	cmdInvoke.MarkPersistentFlagRequired("command")
	timer = cmdInvoke.Flags().Uint8P("timer", "t", 5, "sleep between command execution")

	cmdSummon.Flags().StringVarP(&port, "port", "p", "8080", "give a port to exorcist")

	cmdBanish.Flags().StringVarP(&name, "name", "n", "", "give a name to invocation")

	rootCmd := &cobra.Command{Use: "exorcist"}
	rootCmd.AddCommand(cmdSummon)
	rootCmd.AddCommand(cmdInvoke)
	rootCmd.AddCommand(cmdRecite)
	rootCmd.AddCommand(cmdBanish)
	rootCmd.Execute()
}

// Expose the registered metrics via HTTP
func initServer(port string) {
	endpoint := "/metrics"
	log.Printf("Starting exorcist on port %s, endpoint: %s", port, endpoint)
	host := ":" + port
	http.Handle(endpoint, promhttp.Handler())
	log.Fatal(http.ListenAndServe(host, nil))
}

func invoke(metric string, command string, timer uint8) {
	ritual := rituals.Ritual{Command: command, Timer: timer}
	rituals.AddRitual(metric, ritual)
}
