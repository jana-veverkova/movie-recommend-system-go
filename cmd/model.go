package cmd

import (
	"fmt"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/models"
	"github.com/spf13/cobra"
)

var modelv0TrainCmd = &cobra.Command{
	Use:   "model-train",
	Short: "Trains model based on given arguments.",
	ValidArgs: []string{"modelv0", "modelv2", "train", "test", "edx", "holdout_test"},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return err
		}

		return cobra.OnlyValidArgs(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := getModelByName(args[0])
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		dataSourceUrl, err := getDataSourcePath(args[1])
		if err != nil {
			printErrorWithStack(err)
			return err
		}
		
		err = models.Train(m, dataSourceUrl)
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		return nil
	},
}

var modelv0EvaluateCmd = &cobra.Command{
	Use:   "model-evaluate",
	Short: "Evaluates model base on arguments.",
	Long: "Evaluates model. Arg 1 => model name, arg 2 => dataset used for training, arg 3 => dataset used for prediction.",
	ValidArgs: []string{"modelv0", "modelv2", "train", "test", "edx", "holdout_test"},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(3)(cmd, args); err != nil {
			return err
		}

		return cobra.OnlyValidArgs(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := getModelByName(args[0])
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		trainedOnUrl, err := getDataSourcePath(args[1])
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		dataSourceUrl, err := getDataSourcePath(args[2])
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		summary, err := models.Evaluate(m, trainedOnUrl, dataSourceUrl)
		if err != nil {
			printErrorWithStack(err)
			return err
		}

		fmt.Println("Summary:")
		fmt.Printf("   rmse: %.2f \n", summary.Rmse)

		return nil
	},
}
