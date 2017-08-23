package cflib

type StopObserver interface {
	Stop()
}

type StopPublisher struct {
	listeners []StopObserver
}

func (s *StopPublisher) NotifySubscribers() {
	for _, listener := range s.listeners {
		listener.Stop()
	}
}

func (s *StopPublisher) Add(subscriber StopObserver) {
	s.listeners = append(s.listeners, subscriber)
}
