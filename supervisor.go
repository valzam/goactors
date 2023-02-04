package goactor

type Supervisor interface {
	Register(a Actor)
	shutdown()
}

type RecoverSupervisor struct {
	actors []Actor
}

func NewRecoverSupervisor() (*RecoverSupervisor, func()) {
	s := &RecoverSupervisor{}
	return s, func() { s.shutdown() }
}

func (s *RecoverSupervisor) Register(a Actor) {
	s.actors = append(s.actors, a)
	go func() {
		for !a.IsStopped() {
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
