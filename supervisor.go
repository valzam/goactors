package goactor

import (
	"sync"
)

type Supervisor interface {
	Register(a baseActor)
	shutdown()
}

type RecoverSupervisor struct {
	actors []baseActor
}

func NewRecoverSupervisor() (*RecoverSupervisor, func()) {
	s := &RecoverSupervisor{}
	return s, func() { s.shutdown() }
}

func (s *RecoverSupervisor) Register(a baseActor) {
	s.actors = append(s.actors, a)

	wg := sync.WaitGroup{}
	go func() {
		for !a.wasStopped() {
			wg.Add(1)
			func() {
				defer func() {
					if r := recover(); r != nil {
						println("restarting actor")
					}
				}()
				wg.Done()
				a.start()
			}()
		}
	}()

	// Wait for the goroutine to call `start` once
	wg.Wait()
}

func (s *RecoverSupervisor) shutdown() {
	for _, a := range s.actors {
		a.stop()
	}
}
