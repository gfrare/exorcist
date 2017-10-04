package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var cmdFake = &cobra.Command{
		Use:   "fake",
		Short: "Print fake message",
		Long:  "A fake for a day",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("This is a fake message")
		},
	}

	var rootCmd = &cobra.Command{Use: "exorcist"}
	rootCmd.AddCommand(cmdFake)
	rootCmd.Execute()

}
