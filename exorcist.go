package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gfrare/exorcist/god_eye"
	"github.com/gfrare/exorcist/rituals"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

func main() {
	var name string
	var command string
	var timer *uint8
	var port string

	cmdExorcism := &cobra.Command{
		Use:   "exorcism",
		Short: "Begin an exorcism",
		Long:  "Begin an exorcism",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("exorcism", args)
		},
	}

	cmdSummon := &cobra.Command{
		Use:   "summon",
		Short: "Summon the exorcist",
		Long:  "Summon the exorcist",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("summon", args)
			fmt.Println("flag", port)

			go godEye.Watch()

			initServer(port)
		},
	}

	cmdInvoke := &cobra.Command{
		Use:   "invoke",
		Short: "Invoke a daemon",
		Long:  "Invoke a daemon",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("invocation:", args)
			fmt.Println("name:", name)
			fmt.Println("command:", command)
			fmt.Println("timer:", *timer)
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

	cmdInvoke.Flags().StringVarP(&name, "name", "n", "", "give a name to invocation")
	cmdInvoke.MarkPersistentFlagRequired("name")
	cmdInvoke.Flags().StringVarP(&command, "command", "c", "", "command to execute")
	cmdInvoke.MarkPersistentFlagRequired("command")
	timer = cmdInvoke.Flags().Uint8P("timer", "t", 5, "sleep between command execution")

	cmdSummon.Flags().StringVarP(&port, "port", "p", "8080", "give a port to exorcist")

	rootCmd := &cobra.Command{Use: "exorcist"}
	rootCmd.AddCommand(cmdExorcism)
	rootCmd.AddCommand(cmdSummon)
	rootCmd.AddCommand(cmdInvoke)
	rootCmd.AddCommand(cmdRecite)
	rootCmd.AddCommand(cmdBanish)
	rootCmd.Execute()
}

// Expose the registered metrics via HTTP
func initServer(port string) {
	host := ":" + port
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(host, nil))
}

func invoke(metric string, command string, timer uint8) {
	ritual := rituals.Ritual{Command: command, Timer: timer}
	rituals.AddRitual(metric, ritual)
}
