package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/otakakot/sample-go-unit-test/internal/model"
	"github.com/otakakot/sample-go-unit-test/internal/repository"
)

type Usecase struct {
	repository repository.Repository
}

func New(
	repository repository.Repository,
) *Usecase {
	return &Usecase{
		repository: repository,
	}
}

func (uc *Usecase) Create(
	ctx context.Context,
	name string,
) (model.Model, error) {
	mdl := model.Model{
		ID:   uuid.NewString(),
		Name: name,
	}

	if err := uc.repository.Save(ctx, mdl); err != nil {
		return model.Model{}, fmt.Errorf("failed to save model: %w", err)
	}

	return mdl, nil
}

func (uc *Usecase) Read(
	ctx context.Context,
	id string,
) (model.Model, error) {
	mdl, err := uc.repository.Find(ctx, id)
	if err != nil {
		return model.Model{}, fmt.Errorf("failed to find model: %w", err)
	}

	return mdl, nil
}
