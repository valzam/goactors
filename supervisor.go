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
	ctx              context.Context
	mu               sync.Mutex
	actors           map[ActorID]genericActor
	stopActorReqChan chan ActorID
}

func NewRecoverSupervisor(ctx context.Context) (*RecoverSupervisor, func()) {
	s := &RecoverSupervisor{
		actors: make(map[ActorID]genericActor, defaultSupervisorCapacity),
	}
	cancelCtx, cancel := context.WithCancel(ctx)

	go func() {
		select {
		case id := <-s.stopActorReqChan:
			s.stopActor(id)
		case <-cancelCtx.Done():
			s.shutdown()
		}
	}()

	return s, cancel
}

func (s *RecoverSupervisor) RegisterActor(a genericActor) {
	s.mu.Lock()

	actorUUID := uuid.New()
	s.actors[actorUUID] = a

	s.mu.Unlock()

	go func() {
		for {
			select {
			case <-a.shutdown():
				return
			default:
				defer func() {
					if r := recover(); r != nil {
						println("restarting actor")
						time.Sleep(50 * time.Millisecond)
					}
				}()
				a.start()
			}
		}
	}()

	// Wait for the actor to indicate that it is running
	<-a.running()
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
