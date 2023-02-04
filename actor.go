package goactor

import "context"

type Actor interface {
	start()
	stop()
	IsStopped() bool
}

type actor[S any, T any] struct {
	// Admin
	ctx      context.Context
	stopFunc context.CancelFunc
	stopped  bool

	// Processing
	inputChan   chan T
	state       S
	processFunc func(state *S, msg T)
}

func (c *actor[S, T]) start() {
actorLoop:
	for {
		select {
		case i := <-c.inputChan:
			c.processFunc(&c.state, i)
		case <-c.ctx.Done():
			println("shutting down")
			close(c.inputChan)
			break actorLoop
		}
	}
}

func (c *actor[S, T]) stop() {
	c.stopped = true
	c.stopFunc()
}

func (c *actor[S, T]) IsStopped() bool {
	return c.stopped
}
