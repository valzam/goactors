package base

import "context"

type AsyncActor[S any, T any, R any] interface {
	Actor
	Send(msg T)
	Get() R
}

type BasicAsyncActor[S any, T any, R any] struct {
	BasicActor[S, T]
	responseFunc func(state S) R
}

func NewBasicAsyncActor[S any, T any, R any](ctx context.Context,
	cancelFunc func(),
	initialState S,
	processFunc func(context.Context, S, T) S,
	getResponse func(S) R) *BasicAsyncActor[S, T, R] {

	return &BasicAsyncActor[S, T, R]{
		BasicActor: BasicActor[S, T]{
			ctx:         ctx,
			stopFunc:    cancelFunc,
			inputChan:   make(chan T),
			state:       initialState,
			processFunc: processFunc,
		},
		responseFunc: getResponse,
	}
}

func (c *BasicAsyncActor[S, T, R]) Send(msg T) {
	go func() {
		c.inputChan <- msg
	}()
}

func (c *BasicAsyncActor[S, T, R]) Get() R {
	return c.responseFunc(c.state)
}
