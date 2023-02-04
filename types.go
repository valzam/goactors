package goactor

import "github.com/google/uuid"

type genericActor interface {
	start()
	stop()
	wasStopped() bool
	isRunning() bool
}

type ActorID = uuid.UUID
