package semaphore

type ChannelBasedSemaphore struct {
	sem chan struct{}
}

func NewChannelBasedSemaphore(limit int) Semaphore {
	return &ChannelBasedSemaphore{
		sem: make(chan struct{}, limit),
	}
}

func (s *ChannelBasedSemaphore) Acquire() {
	s.sem <- struct{}{}
}

func (s *ChannelBasedSemaphore) Release() {
	<-s.sem
}
