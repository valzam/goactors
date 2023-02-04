package main

import (
	"context"
	"fmt"
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
	return goactor.NewActor[CounterActorState, CounterMsgIncr, CounterResp](
		ctx,
		CounterActorState{counter: 0},
		incrementCounter,
	)
}

func incrementCounter(_ context.Context, state *CounterActorState, msg CounterMsgIncr) CounterResp {
	fmt.Printf("incrementing counter by %d\n", msg.by)
	state.counter += msg.by

	return getCounterState(state)
}

func getCounterState(ca *CounterActorState) CounterResp {
	return CounterResp{CurrentValue: ca.counter}
}
