package main

import (
	"context"
	"goactor"
)

type CounterActor = goactor.Actor[CounterActorState, CounterMsgIncr, CounterResp]
type CounterActorState struct {
	counter int
}

type CounterMsgIncr struct {
	by int
}

type CounterResp struct {
	CurrentValue int
}

func NewCounter(ctx context.Context) CounterActor {
	ctx, cancel := context.WithCancel(ctx)
	return goactor.NewActor[CounterActorState, CounterMsgIncr, CounterResp](
		ctx,
		cancel,
		CounterActorState{counter: 0},
		incrementCounter,
	)
}

func incrementCounter(_ context.Context, state *CounterActorState, msg CounterMsgIncr) CounterResp {
	state.counter += msg.by

	return getCounterState(state)
}

func getCounterState(ca *CounterActorState) CounterResp {
	return CounterResp{CurrentValue: ca.counter}
}
