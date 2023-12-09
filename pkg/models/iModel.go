package models

import "github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"

type IModel interface {
	predict(data *datarepository.DataSet) ([]float32, error)
	train(data *datarepository.DataSet, params trainParams) error
}