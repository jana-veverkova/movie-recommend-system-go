package cmd

import (
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/models"
	"github.com/spf13/cobra"
)

var modelTrainCmd = &cobra.Command{
	Use:       "model-train",
	Short:     "Trains model given using train data and evaluates based on test data.",
	ValidArgs: []string{"modelv0", "modelv2", "modelv4", "train", "test", "edx", "holdout_test"},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(3)(cmd, args); err != nil {
			return err
		}

		return cobra.OnlyValidArgs(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		trainDataSourceUrl, err := getDataSourcePath(args[1])
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		testDataSourceUrl, err := getDataSourcePath(args[2])
		if err != nil {
			printErrorWithStack(err)
			return err
		}
		
		err = models.Train(args[0], trainDataSourceUrl, testDataSourceUrl)
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		return nil
	},
}
