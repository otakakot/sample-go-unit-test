package repository

import (
	"context"

	model "github.com/otakakot/sample-go-unit-test/internal/model"
)

//go:generate mockgen -source repository.go -destination repository_mock.go -package repository

type Repository interface {
	Save(context.Context, model.Model) error
	Find(context.Context, string) (model.Model, error)
}
