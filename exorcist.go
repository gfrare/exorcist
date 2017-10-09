package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	cmdFake := &cobra.Command{
		Use:   "fake",
		Short: "Print fake message",
		Long:  "A fake for a day",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("This is a fake message")
		},
	}

	cmdInvoke := &cobra.Command{
		Use:   "invoke",
		Short: "Invoke a daemon",
		Long:  "Invoke a daemon",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("invocation", args)
			runCommand(args)
		},
	}

	cmdEnumerate := &cobra.Command{}

	cmdBanish := &cobra.Command{}

	rootCmd := &cobra.Command{Use: "exorcist"}
	rootCmd.AddCommand(cmdFake)
	rootCmd.AddCommand(cmdInvoke)
	rootCmd.AddCommand(cmdEnumerate)
	rootCmd.AddCommand(cmdBanish)
	rootCmd.Execute()

}

func runCommand(args []string) {
	str := strings.Join(args, " ")
	out, err := exec.Command("sh", "-c", str).Output()
	if err != nil {
		log.Fatal(err)
	}

	i, err := strconv.ParseFloat(strings.TrimSpace(string(out[:])), 64)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(i)
}
