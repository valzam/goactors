package goactor

import "github.com/google/uuid"

type genericActor interface {
	start()
	stop()
	running() <-chan struct{}
	shutdown() <-chan struct{}
}

type ActorID = uuid.UUID
