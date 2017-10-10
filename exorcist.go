package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gfrare/exorcist/rituals"
	"github.com/gfrare/exorcist/salms"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

func main() {
	var name string
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
			salms.InitAndExecuteMetrics()
			initServer(port)
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
	cmdSummon.Flags().StringVarP(&port, "port", "p", "8080", "give a port to exorcist")

	rootCmd := &cobra.Command{Use: "exorcist"}
	rootCmd.AddCommand(cmdExorcism)
	rootCmd.AddCommand(cmdSummon)
	rootCmd.AddCommand(cmdInvoke)
	rootCmd.AddCommand(cmdRecite)
	rootCmd.AddCommand(cmdBanish)
	rootCmd.Execute()
}

func initServer(port string) {
	// Expose the registered metrics via HTTP.
	host := ":" + port
	// http.Handle("/metrics", handlers.LoggingHandler(os.Stdout, promhttp.Handler()))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(host, nil))
}

func invoke(metric string, command string, timer uint8) {
	ritual := rituals.Ritual{Command: command, Timer: timer}
	rituals.AddRitual(metric, ritual)
}
