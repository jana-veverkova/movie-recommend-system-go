package models

import (
	"fmt"

	"github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"
	"github.com/jana-veverkova/movie-recommend-system-go/pkg/modelevaluation"
	"github.com/pkg/errors"
)

func Train(modelType string, trainDataSource string, testDataSource string) error {
	trainData := datarepository.GetData(trainDataSource)
	testData := datarepository.GetData(testDataSource)

	modelBuilder := ModelBuilder{}
	model := modelBuilder.buildModel(modelType)

	for k := 0; k < 21; k = k + 2 {
		err := model.train(trainData, trainParams{lambda: float32(k)})
		if err != nil {
			return errors.WithStack(err)
		}

		predictions, err := model.predict(testData)
		if err != nil {
			return errors.WithStack(err)
		}

		summary, err := modelevaluation.Summarize(testData.Ratings, predictions)
		if err != nil {
			return errors.WithStack(err)
		}

		fmt.Printf("Parameters: Lambda: %d => ", k)
		fmt.Printf("Summary: RMSE: %.4f, MAE: %.4f \n", summary.Rmse, summary.Mae)
	}

	return nil
}
