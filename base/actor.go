package base

import "context"

type Actor interface {
	start()
	stop()
	IsStopped() bool
}

type BasicActor[S any, T any] struct {
	// Admin
	ctx      context.Context
	stopFunc context.CancelFunc
	stopped  bool

	// Processing
	inputChan   chan T
	state       S
	processFunc func(ctx context.Context, currentState S, msg T) S
}

func (c *BasicActor[S, T]) start() {
actorLoop:
	for {
		select {
		case i := <-c.inputChan:
			c.state = c.processFunc(c.ctx, c.state, i)
		case <-c.ctx.Done():
			println("shutting down")
			close(c.inputChan)
			break actorLoop
		}
	}
}

func (c *BasicActor[S, T]) stop() {
	c.stopped = true
	c.stopFunc()
}

func (c *BasicActor[S, T]) IsStopped() bool {
	return c.stopped
}
