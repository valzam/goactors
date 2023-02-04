package main

import (
	"context"
	"fmt"
	"goactor"
	"time"
)

type PrinterActor = goactor.TickerActor[PrinterActorState]
type PrinterActorState struct{}

func NewPrinter(ctx context.Context) PrinterActor {
	return goactor.NewTickerActor[PrinterActorState](
		ctx,
		PrinterActorState{},
		printValue,
		100*time.Millisecond,
	)
}

func printValue(_ context.Context, _ *PrinterActorState) {
	fmt.Println("background printer invoked")
}
