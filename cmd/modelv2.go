package cmd

import (
	"fmt"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/models"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/models/modelv2"
	"github.com/spf13/cobra"
)

var modelv2TrainCmd = &cobra.Command{
	Use:   "modelv2-train",
	Short: "Trains modelv2 based on data source given in argument.",
	ValidArgs: []string{"train", "test", "edx", "holdout_test"},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		return cobra.OnlyValidArgs(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		dataSourceUrl, err := getDataSourcePath(args[0])
		if err != nil {
			printErrorWithStack(err)
			return err
		}
		
		var m modelv2.Modelv2
		err = models.Train(&m, dataSourceUrl)
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		return nil
	},
}

var modelv2EvaluateCmd = &cobra.Command{
	Use:   "modelv2-evaluate",
	Short: "Evaluates modelv2.",
	Long: "Evaluates modelv2. Argument 1 => dataset used for training, arg 2 => dataset used for prediction.",
	ValidArgs: []string{"train", "test", "edx", "holdout_test"},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return err
		}

		return cobra.OnlyValidArgs(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		trainedOnUrl, err := getDataSourcePath(args[0])
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		dataSourceUrl, err := getDataSourcePath(args[1])
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		var m modelv2.Modelv2
		summary, err := models.Evaluate(&m, trainedOnUrl, dataSourceUrl)
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		fmt.Println("Summary:")
		fmt.Printf("   rmse: %.2f \n", summary.Rmse)

		return nil
	},
}
