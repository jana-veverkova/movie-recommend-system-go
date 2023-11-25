package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "movie-recommend-system",
	Short: "Description",
}

func Execute() {
	rootCmd.AddCommand(trainTestSplitCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func printErrorWithStack(err error) {
	if err == nil {
		return
	}

	fmt.Printf("%+v\n", err)
}