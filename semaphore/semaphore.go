package semaphore

type Semaphore interface {
	// Acquire blocks until the semaphore can be acquired.
	Acquire()

	// Release releases the semaphore, allowing other goroutines to acquire it.
	Release()
}
