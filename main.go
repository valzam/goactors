package main

import (
	"context"
	"fmt"
	"goactor/base"
	"time"
)

func main() {
	s, shutdown := base.NewRecoverSupervisor()
	defer func() {
		shutdown()
		time.Sleep(1 * time.Second)
	}()

	var c = NewCounter(context.Background())
	s.Register(c)

	c.Send(CounterMsgIncr{by: 1})
	c.Send(CounterMsgIncr{by: 2})
	c.Send(CounterMsgIncr{by: 1})

	time.Sleep(1 * time.Second)
	println(fmt.Sprintf("counter actor state: %v", c.Get()))
}
