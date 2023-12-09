package models

import "github.com/jana-veverkova/movie-recommend-system-go/pkg/datarepository"

type trainParams struct {
	lambda float32
}

type Model struct {
}

func (m *Model) predict(data *datarepository.DataSet) ([]float32, error) {
	return make([]float32, len(data.Ratings)), nil
}

func (m *Model) train(data *datarepository.DataSet, params trainParams) error {
	return nil
}