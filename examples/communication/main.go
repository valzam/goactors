package main

import (
	"context"
	"goactor"
	"time"
)

func main() {
	s, shutdown := goactor.NewRecoverSupervisor()
	defer func() {
		shutdown()
		time.Sleep(1 * time.Second)
	}()

	ctx := context.Background()
	gen := NewGenerator(ctx)
	s.Register(gen)

	printer := NewPrinter(ctx)
	s.Register(printer)

	gen.Send(GeneratorMsg{}, printer.Ref())
	gen.Send(GeneratorMsg{}, printer.Ref())
	gen.Send(GeneratorMsg{}, printer.Ref())

	time.Sleep(1 * time.Second)
}
