package main

import (
	"context"
	"fmt"
	"goactor"
	"time"
)

func main() {
	ctx := context.Background()

	s, shutdown := goactor.NewRecoverSupervisor(ctx)
	defer func() {
		shutdown()
		time.Sleep(1 * time.Second)
	}()

	var c = NewCounter(ctx)
	s.RegisterActor(c)

	// Passing a nil channel means the actor won't try to send a response
	c.Send(CounterMsgIncr{by: 1}, nil)
	c.Send(CounterMsgIncr{by: 2}, nil)
	c.Send(CounterMsgIncr{by: 1}, nil)

	time.Sleep(1 * time.Second)

	resChan := make(chan CounterResp)
	c.Send(CounterMsgIncr{}, resChan)
	res := <-resChan
	close(resChan)

	fmt.Printf("counter actor state: %v\n", res)
}
