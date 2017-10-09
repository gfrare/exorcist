package main

import (
	"fmt"

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

	cmdInvoke := &cobra.Command{}

	cmdEnumerate := &cobra.Command{}

	cmdBanish := &cobra.Command{}

	rootCmd := &cobra.Command{Use: "exorcist"}
	rootCmd.AddCommand(cmdFake)
	rootCmd.AddCommand(cmdInvoke)
	rootCmd.AddCommand(cmdEnumerate)
	rootCmd.AddCommand(cmdBanish)
	rootCmd.Execute()

}
