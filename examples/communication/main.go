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

	gen := NewGenerator(ctx)
	s.RegisterActor(gen)

	printer := NewPrinter(ctx)
	s.RegisterActor(printer)

	gen.Send(GeneratorMsg{}, printer.Ref())
	gen.Send(GeneratorMsg{}, printer.Ref())
	gen.Send(GeneratorMsg{}, printer.Ref())

	time.Sleep(1 * time.Second)
}
