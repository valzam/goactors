package goactor

import (
	"context"
	"time"
)

type TickerActor[S any] interface {
	genericActor
}

type tickerActor[S any] struct {
	// Admin
	ctx          context.Context
	cancelCtx    context.Context
	cancelFunc   context.CancelFunc
	stopped      bool
	runningChan  chan struct{}
	shutdownChan chan struct{}

	// Processing
	state       S
	processFunc func(context.Context, *S)
	ticker      *time.Ticker
}

func NewTickerActor[S any](ctx context.Context,
	initialState S,
	processFunc func(context.Context, *S),
	invokeEvery time.Duration,
) TickerActor[S] {
	cancelCtx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(invokeEvery)

	return &tickerActor[S]{
		ctx:          ctx,
		cancelCtx:    cancelCtx,
		cancelFunc:   cancel,
		state:        initialState,
		processFunc:  processFunc,
		ticker:       ticker,
		runningChan:  make(chan struct{}),
		shutdownChan: make(chan struct{}),
	}
}

func (c *tickerActor[S]) start() {
	close(c.runningChan)

	for {
		select {
		case <-c.ticker.C:
			c.processFunc(c.ctx, &c.state)
		case <-c.cancelCtx.Done():
			println("shutting down actor")
			return
		}
	}
}

func (c *tickerActor[S]) stop() {
	c.cancelFunc()
}

func (c *tickerActor[S]) running() <-chan struct{} {
	return c.runningChan
}

func (c *tickerActor[S]) shutdown() <-chan struct{} {
	return c.shutdownChan
}
