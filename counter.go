package main

import (
	"context"
	"goactor/base"
)

type CounterActor = base.AsyncActor[CounterActorState, CounterMsgIncr, CounterResp]
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
	return base.NewBasicAsyncActor[CounterActorState, CounterMsgIncr, CounterResp](
		ctx,
		cancel,
		CounterActorState{counter: 0},
		incrementCounter,
		getCounterState,
	)
}

func incrementCounter(_ context.Context, state CounterActorState, msg CounterMsgIncr) CounterActorState {
	state.counter += msg.by
	return state
}

func getCounterState(ca CounterActorState) CounterResp {
	return CounterResp{CurrentValue: ca.counter}
}
