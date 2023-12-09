package models

import "strings"

type ModelBuilder struct {
}

func (b *ModelBuilder) buildModel(modelType string) IModel {
	model := Model{}

	if strings.ToLower(modelType) == "modelv0" {

		return &Intercept{
			model: &model,
		}

	} else if strings.ToLower(modelType) == "modelv2" {

		return &UserEffect{
			model: &MovieEffect{
				model: &Intercept{
					model: &model,
				},
			},
		}

	} else if strings.ToLower(modelType) == "modelv4" {

		return &UserEffect{
			model: &MovieEffect{
				model: &ActorEffect{
					model: &DirectorEffect{
						model: &Intercept{
							model: &model,
						},
					},
				},
			},
		}
	}

	return nil
}
