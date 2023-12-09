package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "movie-recommend-system",
	Short: "Description",
}

func Execute() {
	rootCmd.AddCommand(trainTestSplitCmd)
	rootCmd.AddCommand(modelTrainCmd)
	rootCmd.AddCommand(downloadBudgetCmd)
	rootCmd.AddCommand(processBudgetCmd)
	rootCmd.AddCommand(downloadCastCmd)
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

func getDataSourcePath(dataSourceArg string) (string, error) {
	dataSourceUrl := ""
	switch dataSourceArg {
	case "train":
		dataSourceUrl = "data/trainTest/train.csv"
	case "test":
		dataSourceUrl = "data/trainTest/test.csv"
	case "edx":
		dataSourceUrl = "data/processed/edx.csv"
	case "holdout_test":
		dataSourceUrl = "data/processed/final_holdout_test.csv"
	}

	if dataSourceUrl == "" {
		return dataSourceUrl, errors.New("DataSourceUrl is not found.")
	}

	return dataSourceUrl, nil
}