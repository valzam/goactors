package main

import (
	"context"
	"fmt"
	"goactor"
	"time"
)

func main() {
	s, shutdown := goactor.NewRecoverSupervisor()
	defer func() {
		shutdown()
		time.Sleep(1 * time.Second)
	}()

	var c = NewCounter(context.Background())
	s.Register(c)

	c.Send(CounterMsgIncr{by: 1}, nil)
	c.Send(CounterMsgIncr{by: 2}, nil)
	c.Send(CounterMsgIncr{by: 1}, nil)

	time.Sleep(1 * time.Second)
	println(fmt.Sprintf("counter actor state: %v", c.GetState()))
}
