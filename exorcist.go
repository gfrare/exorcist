package main

import (
	"fmt"
	"strings"

	"github.com/gfrare/exorcist/rituals"
	"github.com/spf13/cobra"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

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
	ritual := rituals.Ritual{Command: command, Timer: timer}
	rituals.AddRitual(metric, ritual)
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
