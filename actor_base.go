package goactor

import "context"

type baseActor interface {
	start()
	prepareStart()
	stop()
	IsStopped() bool
	IsRunning() bool
}

type baseActorImpl[S any, T any] struct {
	// Admin
	ctx      context.Context
	stopFunc context.CancelFunc
	stopped  bool
	running  bool

	// Processing
	inputChan   chan T
	state       S
	processFunc func(state *S, msg T)
}

func (c *baseActorImpl[S, T]) prepareStart() {
	c.running = true
}

func (c *baseActorImpl[S, T]) start() {
	if c.IsStopped() {
		panic("cannot restart actor")
	}

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

func (c *baseActorImpl[S, T]) stop() {
	c.stopped = true
	c.running = false
	c.stopFunc()
}

func (c *baseActorImpl[S, T]) IsStopped() bool {
	return c.stopped
}

func (c *baseActorImpl[S, T]) IsRunning() bool {
	return c.running
}
