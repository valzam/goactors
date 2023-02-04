package main

import (
	"context"
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

	var c = NewPrinter(ctx)
	s.RegisterActor(c)

	time.Sleep(2 * time.Second)
}
