package goactor

import "context"

type AsyncActor[S any, I any, R any] interface {
	baseActor
	Send(msg I, resp chan R)
}

type msg[T any, R any] struct {
	input      T
	returnChan chan R
}

type asyncActor[S any, I any, R any] struct {
	actor[S, msg[I, R]]
	processFunc func(context.Context, *S, I) R
	returnChan  chan R
}

func NewAsyncActor[S any, I any, R any](ctx context.Context,
	cancelFunc func(),
	initialState S,
	processFunc func(context.Context, *S, I) R) AsyncActor[S, I, R] {

	processFuncImpl := func(state *S, msg msg[I, R]) {
		res := processFunc(ctx, state, msg.input)

		// Spin up new goroutine to avoid blocking the actor if the return channel is full
		go func() {
			// When returnChan == nil the caller doesn't care about the response
			// We shouldn't send to a nil channel because it will block forever
			if msg.returnChan != nil {
				msg.returnChan <- res
			}
		}()
	}

	return &asyncActor[S, I, R]{
		actor: actor[S, msg[I, R]]{
			ctx:         ctx,
			stopFunc:    cancelFunc,
			inputChan:   make(chan msg[I, R]),
			state:       initialState,
			processFunc: processFuncImpl,
		},
		processFunc: processFunc,
	}
}

func (c *asyncActor[S, I, R]) Send(input I, resp chan R) {
	if !c.IsRunning() {
		panic("cannot send to actor that isn't running")
	}

	go func() {
		c.inputChan <- msg[I, R]{input, resp}
	}()
}
