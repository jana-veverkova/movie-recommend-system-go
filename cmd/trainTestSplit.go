package cmd

import (
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/traintestsplit"
	"github.com/spf13/cobra"
)

var trainTestSplitCmd = &cobra.Command{
	Use:   "train-test-split",
	Short: "Splits data into train and test set.",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := traintestsplit.Split("data/processed/edx.csv", "data/trainTest")
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		return nil
	},
}
