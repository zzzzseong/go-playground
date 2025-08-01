package semaphore

import (
	"context"
	"golang.org/x/sync/semaphore"
	"sync"
	"testing"
)

// TestMySemaphore tests the ChannelBasedSemaphore implementation.
func TestMySemaphore(t *testing.T) {
	// Create a new semaphore with a limit of 10 goroutines.
	sem := NewChannelBasedSemaphore(10)

	// Use a WaitGroup to wait for all goroutines to finish.
	var wg sync.WaitGroup

	numberOfTasks := 1_000_000
	wg.Add(numberOfTasks)

	// 1️⃣ Without mutex, this should cause a race condition.
	count := 0
	for i := 0; i < numberOfTasks; i++ {
		go func(i int) {
			defer wg.Done()

			sem.Acquire()
			defer sem.Release()

			count += 1
		}(i)
	}
	wg.Wait()

	if count == numberOfTasks {
		t.Errorf("Expected a race condition, but it did not occur. The count is equal to the number of tasks.")
	}

	var mu sync.Mutex
	wg.Add(numberOfTasks)

	// 2️⃣ With mutex, this should ensure that the count is correct.
	count = 0
	for i := 0; i < numberOfTasks; i++ {
		go func(i int) {
			defer wg.Done()

			sem.Acquire()
			defer sem.Release()

			mu.Lock()
			count += 1
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	if count != numberOfTasks {
		t.Errorf("Expected count to be %d, but got %d", numberOfTasks, count)
	}
}

func TestGoSemaphore(t *testing.T) {
	// Create a new semaphore with a limit of 10 goroutines.
	sem := semaphore.NewWeighted(10)
	ctx := context.Background()

	// Use a WaitGroup to wait for all goroutines to finish.
	var wg sync.WaitGroup

	numberOfTasks := 1_000_000
	wg.Add(numberOfTasks)

	// 1️⃣ Without mutex, this should cause a race condition.
	count := 0
	for i := 0; i < numberOfTasks; i++ {
		go func(i int) {
			defer wg.Done()

			// Acquire the semaphore with a context.
			if err := sem.Acquire(ctx, 1); err != nil {
				t.Errorf("Failed to acquire semaphore: %v", err)
				return
			}
			defer sem.Release(1)

			count += 1
		}(i)
	}
	wg.Wait()

	if count == numberOfTasks {
		t.Errorf("Expected a race condition, but it did not occur. The count is equal to the number of tasks.")
	}

	var mu sync.Mutex
	wg.Add(numberOfTasks)

	// 2️⃣ With mutex, this should ensure that the count is correct.
	count = 0
	for i := 0; i < numberOfTasks; i++ {
		go func(i int) {
			defer wg.Done()

			// Acquire the semaphore with a context.
			if err := sem.Acquire(ctx, 1); err != nil {
				t.Errorf("Failed to acquire semaphore: %v", err)
				return
			}
			defer sem.Release(1)

			mu.Lock()
			count += 1
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	if count != numberOfTasks {
		t.Errorf("Expected count to be %d, but got %d", numberOfTasks, count)
	}
}
