package booking

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/rizkyrachman00/miniBookTicket-go-concept-concurrent/internal/adapters/redis"
)

func TestConcurrentBooking_ExactlyOneWins(t *testing.T) {
	store := NewRedisStore(redis.NewClient("localhost:6379"))
	svc := NewService(store)

	const numGoroutines = 100_000 // 100k users trying to book a seat at the same time

	var (
		succeses atomic.Int64
		failures atomic.Int64
		wg       sync.WaitGroup
	)

	wg.Add(numGoroutines)
	for i := range numGoroutines {
		go func(userNum int) {
			defer wg.Done()
			err := svc.Book(Booking{
				MovieID: "screen-1",
				SeatID:  "A1",
				UserID:  uuid.New().String(),
			})

			if err == nil {
				succeses.Add(1)
			} else {
				failures.Add(1)
			}

		}(i)
	}
	wg.Wait()

	if got := succeses.Load(); got != 1 {
		t.Errorf("Expected exactly 1 success, got %d", got)
	}

	if got := failures.Load(); got != int64(numGoroutines-1) {
		t.Errorf("Expected %d failures, got %d", numGoroutines-1, got)
	}
}
