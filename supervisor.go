package goactor

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

const defaultSupervisorCapacity = 256

type Supervisor interface {
	RegisterActor(a genericActor) ActorID
	StopActor(id ActorID)
	shutdown()
}

type RecoverSupervisor struct {
	mu               sync.Mutex
	actors           map[ActorID]genericActor
	stopActorReqChan chan ActorID
}

func NewRecoverSupervisor(ctx context.Context) (*RecoverSupervisor, func()) {
	s := &RecoverSupervisor{
		actors: make(map[ActorID]genericActor, defaultSupervisorCapacity),
	}

	go func() {
		select {
		case id := <-s.stopActorReqChan:
			s.stopActor(id)
		case <-ctx.Done():
			s.shutdown()
		}
	}()

	return s, s.shutdown
}

func (s *RecoverSupervisor) RegisterActor(a genericActor) {
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
	s.stopActorReqChan <- id
}

func (s *RecoverSupervisor) StopActorAfter(id ActorID, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		s.stopActorReqChan <- id
	}()
}

func (s *RecoverSupervisor) stopActor(id ActorID) {
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
	println("shutting down supervisor")
	for _, a := range s.actors {
		a.stop()
	}
}
