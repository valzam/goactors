package goactor

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
	a.prepareStart()

	go func() {
		for !a.isStopped() {
			func() {
				defer func() {
					if r := recover(); r != nil {
						println("restarting actor")
					}
				}()
				a.start()
			}()
		}
	}()
}

func (s *RecoverSupervisor) shutdown() {
	for _, a := range s.actors {
		a.stop()
	}
}
