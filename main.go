package main

import (
	"context"
	"log"

	"github.com/otakakot/sample-go-unit-test/internal/gateway"
	"github.com/otakakot/sample-go-unit-test/internal/usecase"
)

func main() {
	name := "otakakot"

	ctx := context.Background()

	gtw := gateway.New()

	uc := usecase.New(gtw)

	mdl, err := uc.Create(ctx, name)
	if err != nil {
		panic(err)
	}

	log.Printf("saved model: %+v", mdl)

	got, err := uc.Read(ctx, mdl.ID)
	if err != nil {
		panic(err)
	}

	log.Printf("got model: %+v", got)
}
