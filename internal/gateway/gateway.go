package gateway

import (
	"context"
	"fmt"

	"github.com/otakakot/sample-go-unit-test/internal/model"
	"github.com/otakakot/sample-go-unit-test/internal/repository"
)

var _ repository.Repository = (*Gateway)(nil)

type Gateway struct {
	models map[string]model.Model
}

func New() *Gateway {
	return &Gateway{
		models: map[string]model.Model{},
	}
}

func (gtw *Gateway) Find(
	ctx context.Context,
	id string,
) (model.Model, error) {
	mdl, ok := gtw.models[id]
	if !ok {
		return model.Model{}, fmt.Errorf("not found model, id: %s", id)
	}

	return mdl, nil
}

func (gtw *Gateway) Save(
	ctx context.Context,
	mdl model.Model,
) error {
	gtw.models[mdl.ID] = mdl

	return nil
}
