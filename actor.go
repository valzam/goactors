package goactor

import (
	"context"
)

type Actor[S any, I any, R any] interface {
	genericActor
	Send(msg I, resp chan R) bool
	Ref() chan I
}

type msg[T any, R any] struct {
	input      T
	returnChan chan R
}

type asyncActor[S any, I any, R any] struct {
	// Admin
	ctx          context.Context
	cancelCtx    context.Context
	cancelFunc   context.CancelFunc
	runningChan  chan struct{}
	shutdownChan chan struct{}

	// Processing
	state S

	inputChan   chan msg[I, R]
	inputChanFF chan I

	processFunc   func(state *S, msg msg[I, R])
	processFuncFF func(state *S, msg I)
}

func NewActor[S any, I any, R any](ctx context.Context,
	initialState S,
	processFunc func(context.Context, *S, I) R,
) Actor[S, I, R] {
	cancelCtx, cancel := context.WithCancel(context.Background())

	// Default processing function that sends a response to the response channel passed in through Send
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

	// Fire and forget channel to accept messages from other actors without sending a response back
	processFuncFFImpl := func(state *S, msg I) {
		processFunc(ctx, state, msg)
	}

	return &asyncActor[S, I, R]{
		ctx:           ctx,
		cancelCtx:     cancelCtx,
		cancelFunc:    cancel,
		inputChan:     make(chan msg[I, R]),
		inputChanFF:   make(chan I),
		state:         initialState,
		processFunc:   processFuncImpl,
		processFuncFF: processFuncFFImpl,
		runningChan:   make(chan struct{}),
		shutdownChan:  make(chan struct{}),
	}
}

func (c *asyncActor[S, I, R]) Send(input I, resp chan R) bool {
	// TODO Should this error if the actor hasn't been started yet?
	// Problem: Would need a lock on the `running` variable for every send

	c.inputChan <- msg[I, R]{input, resp}
	return true
}

func (c *asyncActor[S, I, R]) Ref() chan I {
	return c.inputChanFF
}

func (c *asyncActor[S, I, R]) start() {
	close(c.runningChan)

	for {
		select {
		case i := <-c.inputChanFF:
			c.processFuncFF(&c.state, i)
		case i := <-c.inputChan:
			c.processFunc(&c.state, i)
		case <-c.cancelCtx.Done():
			println("shutting down actor")
			close(c.inputChan)
			close(c.shutdownChan)
			return
		}
	}
}

func (c *asyncActor[S, I, R]) stop() {
	c.cancelFunc()
}

func (c *asyncActor[S, I, R]) running() <-chan struct{} {
	return c.runningChan
}

func (c *asyncActor[S, I, R]) shutdown() <-chan struct{} {
	return c.shutdownChan
}
