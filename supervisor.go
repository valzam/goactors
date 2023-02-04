package goactor

import (
	"sync"

	"github.com/google/uuid"
)

const defaultSupervisorCapacity = 256

type ActorID = uuid.UUID

type Supervisor interface {
	Register(a baseActor) ActorID
	shutdown()
}

type RecoverSupervisor struct {
	mu     sync.Mutex
	actors map[ActorID]baseActor
}

func NewRecoverSupervisor() (*RecoverSupervisor, func()) {
	s := &RecoverSupervisor{
		actors: make(map[ActorID]baseActor, defaultSupervisorCapacity),
	}
	return s, func() { s.shutdown() }
}

func (s *RecoverSupervisor) Register(a baseActor) {
	s.mu.Lock()

	actorUUID := uuid.New()
	s.actors[actorUUID] = a

	s.mu.Unlock()

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

func (s *RecoverSupervisor) StopActor(id ActorID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	a, ok := s.actors[id]
	if !ok {
		return
	}

	a.stop()
	delete(s.actors, id)
}

func (s *RecoverSupervisor) shutdown() {
	for _, a := range s.actors {
		a.stop()
	}
}
