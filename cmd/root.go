package cmd

import (
	"fmt"
	"os"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/models"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/models/modelv0"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/models/modelv2"
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
	rootCmd.AddCommand(modelEvaluateCmd)
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

func getModelByName(modelName string) (models.Model, error) {
	var m1 modelv0.Modelv0
	var m2 modelv2.Modelv2

	ms := []models.Model{
		&m1, &m2,
	}

	var m models.Model
	for _, m = range(ms) {
		if m.GetName() == modelName {
			return m, nil
		}
	}
	
	if m == nil {
		return nil, errors.New("Model not found.")
	}

	return nil, nil
}