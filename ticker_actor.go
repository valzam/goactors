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
	ctx        context.Context
	cancelCtx  context.Context
	cancelFunc context.CancelFunc
	stopped    bool
	running    bool

	// Processing
	state       S
	processFunc func(context.Context, *S)
	ticker      *time.Ticker
}

func NewTickerActor[S any](ctx context.Context,
	initialState S,
	processFunc func(context.Context, *S),
	invokeEvery time.Duration) TickerActor[S] {
	cancelCtx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(invokeEvery)

	return &tickerActor[S]{
		ctx:         ctx,
		cancelCtx:   cancelCtx,
		cancelFunc:  cancel,
		state:       initialState,
		processFunc: processFunc,
		ticker:      ticker,
	}
}

func (c *tickerActor[S]) start() {
	if c.wasStopped() {
		panic("cannot start actor that was manually stopped before")
	}

	c.running = true

	for {
		select {
		case <-c.ticker.C:
			c.processFunc(c.ctx, &c.state)
		case <-c.cancelCtx.Done():
			println("shutting down actor")
			c.stopped = true
			c.running = false
			return
		}
	}
}

func (c *tickerActor[S]) stop() {
	c.cancelFunc()
}

func (c *tickerActor[S]) wasStopped() bool {
	return c.stopped
}

func (c *tickerActor[S]) isRunning() bool {
	return c.running
}
