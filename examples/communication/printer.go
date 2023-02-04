package main

import (
	"context"
	"fmt"
	"goactor"
)

type PrinterActor = goactor.Actor[PrinterActorState, GeneratorResp, PrinterActorResp]
type PrinterActorState struct{}

type PrinterActorResp struct{}

func NewPrinter(ctx context.Context) PrinterActor {
	return goactor.NewActor[PrinterActorState, GeneratorResp, PrinterActorResp](
		ctx,
		PrinterActorState{},
		printValue,
	)
}

func printValue(_ context.Context, _ *PrinterActorState, msg GeneratorResp) PrinterActorResp {
	fmt.Printf("got value %d\n", msg.Value)

	return PrinterActorResp{}
}
