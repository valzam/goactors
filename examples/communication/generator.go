package main

import (
	"context"
	"goactor"
	"math/rand"
)

type GeneratorActor = goactor.Actor[GeneratorActorState, GeneratorMsg, GeneratorResp]
type GeneratorActorState struct{}

type GeneratorMsg struct{}

type GeneratorResp struct {
	Value int
}

func NewGenerator(ctx context.Context) GeneratorActor {
	return goactor.NewActor[GeneratorActorState, GeneratorMsg, GeneratorResp](
		ctx,
		GeneratorActorState{},
		generateValue,
	)
}

func generateValue(_ context.Context, _ *GeneratorActorState, _ GeneratorMsg) GeneratorResp {
	v := rand.Int()
	return GeneratorResp{v}
}
